+++
date = "2021-12-10T14:28:23+07:00"
author = "duchh"
description = "Write the way to set up nats multi-tenant."
title = "Set up nats multi-tenant in golang"
categories = ["Golang", "Nats-Jetstream"]
tags = ["Golang", "Nats-Jetstream", "Multi-tenant"]
slug = "set-up-nats-multi-tenant-in-golang" 
+++

***In this blog post, we will learn how to set up Nats Multi-tenant in golang.***

#### **Nats-Jetstream**

Nats has a built-in distributed persistence system called Jetstream which features new functionalities and higher 
qualities of service on top of the base Core NATS functionalities and qualities of service.

Jetstream was created to solve the problems identified with streaming in technology today. Some technologies address 
these better than others, but no current streaming technology is truly multi-tenant, horizontally scalable,... Today I just 
talk about the way set up nats multi-tenant in golang.

#### **Concepts**

Look in this example, I will explain below.

```bash
accounts {
    A: {
        jetstream: enabled
        users = [
            {
                user: "Admin",
                password: "123456",
                permissions: {
                    publish: {
                        allow: [
                            "$JS.API.INFO", # allow user can send request API to nats server
                            "$JS.API.STREAM.>", # allow user do anything with all streams in nats server
                            "$JS.API.CONSUMER.>" # allow user do anything with all consumers in nats server
                        ]
                    }
                }
            },
            {
                user: "Bob",
                password: "123456",
                permissions: {
                    publish: {
                        allow: [
                            "$JS.API.INFO", # allow user can send request API to nats server
                            "$JS.API.STREAM.*.student", # allow user interact with stream (example: create, delete,...)
                            "student.Created" # allow user publish messages with subject `student.Created`
                        ]
                    }
                }
            },
            {
                user: "Tom",
                password: "123456",
                permissions: {
                    publish: {
                        allow: [
                            "$JS.API.INFO", # allow user can send request API to nats server
                            "$JS.API.STREAM.NAMES", # allow user get all stream's names
                            "$JS.API.CONSUMER.*.student.*", # allow user interact with consumer (example: create, delete,...)
                            "$JS.API.CONSUMER.DURABLE.CREATE.student.*", # allow user create consumer queue
                            "$JS.ACK.student.>" # allow user can ack messages in `student` stream
                        ]
                    },
                    subscribe: {
                        allow: [
                            "_INBOX.>"
                        ]
                    }
                }
            }
        ]
    }
}
```

##### **Accounts**

Accounts allow grouping of clients, isolating them from clients in other accounts, thus enabling ***multi-tenant*** in 
the server. In the configuration above we have an account named **A** which has three users ***Admin, Bob and Tom***. We 
can add one or more accounts to the configuration, an account has one or more users.

##### **Users**

User identified by the **user**, the **password**, the **permissions**. Users in the account are simply the minimum number 
of services that must work together to provide some functionality.

##### **Permission**

