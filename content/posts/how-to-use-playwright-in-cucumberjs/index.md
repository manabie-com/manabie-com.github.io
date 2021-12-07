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

Cucumber is a tool that supports [Behaviour-Driven Development](https://cucumber.io/docs/bdd)(BDD), If you’re new to Behaviour-Driven Development read [BDD introduction](https://cucumber.io/docs/bdd/) first.

[Cucumber-js](https://github.com/cucumber/cucumber-js) is an open-source software testing tool written in Javascript, while the tests are written in Gherkin, a non-technical and human-readable language.

### Playwright

Playwright is a framework for Web Testing and Automation. It allows testing Chromium, Firefox and WebKit with a single API. Playwright is built to enable cross-browser web automation that is ever-green, capable, reliable and fast.

### Getting Started with Cucumber and Playwright Example

#### Prerequisites and Installations

- Installations:
  - Install [Node.js](https://nodejs.org/en/) (12 or higher)
  - Install Cucumber modules with [yarn](https://yarnpkg.com/en/) or [npm](https://www.npmjs.com/)
    - yarn:

        ```bash
        yarn add -D @cucumber/cucumber
        ```
    - npm:

        ```bash
        npm i -D @cucumber/cucumber
        ```

  - Install Playwright
    - yarn:
        ```bash
        yarn add playwright
        ```
    - npm:
        ```bash
        npm i -D playwright
        ```
    - Add the following files:
      - `features/search_manabie_on_google.feature`

        ```feature
        Feature: Search Manabie on google

            Scenario: Bob try to find Manabie on Google
                Given Bob go to Google website
                When Bob search Manabie
                Then Manabie appears on result list
        ```
      - `features/support/world.js`

        ```javascript
        const { setWorldConstructor } = require("@cucumber/cucumber");
        const playwright = require('playwright');

        class CustomWorld {
            constructor() {
                this.variable = 0;
            }

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
        const assert = require("assert").strict;

        Given("Bob go to Google website", async function () {
            await this.openUrl('http://google.com/');
        });

        When("Bob search Manabie", async function () {
            await this.page.click('[aria-label="Tìm kiếm"]');
            await this.page.fill('[aria-label="Tìm kiếm"]', 'manabie');
            await this.page.press('[aria-label="Tìm kiếm"]', 'Enter')
        });

        Then("Manabie appears on result list", async function () {
            await this.page.waitForSelector('text=https://www.manabie.vn');
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

### References

- <https://github.com/cucumber/cucumber-js>
- <https://github.com/microsoft/playwright>
