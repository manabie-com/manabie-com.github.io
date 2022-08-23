+++
date = "2022-08-23T17:10:49+07:00"
author = "duongcongtoai"
description = "TODO"
title = "J4 stress test framework"
categories = ["DevSecOps", "Testing"]
tags = ["Kubernetes", "integration-test", "Golang", "microservices"]
slug = "j4-stress-test-framework"
+++
## What is stress test
## Why stress test

At Manabie, we have already implemented bunch of features for the business, mostly to support ERP in education domain. The traffic on production is not that high. For each cluster (of each clients), it barely reach 100 rpc on our Grafana dashboard, and everything looks safe. 
But no, we realize that in our roadmap, even though each cluster traffic is not that significant, but when they all combined into one single multi-tenant cluster (which is a topic for another day), we may encouter performance issues unprepared. And thus we have an epic to plan for stress test developing in our Kanban dashboard.

## Why another framework

What we want:
- Programmable interface (Golang): In Manabie, we write a lot of integration test using Golang, it is already a normal routine of every developer, and stress testing script is just another integration test script, but with high load and orchestrated by the stress test framework. Thus, we want our developer to feel indifference between writing an integration test and stress test. Another reason for using programmable interface is that we want to inject custom metadata into the stressed test requests such as prometheus client side metrics, or Jaeger trace.
- Scalable deployment together with K8S: this is a quality we want in a framework (even though we may don't need it right now), that it can scaled and integrate well with K8S eco-system. We want to provide utility to our developer to stress test their service without complex setup: how many request to simulate, how fast the rpc increase, gradually or exponentially,... And the framework should behave correctly with what it promise. In order to do that, it must be deployed in cluster mode, into multiple containers with certain level of fault tolerant (aka one worker goes down, the workload must be handled by another newly spawned worker)


K6s and Locust are good options but they do not provide the option for a Golang friendly programmable interface. You may argue that language does not matter, but for us in this case it does.

## J4 (Jarvan)
![J4](./images/j4-icon.png)
### Core concept
Scenario
Task allocator
Task executor
Rampup cycle
Rampdown cycle
Hold cycle
### Cluster mode
### Custom script
#### Simple requests
#### Long living requests
### Observability
#### Client side metrics gathering
#### J4 internal metrics

## Room for improvement
## Reference

