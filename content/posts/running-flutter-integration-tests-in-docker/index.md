+++
date = "2021-12-28T09:49:36+07:00"
author = "phamnhuvu-dev"
description = "Why and How we run integration tests in Docker"
title = "Running Flutter integration tests in Docker"
categories = ["DevSecOps", "Testing"]
tags = ["docker", "integration-test", "flutter", "dart"]
slug = "running-flutter-integration-tests-in-docker"
+++
# Running Flutter integration tests in Docker

## Prerequisite
- Linux OS
- Installing:
    - Flutter
    - Docker
    - Docker compose
    - KVM https://help.ubuntu.com/community/KVM/Installation

## Technologies

- Flutter:
    - Google's SDK for crafting beautiful, fast user experiences for mobile, web, and desktop from a single codebase. Flutter works with existing code, is used by developers and organizations around the world, and is free and open source.
    - https://github.com/flutter/flutter

- Docker:
    - An open platform for developing, shipping, and running applications
    - https://docs.docker.com/get-started/overview/

- Android Emulator container:
    - Allowing you to find and run the right version of the Emulator without the headache of dependency management, which makes it easy to scale automated tests as part of a CI/CD system without the upkeep cost of a physical device farm
    - https://github.com/google/android-emulator-container-scripts

## The criteria that we use Docker for running automation tests

- Scaling tests:
    - Docker provides the ability to package and run an application in a loosely isolated environment, so we can run multiple e2e tests parallel.

- Reusing:
    - Our Dockerfile file and docker-compose.yml file is simple because of reusing container from the community.

- Consistent:
    - Allowing developers to work in standardized environments using local containers which provide applications and services

- Responsive:
    - Docker containers can run on a developer’s local laptop and on cloud providers.

## Explanation

- We create Dockerfile for the flutter-app service because we need to copy our source code to the container to build then run tests
```
FROM cirrusci/flutter:2.8.1

RUN sdkmanager "build-tools;29.0.2"

RUN flutter precache --android --no-web --no-ios --no-universal

COPY pubspec.* .

RUN --mount=type=cache,sharing=locked,target=/flutter flutter pub get

WORKDIR /project

COPY . .
```

- To run the Flutter integration test we need a process for Android Emulator and a process for Flutter drive. That is why we have 2 services is `flutter-app` and `android-emulator` in `docker-compose.yml` file.
```
version: "3.9"
services:
  android-emulator:
    image: us-docker.pkg.dev/android-emulator-268719/images/30-google-x64:30.1.2
    devices:
      - /dev/kvm

  flutter-app:
    depends_on:
      - android-emulator
    build: .
    command: make run-test
```

- But `flutter-app` and `android-emulator` are services separately. It means they are not the same network, so we need to connect `flutter-app` to `android-emulator` by the following command: `adb connect android-emulator:5555`. Reference: https://developer.android.com/studio/command-line/adb#wireless

- After `flutter-app` connects to `android-emulator`, we can run the Flutter integration tests with `flutter drive` command. We combine connect emulator script and flutter drive script in `run-test` command in `Makefile`.
```
run-test:
	sleep 10
	adb connect android-emulator:5555
	adb wait-for-device
	bash waiting-for-boot.sh
	flutter pub get
	flutter drive --target test_driver/app.dart
```

## Run the example
- Cloning the source code on Github: https://github.com/phamnhuvu-dev/flutter_android_docker
- After cloning the source code, move to the root of the source code.
- Run `make run-docker-test`
- If you see the following log after running the above command. You succeeded.

```
All tests passed!
```

## Summary:
- Using docker for automation tests is used widely in the software industry, not only Mobile field.
- With Docker's ability, scaling tests are easier and faster, helping us detect breakings and bugs quickly to fix and maintain.
- Maybe you have a question that almost all mobile developers use macOS to work. How can run the automation tests on macOS before shipping the source code to the cloud?
    1. We can't run `android-emulator x64` container on macOS because the macOS doesn't have KVM.
    2. Actually, we can run `android-emulator arm` with `-no-accel` but the performance is very bad => can't use on Intel chip macOS
    3. Running `android-emulator arm` on ARM chip macOS:
        - Google is working on that: https://github.com/google/android-emulator-container-scripts/issues/211
        - Self-building `android-emulator arm` image https://github.com/google/android-emulator-container-scripts. I am using Intel macOS so I don't have ARM macOS to do this. I am waiting for the next M1X Mac mini with 32GB ram, 16GB ram is not enough for me for something paralleling and scaling.
    4. Connecting `flutter-app` service from container to Android Emulator on macOS, yes we can do that I will write a blog for this.

