+++
date = "2021-12-29T10:45:33Z"
author = "nploi"
description = "This tutorial helps you run web automation tests, using Playwright in Cucumber."
title = "How to use Playwright in cucumberjs"
categories = ["DevSecOps", "Testing"]
tags = ["bdd", "end-to-end", "cucumber", "automation", "test"]
slug = "how-to-use-playwright-in-cucumberjs"
+++

In our previous blog article [HERE](https://blog.manabie.io/2021/12/why-we-use-cucumber-for-end-to-end-testing/), we explained why we're using Cucumber for E2E test at Manabie, and provided a brief understanding about Cucumber.

This time, let's dive a little bit deeper!

Before we start, let's cover a brief introduction to Cucumber again, and afterward, about Playwright.
### [Cucumber](https://cucumber.io/)

Cucumber is a tool that supports [Behaviour-Driven Development](https://cucumber.io/docs/bdd)(BDD), If you’re new to Behaviour-Driven Development read [BDD introduction](https://cucumber.io/docs/bdd/) first.

#### [Cucumberjs](https://github.com/cucumber/cucumber-js)
-  is an open-source software testing tool written in Javascript, while the tests are written in Gherkin, a non-technical and human-readable language.

#### Gherkin Syntax
Gherkin uses a set of special keywords to give structure and meaning to executable specifications. Each keyword is translated to many spoken languages; in this reference we’ll use English.

Each line that isn’t a blank line has to start with a Gherkin keyword, followed by any text you like. The only exceptions are the feature and scenario descriptions.

The primary keywords are:

- Feature
- Rule (as of Gherkin 6)
- Example (or Scenario)
- Given, When, Then, And, But for steps (or *)
- Background
- Scenario Outline (or Scenario Template)
- Examples (or Scenarios)

There are a few secondary keywords as well:

- `"""` (Doc Strings)
- `|` (Data Tables)
- `@` (Tags)
- `#` (Comments)

### [Playwright](https://github.com/microsoft/playwright)

Playwright is a framework for Web Testing and Automation. It allows testing Chromium, Firefox and WebKit with a single API. Playwright is built to enable cross-browser web automation that is ever-green, capable, reliable and fast.
#### Capabilities
Playwright is built to automate the broad and growing set of web browser capabilities used by Single Page Apps and Progressive Web Apps.

- Scenarios that span multiple page, domains and iframes
- Auto-wait for elements to be ready before executing actions (like click, fill)
- Intercept network activity for stubbing and mocking network requests
- Emulate mobile devices, geolocation, permissions
- Support for web components via shadow-piercing selectors
- Native input events for mouse and keyboard
- Upload and download files

### Getting Started with Cucumber and Playwright Example

#### Prerequisites and Installations
- Prerequisites:
  - [Node.js](https://nodejs.org/en/) (12 or higher)
- Installations:
  - Install Cucumber modules with [yarn](https://yarnpkg.com/en/) or [npm](https://www.npmjs.com/)
    - yarn:

        ```bash
        yarn add @cucumber/cucumber
        ```
    - npm:

        ```bash
        npm i @cucumber/cucumber
        ```

  - Install Playwright
    - yarn:
        ```bash
        yarn add playwright
        ```
    - npm:
        ```bash
        npm i playwright
        ```
    - Add the following files:
      - `features/search_job_openings_at_manabie.feature`

        ```feature
        Feature: Search job openings at Manabie

            Scenario: Bob search job openings at Manabie
                Given Bob opens Manabie website
                When Bob goes to Careers section
                Then Bob sees all job openings at Manabie
        ```
      - `features/support/world.js`

        ```javascript
        const { setWorldConstructor } = require("@cucumber/cucumber");
        const playwright = require('playwright');

        class CustomWorld {
            async openUrl(url) {
                const browser = await playwright.chromium.launch({
                    headless: false,
                });
                const context = await browser.newContext();
                this.page = await context.newPage();
                await this.page.goto(url);
            }
        }

        setWorldConstructor(CustomWorld);
        ```

      - `features/support/steps.js`

        ```javascript
        const { Given, When, Then } = require("@cucumber/cucumber");

        Given("Bob opens Manabie website", { timeout: 60 * 1000 }, async function () {
            await this.openUrl('http://manabie.com/');
        });

        When("Bob goes to Careers section", async function () {
            await this.page.click('text=Careers');
        });

        Then("Bob sees all job openings at Manabie", async function () {
            await this.page.click('text=View Openings');
            await this.page.waitForSelector('text=Our Openings');
        });
        ```

    - Run:

        ```bash
        ./node_modules/.bin/cucumber-js --exit
        ```

    - After run:

        ```bash
        1 scenario (1 passed)
        3 steps (3 passed)
        0m03.739s (executing steps: 0m03.729s)
        ```

### Conclusion

Cucumber and Playwright are great frameworks. I hope this article will be of some use to you. [Here is source code](/content/posts/how-to-use-playwright-in-cucumberjs/example), thank you.

Curious for more about how we're tackling testing problems at Manabie? Join us as an Engineer [HERE](https://manabie.breezy.hr/), or follow our blog for future articles.

### References

- <https://github.com/cucumber/cucumber-js>
- <https://github.com/microsoft/playwright>
