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
the server. In the configuration above we have an account named **A** which has third users ***Admin, Bob and Tom***. We 
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

#### **Let's run a test**

You can go to [my demo]()

##### **Setup config and docker-compose**


