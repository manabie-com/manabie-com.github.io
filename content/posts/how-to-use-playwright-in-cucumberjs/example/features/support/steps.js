const { Given, When, Then } = require("@cucumber/cucumber");

Given("Bob go to Manabie website", { timeout: 60 * 1000 }, async function () {
    await this.openUrl('http://manabie.com/');
});

When("Bob click Careers", async function () {
    await this.page.click('text=Careers');
});

Then("See all job openings at Manabie", async function () {
    await this.page.click('text=View Openings');
    await this.page.waitForSelector('text=Our Openings');
});