+++
date = "2022-04-01T10:00:00+07:00"
author = "duchh"
description = "Explaining how Raft work to achieve consensus between nodes"
title = "Raft consensus algorithms"
categories = ["Raft", "Consensus algorithms"]
tags = ["Raft", "Consensus-algorithms"]
slug = "raft-consensus-algorithms"
+++

***In this blog post, we will learn about Raft and the way it works to achieve consensus between nodes.***

#### **1. Consensus algorithms**

Consensus algorithms allow a collection of machines to work as a coherent group that can survive the failures of some of its members.
Consensus algorithms guarantee is cluster availability, for example, when the sever leader is crashes the cluster will have another 
server instead as soon as possible and all of servers in the cluster will consensus data together, because the server leader always 
replicate logs to other servers. Consensus algorithms properties include: safety, availability, integrity.

![alt](./images/master-slave.png)

#### **2. Raft algorithms**
Raft is a consensus algorithms that uses a replicated log data structure. Raft implements this by first electing a server leader, then giving the 
leader complete responsibility for managing the replicated log. The leader accepts log entries from client, replicate them on other servers, 
and tells servers when it is safe to apply log entries to their state machines.

##### **2.1. Term**
* **Term** is number with consecutive integers. Each term begins with an election, in which one ore more servers attempt to become leader. 
if a server wins the election, then it servers as leader for the rest of the term, it will replicated log to other servers until it is crashes 
(this term end). Or the term will end when the cluster doesn't have any leader is chosen (this case is `split vote` i will explain the 
reason on 3.1.I)

![alt](./images/term.png)

##### **2.2. State**
All nodes in Raft have some states such as:
* **Follower**: This state just interacts with the leader and the candidate. If the leader heartbeat occurs timeout (it is crashes), 
the follower will change state to candidate or it will vote for another the candidate to become the leader when it is received RequestVote.
* **Candidate**: This state have a mission is try to become the leader. When a follower become a candidate, the term will be increased by one, 
then send RequestVote to other servers, if it has vote number over a half of total servers in the cluster, it will be become the leader, 
if not it will back to the follower.
* **Leader**: Only the leader can interact with the client. The leader has mission is replicate data for all nodes in the cluster. In the cluster 
just only has a leader exist. If the leader occurs crash, it will back to Follower.

##### **2.3. Log replication**

![alt](./images/log.png)

Logs are composed of entries, which are numbered sequentially. Each entry contains the term in which it was created (the number in each box) 
and a command for the state machine. An entry is considered committed if it is safe for that entry to be applied to state machines. As you can see 
on the image above, committed entries from 1 to 7, because the cluster has three nodes with the same value in total 5 nodes.

* **Data**: log index, term number, the value is contained on each entry.
* **Log committed**: `Log committed` means logs are replicated on over a half of total servers in the cluster. These logs will be 
replicated for all servers in the cluster. 
* **Log up-to-date**: Raft will compare term and log entry between two nodes when leader election is happen, for example, we have node A and B.
If node A has term is greater node B, so node A is up-to-date. Another way, if the term of two nodes is equal, if node B has log entry longer than 
node A, so node B is up-to-date.

##### **2.4. Request type**
In order to interact between nodes in the cluster, Raft uses some request-types such as:

* **RequestVote**: it uses for, the candidate will send RequestVote to other nodes for voting it become the leader. As you can see the image below, 
the request has arguments, results and implementation conditions. Regarding condition (2), if the result (votedFor) is null but the candidate is 
up-to-date, the vote will be for it.

![alt](./images/request-vote.png)

* **AppendEntries**: The leader uses this request to replicate log entries for all nodes in the cluster. This request is also has mission is send 
the leader heartbeat to other servers.

![alt](./images/append-entries.png)

##### **2.5. Safety**

Raft always guarantees these properties:

* **Election Safety**: at most one leader can be elected in a given term.
* **Leader Append Only**: a leader never overwrites or deletes entries in its logs. It only append new entries.
* **Log Matching**: If two logs contain an entry with the same index and term, then the logs are identical in all entries up through 
the given index
* **Leader Completeness**: If a log entry is committed in a given term, then that entry will be present in the logs of the leader for all 
higher-numbered terms.
* **State Machine Safety**: If a server has applied a log entry at a given index to its state machine, no other server will apply a different 
log entry for the same index. 

