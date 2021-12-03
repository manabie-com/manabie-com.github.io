# Why we use Cucumber for end-to-end testing?

Before knowing why we use cucumber for End-to-end testing. Let understand what is end-to-end testing and cucumber

## What is End to end testing?

End to end testing (E2E testing) is a software testing method that validates entire software from beginning to end.

The purpose of end-to-end testing is to simulate what a real user scenario and testing whole software for dependencies, data integrity and communication with other systems, interfaces, network connectivity and databases to exercise complete production like scenario.

## What is Cucumber?

What is the world said:

> Cucumber is a software tool that supports behavior-driven development (BDD).Central to the Cucumber BDD approach is its ordinary language parser called Gherkin. It allows expected software behaviors to be specified in a logical language that customers can understand. As such, Cucumber allows the execution of feature documentation written in business-facing text. It is often used for testing other software. It runs automated acceptance tests written in a behavior-driven development (BDD) style.

### What is Behavior-driven development?

The world said:

> In software engineering, behavior-driven development (BDD) is an agile software development process that encourages collaboration among developers, quality assurance testers, and customer representatives in a software project

BDD is

-   specification by example.
-   focuses on behavior first
-   a refinement of the Agile process, not an overhaul.
-   is a paradigm shift

For more detail: https://automationpanda.com/2017/01/25/bdd-101-introducing-bdd/

## Why we use Cucumber for end-to-end testing?

Cucumber is following syntax called Gherkin to explain and validate executable specifications written in plain text.

To define Gherkin in a simple way, it is:

-   Define test cases as plain text, using Gherkin language.
-   It is designed to be non-technical and human-readable, becomes a ubiquitous language between tech and non-tech peoples. And collectively describes use cases relating to a software system.

For example:

```feature
Scenario: Bob applies for a Frontend engineer at Manabie
  Given Bob receives the test challenge for Frontend engineer at Manabie
  When Bob submits the challenge after complete
  Then Manabie reviews his submission
```

Gherkin also has their rule and syntax. As you see Scenario, Given, When, and Then is Gherkin keyword

Each keyword represents each purpose.

As non-technical person, they can understand it well as a plain text with step by step.

-   Which they have
-   Which they do
-   Which they're result

The technical person also has same look, but needs to implement some things to can integrate with our system. It's called step definitions in Cucumber.

In this step, you can use any language to do like JS, Golang, ...

End-to-end testing is not only software testing. But also is documentation, source of truth

That's why cucumber becomes is one of the best choices for us

-   Implementation of end-to-end testing quickly with many plugin supports
-   Behavior-driven development. Document how the system actually behaves.
-   Using Gherkin language supports human-readable easy

## Reference

1. [Cucumber](<https://en.wikipedia.org/wiki/Cucumber_(software)>)
2. [Behavior-driven development Wiki](https://en.wikipedia.org/wiki/Behavior-driven_development)
3. [Introduce BDD](https://automationpanda.com/2017/01/25/bdd-101-introducing-bdd/)
4. [BDD in Cucumber](https://cucumber.io/docs/bdd/)
5. [End-to-end testing](https://www.guru99.com/end-to-end-testing.html)
