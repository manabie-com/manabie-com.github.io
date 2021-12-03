# Why we use Cucumber for end-to-end testing?

Before knowing why we use cucumber for End-to-end testing. Let understand what is end-to-end testing and cucumber

## What is End to end testing?

End to end testing (E2E testing) is a software testing method that validates entire software from beginning to end.

The purpose of end-to-end testing is to simulate what a real user scenario and testing whole software for dependencies, data integrity and communication with other systems, interfaces, network connectivity and databases to exercise complete production like scenario.

## What is Cucumber?

What is the world said:

> Cucumber is a software tool that supports behavior-driven development (BDD).Central to the Cucumber BDD approach is its ordinary language parser called Gherkin. It allows expected software behaviors to be specified in a logical language that customers can understand. As such, Cucumber allows the execution of feature documentation written in business-facing text. It is often used for testing other software. It runs automated acceptance tests written in a behavior-driven development (BDD) style.

## What is BDD?



## Why we use Cucumber for end-to-end testing?

Cucumber is following syntax called Gherkin to explain and validate executable specifications written in plain text.

To define Gherkin in a simple way, it is:
-   Define test cases as plain text, using Gherkin language.
-   It is designed to be non-technical and human-readable, becomes a ubiquitous language between tech and non-tech peoples. And collectively describes use cases relating to a software system.

For example:

```feature
Scenario: Bob applies for a Frontend engineer at Manabie
  Given Bob receives the test challenge for Frontend engineer at Manabie
  When Bob submit the challenge after complete
  Then Manabie reviews his submission
```

Gherkin also has their rule and syntax. As you see Scenario, Given, When and Then is Gherkin keyword

Each keyword represents for each purpose.

As non-technical person, they can understand it well as a plain text with step by step.

-   Which they have
-   Which they do
-   Which they're result

Technical person also have same look, but need to implement some things to can integrate with our system. It's called step definitions in Cucumber.

In this steps, you can use any language to do like: JS, Golang, ...


End-to-end testing is not only software testing. But also is documentation, source of truth

Cucumber provides easy way to

-   Implementation end-to-end testing quickly with many plugin supports
-   Behavior-driven development. Document how the system actually behaves.
-   Using Gherkin language supports human-readable easy


## Reference
1. 
