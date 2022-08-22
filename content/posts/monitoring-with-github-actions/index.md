+++
date = "2022-03-30T17:10:49+07:00"
author = "ds0nt"
description = "Using Github Actions to run some of your monitoring tasks on a schedule."
title = "Using Github Actions for Daily Monitoring Tasks"
categories = ["DevSecOps", "Automation", "CI"]
tags = ["DevOps", "Github Actions", "Monitoring", "CRON", "Schedule"]
slug = "github-actions-for-monitoring-tasks"
+++

# Using Github Actions for Daily Monitoring Tasks

Github actions has a CRON schedule trigger we can use run scripts on a custom schedule. It's pretty cool that Github Actions can run jobs on a CRON schedule, and we can use this ability to cover some of our dev-ops needs.

For example, we used this to check if any of our domain certificates were close to their expiry. Each day, a github action  runs at 10:45, checks if our certificates are dangerously close to expiring, and posts to our `#monitoring` Slack if any of them are.


### The Workflow

This is a simple monitoring flow really; Run Check -> Alert if Problem.

To do this, were going to use three Github actions.

1. an action to checkout our code
2. an action to execute a script
3. and an action to send slack message based on our condition.

```yaml
# github/workflows/check_weather.yml`

name: check_weather
on:
  schedule: # run at 010:45 UTC daily
    - cron: '45 10 * * *'
  workflow_dispatch: # runnable manually.

jobs:
  check-weather:
    runs-on: ubuntu-latest
    steps:

      # checkout code
      - uses: actions/checkout@v2
     
      # run script
      - name: Check Weather
        shell: bash
        run: ./check-weather.sh
      
      # post results
      - if: failure()
        name: Slack Notification
        uses: rtCamp/action-slack-notify@v2
        env:
          SLACK_WEBHOOK: https://hooks.slack.com/services/*****/*****/*****
          SLACK_CHANNEL: '#monitoring'
          SLACK_USERNAME: "weather-checker"
          SLACK_TITLE: Bad Weather Warning
```

This job runs a `./check-weather.sh` script every day at 10:45AM UTC and sends a message on slack if we want it to.

We're using this action [rtCamp/action-slack-notify](https://github.com/rtCamp/action-slack-notify) to send the slack notification. We set most of the variables for the message already, but we still want to set the `SLACK_MESSAGE` variable from our weather-checking script. It's not hard to set env vars from prior steps, check it out [in the docs](https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-environment-variable). 


### Writing the Check Script

We need to perform some check (like the weather), and exit with status 0 if it's ok, or exit with a non-zero status if theres a problem.

```bash
# check-weather.sh

#!/bin/bash

# fetch the temperature
TEMPERATURE=$(curl 'wttr.in/Ho_Chi_Minh?format=3' | cut -d' ' -f4 | grep -oE '[+-][0-9]+')

# test if the temperature is ok
if [[ $TEMPERATURE -le 20 ]]; then
	MESSAGE="Dress warm, it's a chilly $TEMPERATURE outside"
elif [[ $TEMPERATURE -ge 30 ]]; then
	MESSAGE="Its going to be hot. The temperature is $TEMPERATURE"
fi


# if it was ok, we can exit 0
[[ -n $MESSAGE ]] || exit 0

# otherwise, we set the SLACK_MESSAGE env var for the next step, and exit 1.
echo 'SLACK_MESSAGE<<EOF' >> $GITHUB_ENV
echo -e "$MESSAGE" >> $GITHUB_ENV
echo 'EOF' >> $GITHUB_ENV
exit 1
```
Once your ready to run your workflow, make sure it's merged into your main github branch, for us it's `develop`.

Then it will run on it's CRON schedule, and you can also dispatch it manually from the actions page to run the check immediately.

### Conclusion

Github Actions is capable of running jobs on a timer. It's native to your github repository so you don't need to muck about with maintaining a deployment of it in the cloud, or writing infrastructure code.

Next time you ask yourself how to deploy one of these simple scripted checks, consider github actions.