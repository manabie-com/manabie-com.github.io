// features/support/steps.js
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