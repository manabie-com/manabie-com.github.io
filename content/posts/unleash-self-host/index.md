+++
date = "2021-12-09T15:11:00+07:00"
author = "anhpngt"
description = "Self-hosting Unleash with Kubernetes"
title = "Self-hosting Unleash with Kubernetes"
categories = ["DevSecOps", "Feature"]
tags = ["Unleash", "Kubernetes", "helm"]
slug = "unleash-self-host"
+++

In this blog post, we will learn how to self-host Unleash in a Kubernetes cluster.

### Why feature toggles?

While developing new features for our end-users, we often encounter these 2 problems:

- Features are not developed and rolled out in a single night. It often takes several days or even
weeks before a feature is fully developed, tested, and deployed. In such cases, we usually have to
deploy a piece of that feature to production, and we would need to hide that piece until the feature
is fully completed and ready.
- Even after the feature reaches end-users, until we are fully confident in the feature, fail-safe
methods are usually employed. One of such methods is to hide the feature away from the users to prevent
it from being used.

With the feature toggle, we then can simply turn the feature's flag off and on to hide or show the
feature to users.

### Why Unleash?

Before [Unleash](https://www.getunleash.io/), we have tried [Firebase Remote Config](https://firebase.google.com/docs/remote-config)
to some success. However, what it did not have was **a local deployment or an emulator**.
This feature was crucial to us, because:

1. Local development: in Manabie, any developer can spin up the entire end-to-end infrastructure
in their own machine and start working on their task without having to worry about breaking any
of the production clusters. However, Firebase Remote Config is a shared instance. Using Firebase
would not meet our separation-of-concern's standards.
2. CI/CD: when running end-to-end tests, it is desirable to run the tests against different feature
toggle configurations. We need to ensure that our code works with both cases of the flag being turned
on and off. It would be disastrous if it does not.

### Deploying Unleash in a Kubernetes cluster

#### 1. Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [minikube](https://minikube.sigs.k8s.io/docs/)
- [helm](https://helm.sh/)

Using `helm` is a bit of an overkill here. However, our Manabie's CI/CD pipeline uses `helm` to deploy
so we will use it here as well.

For this guide, I am using the following versions

{{< gist anhpngt 542f42d65d08d480a5860e8bb790624d version >}}

The versioning requirements are not strict. However, if you encounter any strange errors, you can
try installing the listed versions first.

#### 2. Setting up the project

```sh
$ helm start
$ helm cache add 

```