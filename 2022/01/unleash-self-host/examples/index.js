'use strict';

const unleash = require('unleash-server');

let options = {};

async function startUnleash() {
    await unleash.start(options);
}

startUnleash();