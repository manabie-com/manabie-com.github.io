+++
date = "2021-12-03T14:37:23+07:00"
author = "vctqs1"
description = "What is the purpose of Cucumber and End-to-end testing? Why do we need to invest in them?"
title = "Why we use Cucumber for end-to-end testing?"
categories = ["Automation", "End-to-end", "BDD"]
tags = ["bdd", "end-to-end", "cucumber", "automation", "test"]
slug = "why-we-use-cucumber-for-end-to-end-testing"
+++

# Why we use Cucumber for end-to-end testing?

Before knowing why we use cucumber for End-to-end testing. Let understand what is end-to-end testing and cucumber

## What is End to end testing?

End to end testing (E2E testing) is a software testing method that validates entire software from beginning to end.

The purpose of end-to-end testing is to simulate what a real user scenario and testing whole software for dependencies, data integrity and communication with other systems, interfaces, network connectivity and databases to exercise complete production like scenario.

## What is Cucumber?

What is the world said:

> Cucumber is a software tool that supports behavior-driven development (BDD).Central to the Cucumber BDD approach is its ordinary language parser called Gherkin. It allows expected software behaviors to be specified in a logical language that customers can understand. As such, Cucumber allows the execution of feature documentation written in business-facing text. It is often used for testing other software. It runs automated acceptance tests written in a behavior-driven development (BDD) style.

### What is Behavior-driven development?

In software engineering, behavior-driven development (BDD) is an agile software development process that encourages collaboration among developers, quality assurance testers, and customer representatives in a software project.

BDD extends from TDD (Test driven development). Instead of focus on test-first in software development BDD focuses on behavior.

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

Gherkin also has their rule and syntax. As you see Scenario, Given, When, and Then is Gherkin keyword.

Each keyword represents each purpose.

As non-technical person, they can understand it well as a plain text with step by step.

-   Which they have
-   Which they do
-   Which they're result

The technical person also has same look, but needs to implement some things to can integrate with our system. It's called step definitions in Cucumber.

In this step, you can use any language to do like JS, Golang, ...

End-to-end testing is not only software testing. But also is documentation, source of truth.

**That's why cucumber becomes is one of the best choices for us**
-   Cucumber is a popular tools, we can implement of end-to-end testing quickly with many plugin supports like report.
-   Behavior-driven development focus on behavior and it can be become our document how the system actually behaves for millions/billions test cases. More effective than, it can be training document for our members and customers
-   Using Gherkin language supports human-readable ease. Both technical and non-technical are able to read and understand the products.

## Reference

1. [Cucumber Wiki](<https://en.wikipedia.org/wiki/Cucumber_(software)>)
2. [Behavior-driven development Wiki](https://en.wikipedia.org/wiki/Behavior-driven_development)
3. [Introduce BDD](https://automationpanda.com/2017/01/25/bdd-101-introducing-bdd/)
4. [Cucumber](https://cucumber.io/docs/)
5. [End-to-end testing](https://www.guru99.com/end-to-end-testing.html)
