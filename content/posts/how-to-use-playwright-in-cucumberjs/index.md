+++
date = ""
author = "nploi"
description = "How to use Playwright in cucumberjs"
title = "ow to use Playwright in cucumberjs"
categories = ["e2e test"]
tags = ["k8s", "playwright", "cucumberjs"]
slug = "how-to-use-playwright-in-cucumberjs"
+++

This tutorial helps you run automate your test using Playwright in Cucumber.

Before begin, we will give a brief introduction to Cucumber and Playwright.

### Cucumber
Cucumber is a tool that supports Behaviour-Driven Development(BDD).

Ok, now that you know that BDD is about discovery, collaboration and examples (and not testing), let’s take a look at Cucumber.

Cucumber reads executable specifications written in plain text and validates that the software does what those specifications say. The specifications consists of multiple examples, or scenarios. For example:

```
Scenario: Breaker guesses a word
  Given the Maker has chosen a word
  When the Breaker makes a guess
  Then the Maker is asked to score
```

Each scenario is a list of steps for Cucumber to work through. Cucumber verifies that the software conforms with the specification and generates a report indicating ✅ success or ❌ failure for each scenario.

In order for Cucumber to understand the scenarios, they must follow some basic syntax rules, called [Gherkin](https://cucumber.io/docs/gherkin/).

### Playwright
Playwright can either be used as a part of the Playwright Test test runner (this guide), or as a [Playwright Library](https://playwright.dev/docs/library/).

Playwright Test was created specifically to accommodate the needs of the end-to-end testing. It does everything you would expect from the regular test runner, and more. Playwright test allows to:

- Run tests across all browsers.
- Execute tests in parallel.
- Enjoy context isolation out of the box.
- Capture videos, screenshots and other artifacts on failure.
- Integrate your POMs as extensible fixtures.

### Getting Started with Cucumber and Playwright Example

#### Prerequisites and Installations
