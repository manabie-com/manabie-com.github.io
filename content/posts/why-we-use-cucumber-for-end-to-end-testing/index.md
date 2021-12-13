+++
date = "2021-12-13T14:37:23+07:00"
author = "vctqs1"
description = "What is the purpose of Cucumber and End-to-end testing? Why do we need to invest in them?"
title = "Why we use Cucumber for end-to-end testing?"
categories = ["DevSecOps", "Testing"]
tags = ["bdd", "end-to-end", "cucumber", "automation", "test"]
slug = "why-we-use-cucumber-for-end-to-end-testing"
+++

# Why we use Cucumber for end-to-end testing?

Before getting to know why we use cucumber for End-to-end testing. Let's understand what's end-to-end testing and cucumber are.

## What is End to end testing?

End to end testing (E2E testing) is a software testing method that validates the entire software. From the beginning to the end.

The purpose of end-to-end testing is to simulate how a real user interacts with the application, list out the scenarios, and test the whole software for dependencies, data integrity and communication with other systems, interfaces, network connectivity and databases to exercise complete production like scenario.

## What is Cucumber?

The definition:

> Cucumber is a software tool that supports behavior-driven development (BDD).Central to the Cucumber BDD approach is its ordinary language parser called Gherkin. It allows expected software behaviors to be specified in a logical language that customers can understand. As such, Cucumber allows the execution of feature documentation written in business-facing text. It is often used for testing other software. It runs automated acceptance tests written in a behavior-driven development (BDD) style.

### What is Behavior-driven development?

In software engineering, behavior-driven development (BDD) is an agile software development process which encourage collaboration amongst developers, QAs (Quality Assurance team), and customer representatives in a project.

BDD is an extension from TDD (Test driven development). Instead of focus on test-first, BDD focuses on users behavior.

## Why we use Cucumber for end-to-end testing?

Cucumber is following syntax called Gherkin to explain and validate executable specifications written in plain text.

To define Gherkin in a simple way, it is:

-   Define test cases as plain text, using Gherkin language.
-   It is designed to be non-technical and human-readable, becomes a ubiquitous language between tech and non-tech peoples. And collectively describes use cases relating to a software system.

For example:

```feature
Scenario: Bob try to find Manabie on Google
  Given Bob go to Google website
  When Bob search Manabie
  Then Manabie appears on result list
```

Gherkin also follow rules and syntaxes. As you can see, Gherkin keywords are:

-   Scenario
-   Given
-   When
-   Then
-   ...

Each keyword has their purpose. As a person with no technical background, you can understand it well as plain text.

-   Given: What we have
-   When: What we do
-   Then: What's the expected result

A person with technical background also look at the same thing with different perspective. We need to define steps for Cucumber to carry out the scenarios. It called `Step Definitions` in Cucumber.

In this step, you can use `a wide varieties of programming languages` to implement, such as JS, Golang, etc

End-to-end testing is not only a software testing method, but also a means to documentation, a source of truth.

**That's why cucumber becomes is one of the best choices for us**

-   Cucumber is a popular tool, we can implement end-to-end testing quickly and support many plugins (like report).
-   Behavior-driven development focus on behavior. And it can becomes our document on how the system behaves with millions/billions test cases. More over, it can be a training document for our members and customers.
-   Using Gherkin language supports human-readable ease. Both technical and non-technical are able to read and understand the products.

## Reference

1. [Cucumber Wiki](<https://en.wikipedia.org/wiki/Cucumber_(software)>)
2. [Behavior-driven development Wiki](https://en.wikipedia.org/wiki/Behavior-driven_development)
3. [Introduce BDD](https://automationpanda.com/2017/01/25/bdd-101-introducing-bdd/)
4. [Cucumber](https://cucumber.io/docs/)
5. [End-to-end testing](https://www.guru99.com/end-to-end-testing.html)