Nats has concept [widlcard subject](https://docs.nats.io/nats-concepts/subjects#wildcards) will use a lot in our example, 
so you should understand this concept before going next.

Another concept is [Jetstream wire API Reference](https://docs.nats.io/reference/reference-protocols/nats_api_reference), 
you can easily interact with the Jetstream infrastructure programmatically.

In [permissions](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/authorization#permission-map) 
above we have two section *publish* and *subscribe*. 
* **Publish** consist of subject, list of subjects, API reference the user can publish. All users 
need to have `$JS.API.INFO` permission, because this API allow user can send others API to nats server.

* **Subscribe** consist of subject, list of subjects, API reference the user can subscribe. To receiving messages 
published, the consumer needs to be able to subscribe to the request subjects. The request subject is an inbox. Typically 
inboxes start with the prefix ``_INBOX.`` followed by a generated string. The ``_INBOX.>`` subject matches all subjects 
that begin with ``_INBOX.``.

In example above, Admin can do anything with streams and consumers (create, update, delete,...). Bob is a publisher, 
publishes messages to nats with streams `student`. Bob can create streams by `$JS.API.STREAM.*.student` and publish 
messages by `student.Created`. Tom is a consumer, receives messages and processes them related streams `student`. Tom 
can create *consumer queue* by `$JS.API.CONSUMER.*.student.*` and `$JS.API.CONSUMER.DURABLE.CREATE.student.*`, then 
ack messages by `$JS.ACK.student.>`.

#### **Let's make a test**

You can go to [my demo](https://github.com/manabie-com/manabie-com.github.io/tree/setup-nats-multi-tenant/content/posts/nats-multi-tenant/examples) 
for more detail.

##### **Setup config and docker-compose**

In this demo I just configure one account and three users, you can custom the test with more than accounts. We run command 
`docker-compose up` to start nats server, the logs when start nats successfully.

```bash
n1    | [1] 2021/12/14 08:12:08.709249 [INF] Starting nats-server
n1    | [1] 2021/12/14 08:12:08.709303 [INF]   Version:  2.6.6
n1    | [1] 2021/12/14 08:12:08.709306 [INF]   Git:      [878afad]
n1    | [1] 2021/12/14 08:12:08.709310 [INF]   Name:     NBWOQ2DYVAYSY3HXYVJXGGURKJY6TBPC2JCKEQYND7GTSSJODGYUCMLK
n1    | [1] 2021/12/14 08:12:08.709318 [INF]   Node:     Cq4JgPlu
n1    | [1] 2021/12/14 08:12:08.709321 [INF]   ID:       NBWOQ2DYVAYSY3HXYVJXGGURKJY6TBPC2JCKEQYND7GTSSJODGYUCMLK
n1    | [1] 2021/12/14 08:12:08.709325 [WRN] Plaintext passwords detected, use nkeys or bcrypt
n1    | [1] 2021/12/14 08:12:08.709336 [INF] Using configuration file: /jetstream.config
n1    | [1] 2021/12/14 08:12:08.710457 [INF] Starting JetStream
n1    | [1] 2021/12/14 08:12:08.710696 [INF]     _ ___ _____ ___ _____ ___ ___   _   __  __
n1    | [1] 2021/12/14 08:12:08.710708 [INF]  _ | | __|_   _/ __|_   _| _ \ __| /_\ |  \/  |
n1    | [1] 2021/12/14 08:12:08.710712 [INF] | || | _|  | | \__ \ | | |   / _| / _ \| |\/| |
n1    | [1] 2021/12/14 08:12:08.710714 [INF]  \__/|___| |_| |___/ |_| |_|_\___/_/ \_\_|  |_|
n1    | [1] 2021/12/14 08:12:08.710717 [INF] 
n1    | [1] 2021/12/14 08:12:08.710720 [INF]          https://docs.nats.io/jetstream
n1    | [1] 2021/12/14 08:12:08.710722 [INF] 
n1    | [1] 2021/12/14 08:12:08.710725 [INF] ---------------- JETSTREAM ----------------
n1    | [1] 2021/12/14 08:12:08.710737 [INF]   Max Memory:      4.00 GB
n1    | [1] 2021/12/14 08:12:08.710742 [INF]   Max Storage:     10.00 GB
n1    | [1] 2021/12/14 08:12:08.710745 [INF]   Store Directory: "/data/jetstream"
```

##### **Golang code**
* **Publisher** folder try connect to nats with user *Bob* then checking and creating streams `student` after that publish
three messages with subject `student.Created`. Go to publisher folder then run command `go run publisher.go`.

```bash
2021/12/14 15:21:02 Student with StudentID:1 has been published
2021/12/14 15:21:02 Student with StudentID:2 has been published
2021/12/14 15:21:02 Student with StudentID:3 has been published
```

* **Consumer** folder try connect to nats with user *Tom* then [creating subscription](https://github.com/manabie-com/manabie-com.github.io/blob/a9b1ec1f6b87c3b408516014a06850869c7b30f8/content/posts/nats-multi-tenant/examples/consumer/consumer.go#L29) 
and processing messages with subject `student.Created`. Go to consumer folder then run command `go run consumer.go`.

```bash
2021/12/14 15:26:15 Student with StudentID:1 has been processed
2021/12/14 15:26:15 Student with StudentID:2 has been processed
2021/12/14 15:26:15 Student with StudentID:3 has been processed
```

* What's happen when Tom don't have permission to ack messages? Let's test. We will remove `$JS.ACK.student.>` in user 
Tom. Then restart nats server. As expected, Tom doesn't have permission to ack messages

The logs in consumer will look like that.
```bash
nats: Permissions Violation for Publish to "$JS.ACK.student.durable-push.1.70.97.1639472208769956133.2" on connection [12]
nats: Permissions Violation for Publish to "$JS.ACK.student.durable-push.1.71.98.1639472208770655160.1" on connection [12]
nats: Permissions Violation for Publish to "$JS.ACK.student.durable-push.1.72.99.1639472208771104510.0" on connection [12]
```
The logs in nats server will look like that.
```bash
n1    | [1] 2021/12/14 08:56:56.427316 [ERR] 172.25.0.1:49916 - cid:12 - "v1.13.0:go" - Publish Violation - User "Tom", Subject "$JS.ACK.student.durable-push.1.70.97.1639472208769956133.2"
n1    | [1] 2021/12/14 08:56:56.427353 [ERR] 172.25.0.1:49916 - cid:12 - "v1.13.0:go" - Publish Violation - User "Tom", Subject "$JS.ACK.student.durable-push.1.71.98.1639472208770655160.1"
n1    | [1] 2021/12/14 08:56:56.427369 [ERR] 172.25.0.1:49916 - cid:12 - "v1.13.0:go" - Publish Violation - User "Tom", Subject "$JS.ACK.student.durable-push.1.72.99.1639472208771104510.0"
```
##### **Nats-box**

We can see the connection of users in account by [nats-box](https://github.com/nats-io/nats-box). 
Go to terminal and run `docker run -ti --network host natsio/nats-box`, then we will connect to nats server by Admin 
`nats context save --server=nats://localhost:4223 --user=Admin --password=123456 --select server` if connect successfully, the result like this

```bash
NATS Configuration Context "server"

      Server URLs: nats://localhost:4223
         Username: Admin
         Password: *********
             Path: /root/.config/nats/context/server.json
       Connection: OK
```
Then we run command `nats server report accounts --json`, the result:

```bash
[
  {
    "account": "A",
    "connections": 3,
    "connection_info": [
      {
        "cid": 10,
        "ip": "172.26.0.1",
        "port": 53406,
        "start": "2021-12-15T01:58:17.622191984Z",
        "last_activity": "2021-12-15T01:58:17.626455989Z",
        "rtt": "449µs",
        "uptime": "16m32s",
        "idle": "16m32s",
        "pending_bytes": 0,
        "in_msgs": 4,
        "out_msgs": 4,
        "in_bytes": 168,
        "out_bytes": 731,
        "subscriptions": 1,
        "lang": "go",
        "version": "1.13.0",
        "authorized_user": "Bob",
        "account": "A",
        "subscriptions_list": [
          "_INBOX.a1ZGIA7xxpq6odeC0nRjid.*"
        ]
      },
      {
        "cid": 13,
        "ip": "172.26.0.1",
        "port": 53476,
        "start": "2021-12-15T01:58:58.408294005Z",
        "last_activity": "2021-12-15T01:59:02.421454384Z",
        "rtt": "530µs",
        "uptime": "15m51s",
        "idle": "15m47s",
        "pending_bytes": 0,
        "in_msgs": 26,
        "out_msgs": 26,
        "in_bytes": 125,
        "out_bytes": 2146,
        "subscriptions": 2,
        "lang": "go",
        "version": "1.13.0",
        "authorized_user": "Tom",
        "account": "A",
        "subscriptions_list": [
          "_INBOX.pAs3HKEavN4xPgq0R3ISuv.*",
          "_INBOX.KUZs5Iu1AS627v4hxQDEDT"
        ]
      },
      {
        "cid": 16,
        "ip": "172.26.0.1",
        "port": 54208,
        "start": "2021-12-15T02:14:50.121066448Z",
        "last_activity": "2021-12-15T02:14:50.121668124Z",
        "rtt": "602µs",
        "uptime": "0s",
        "idle": "0s",
        "pending_bytes": 0,
        "in_msgs": 0,
        "out_msgs": 0,
        "in_bytes": 0,
        "out_bytes": 0,
        "subscriptions": 1,
        "name": "NATS CLI Version development",
        "lang": "go",
        "version": "1.11.0",
        "authorized_user": "Admin",
        "account": "A",
        "subscriptions_list": [
          "_INBOX.4FRAGOb5MzOIGS1ZsYAtYI"
        ]
      }
    ],
    "in_msgs": 30,
    "out_msgs": 30,
    "in_bytes": 293,
    "out_bytes": 2877,
    "subscriptions": 4
  }
]
```
You can see the result above, we have an account A and three users (Admin, Bob, Tom) are connecting to nats. Nats say 
`solating them from clients in other accounts`, so we can test by add account **B** in our configuration then connect nats-box 
by **AdminB** then run `nats server report accounts --json`. 

#### **Summary**

Currently, the Manabie team are using this way to apply nats multi-tenant in our backend code, this way quite easy to configure permission for users 
we just go to file config and add permission into a user you want. You can go to [my demo](https://github.com/manabie-com/manabie-com.github.io/tree/setup-nats-multi-tenant/content/posts/nats-multi-tenant/examples) 
and try to run this example. Thank you for your reading!