#### **3. How does raft work?**

##### **3.1 Leader election**

![alt](./images/leader-election.png)

* You can see the image above, I will explain the process of Leader election: the leader uses AppendEntries to send heartbeat to followers, 
timeout occurs when the follower can not receive heartbeat from the leader. This follower will change state to candidate, increase term 
by one, it will vote for itself, then send RequestVote to other servers. If this candidate has vote number over a half of total servers in 
the cluster, it will become leader and begin to send AppendEntries (heartbeat) to other servers, to avoid timeout from another server and leader 
election happen again.

* **We have two case can happen here**:
    * **3.1.I** When the process of leader election end, it may happen that there is no leader at all. The reason is we have two candidates 
    and they have vote number is equal, this case is called **split vote**. In order to solve this problem, Raft chooses the solution 
    `randomize the timeout` of followers (150-300ms). This solution solves that no two candidates exist at the same time. When the leader 
    is crashes, the follower with the smallest timeout will quickly recognize and start leader election. You can see the image below: 

    ![test](./images/leader-election-timeout.png)

    * **3.1.II** During the process leader election, the candidate can be received RequestVote from another candidate. When this case happens,
    Raft will compare log entries of two candidates, which candidate isn't up-to-date, it will stop sending RequestVote and back to follower.

##### **3.2. Log replication**

* **Leader will replicate log for all nodes in the cluster. Raft manages consensus by Leader will overwrite logs are uncommitted**

* When the Client sends a request to the cluster, the request will be received and processed by leader. Then leader sends AppendEntries with two
elements include: [ N: `{log entry, term}`, N-1: `{log entry, term}` ] (**N is nextIndex**). The follower will compare its data at 
latestIndex with **N-1**, if data is the same, follower will append N into its log entry and return success. If data is different, follower 
will return result is fail, then leader decreases index by one and send AppendEntries again with three elements include: [ N: `{log entry, term}`, 
N-1: `{log entry, term}`, N-2: `{log entry, term}` ]. The follower will compare its data at latestIndex with **N-2**. This operation will 
be repeated until data between follower and leader is the same, then leader will overwrite all of logs in follower `(N, N-1, N-2, N-n,...)`. 
    * For example, you can see the image below, the leader has term is 7, nextIndex is 11. The follower (a) is the same log with leader at index 4, so 
    leader will overwrite all logs of follower (a) at index 4. The follower (b) is the same log with leader at index 3, so the leader will 
    overwrite all logs of follower (b) at index 3.

![alt](./images/log-replication.png)

* **What happens when log replication is processing but Leader is crashed?** 

    ![alt](./images/log-replication-up-to-date.png)

    **1.** At (a), S1 is Leader is replicating log with term-2 for S2, S3, S4, S5. But S1 is crashes.

    **2.** At (b), After leader election end, S5 becomes leader with term-3, but S5 crashes immediately.

    **3.** At (c), S5 crashes, S1 restarts, is elected leader with term-4 and continues replication. At this point, we have two case:

    **3.1** At (d), S1 crashes, it hasn't completed replicate term-2 for S3 yet. Then S5 is elected leader by vote from S3, S4 and itself, 
    then S5 will overwrite term-3 for all nodes in the cluster.

    **3.2** At (e), S1 crashes, it has completed replicate term-2 for S3 and committed. At this point, S5 cannot be elected leader 
    because its log isn't up-to-date.

#### **Conclusions**

Raft is the de-facto standard today for achieving consistency in modern distributed systems. It is designed to be easily understandable 
than Paxos algorithm, which is very hard to understand and implement. Any node in the cluster can become the leader. So, it has a 
certain degree of fairness. Some technologies like: Etcd, MongoDB, NATS,... are using it. 

* You can read more about Raft: https://raft.github.io/
* Raft is implemented by Go: https://github.com/yunuskilicdev/distributedsystems/tree/master/src/raft

