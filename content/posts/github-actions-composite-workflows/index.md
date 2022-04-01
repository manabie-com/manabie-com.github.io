+++
date = "2022-04-01T12:00:00+07:00"
author = "anhpngt"
description = "Create a reusable workflow for GitHub Actions"
title = "Create a reusable workflow for GitHub Actions"
categories = ["DevSecOps", "CI/CD", "GitHub Actions", "GitHub"]
tags = ["github", "github action"]
slug = "create-a-reusable-workflow-for-github-actions"
+++

This blog post focuses on how we can create a reusable workflow for GitHub Actions.
The do's and don'ts will be elaborate along the way.

### **1. What are the available options?**

There are two major types of workflows that can be reused in Github Actions:

- [Custom actions](https://docs.github.com/en/actions/using-workflows/reusing-workflows)
- [Reusable workflows](https://docs.github.com/en/actions/creating-actions/about-custom-actions)

Reusable workflows are quite new. In this blog, we decide to use *custom actions*
(or more specifically, composite actions). The main reason is because we want
to run the workflow inside a job, along with other steps (to share caches, tokens, etc...).
Meanwhile, **reusable workflow must be run in a separate job**, so it is not desirable for our use-case.

### **2. Reusing an action**

In Manabie, we are using Google Kubernetes Engine (GKE) to host our services, and
Github Action for CI/CD. We have a lot of Github Action workflows (deployment, testing,
monitoring, etc...), and in every workflow, we have to keep copying-and-pasting the same
step to grant Github Action access to our GKE cluster.

It looks something like this:

```yaml
...
  - name: Authenticate to Google Cloud
    uses: google-github-actions/auth@v0
    with:
      credentials_json: ${{ env.SERVICE_ACCOUNT_KEY }}
      project_id: manabie-stag
      create_credentials_file: true

  - name: Setup Cloud SDK
    uses: google-github-actions/setup-gcloud@v0
    with:
      project_id: manabie-stag

  - name: Setup kubectl
    run: gcloud container clusters get-credentials manabie-stag-cluster --region asia-southeast1
...
```

As copying-and-pasting is troublesome and error prone (you know, DRY principle), 
in this blog, we explore how to create an action that can be reused for many workflows.

#### 2.1. Initializing a new repository

To run the workflow, we need a Github repository. In this blog, the
example repository will be `manabie-com/reusing-workflows`. Replace
the name of the repository with your own.

Clone the repository: 

```bash
git clone manabie-com/reusing-workflows
cd reusing-workflows
```

#### 2.2. Creating a new composite action

A composite action must have its own folder, preferrably inside `.github` directory.
We will create a new action called `setup-kubectl` in the following path.

```bash
touch -p .github/actions/setup-kubectl/action.yaml
```

Note that the file name must be `action.yaml` or `action.yml`. Other names will not work.

```yaml
# .github/actions/setup-kubectl/action.yaml
name: "setup-kubectl"
description: "Setup gcloud and kubenetes"
inputs:
  service_account_key:
    description: "Service account key"
    required: true
runs:
  using: "composite"
  steps:
    - name: Authenticate to Google Cloud
      uses: google-github-actions/auth@v0
      with:
        credentials_json: ${{ inputs.SERVICE_ACCOUNT_KEY }}
        project_id: manabie-stag
        create_credentials_file: true

    - name: Setup Cloud SDK
      uses: google-github-actions/setup-gcloud@v0
      with:
        project_id: manabie-stag

    - name: Setup kubectl
      shell: bash
      run: gcloud container clusters get-credentials manabie-stag-cluster --region asia-southeast1
```

A few things to note:
- We cannot use directly use Github secrets inside a composite action.
The only way to access secrets is to input it via `inputs` fields. In this case,
we ask the caller to input the secret key to `service_account_key` field.

#### 2.3. Using the composite action

In this example, we use `setup-kubectl` in an action name `gke` from **the same repository**.

```bash
touch -p .github/workflows/gke.yaml
```

```yaml
# .github/workflows/gke.yaml
name: gke
on:
  workflow_dispatch:
jobs:
  manabie-stag:
    runs-on: ubuntu-20.04
    steps:
      - name: Checkout source
        uses: actions/checkout@v3

      - name: Setup kubectl
        uses: ./.github/actions/setup-kubectl
        with:
          service_account_key: ${{ secrets.SERVICE_ACCOUNT_KEY }}

      - name: Run some kubectl commands
        run: |
          kubectl get namespaces
          kubectl get nodes
```

We also need to add the `SERVICE_ACCOUNT_KEY` to the Github secrets 
(in the Github repository `Settings` -> `Secrets` -> `Actions`).
The key can be retrieved from the GCP console, in the `IAM and admin` -> `Service accounts` section.
The key is a JSON file that looks like:

```json
{
  "type": "service_account",
  "project_id": "manabie-stag",
  "private_key_id": "...",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "manabie@manabie-stag.iam.gserviceaccount.com",
  "client_id": "...",
  "auth_uri": "https://...",
  "token_uri": "https://...",
  "auth_provider_x509_cert_url": "https://...",
  "client_x509_cert_url": "https://..."
}
```

Then, time to commit the changes.

```bash
git add .
git commit -m "Add setup-kubectl action"
git push
```

Then, go to the repository UI -> `Actions` tab -> Choose `gke` action -> Choose `Run workflow`.

Or, using the [Github CLI](https://cli.github.com/), the command is:

```bash
gh workflow run gke
```

If everything goes right, you should see the workflow accessing
to GKE cluster and printing out the k8s namespaces.

#### 2.4. Improving the action

In the previous example, we used Manabie Staging environment.
What if we have other environments, in a different GKE cluster of 
a different GCP project.

**Let's try to make the action reusable for both Staging and Production.**

There are a few things we would need to consider:
1. We need to allow the caller to choose the environment
2. All the parameters between Staging and Production are different
(GCP project ID, GKE cluster name, region, etc...)

##### **Allowing caller to choose the environment**

It is quite simple, since Github Action has built-in support for it,
using the `inputs` field.

```yaml
# .github/actions/setup-kubectl/action.yaml
inputs:
  environment:
    description: "Environment (stag or prod)"
    required: true
```

##### **Handling different configurations for different environments**

Assuming the environments have already been set up, the project ID,
cluster name, etc... will not changed for each environment.

Thus, we can implement a simple switch case to tell the action which
parameters to use:

```yaml
# .github/actions/setup-kubectl/action.yaml
  steps:
    - name: Get configuration
      uses: actions/github-script@v5
      with:
        script: |
          const env = '${{ inputs.environment }};
          switch (env) {
            case 'stag':
              core.exportVariable('PROJECT_ID', 'manabie-stag');
              core.exportVariable('CLUSTER', 'manabie-stag-cluster');
              core.exportVariable('REGION', 'asia-southeast1');
            case 'prod':
              core.exportVariable('PROJECT_ID', 'manabie-prod');
              core.exportVariable('CLUSTER', 'manabie-prod-cluster');
              core.exportVariable('REGION', 'asia-northeast1');
          }

    - name: Authenticate to Google Cloud
      uses: google-github-actions/auth@v0
      with:
        credentials_json: ${{ inputs.SERVICE_ACCOUNT_KEY }}
        project_id: ${{ env.PROJECT_ID }}
        create_credentials_file: true
```

`core.exportVariable` exports the value to the environment variables, so that
the next steps can use them (with `${{ env.PROJECT_ID }}`)

##### **Handling different secrets for different environments**

The only secret we've been using now is the service account key (`${{ secrets.SERVICE_ACCOUNT_KEY }}`
in the previous step).

But since we have an additional environment, the number of keys is increased to 2, namely:
- `${{ secrets.STAG_SERVICE_ACCOUNT_KEY }}`
- `${{ secrets.PROD_SERVICE_ACCOUNT_KEY }}`

Since Github's composite actions cannot access secrets directly,
we need to have `gke` workflow input them to the composite action `setup-kubectl`:

```yaml
# .github/workflows/gke.yaml
      - name: Setup kubectl
        uses: ./.github/actions/setup-kubectl
        with:
          stag_service_account_key: ${{ secrets.STAG_SERVICE_ACCOUNT_KEY }}
          prod_service_account_key: ${{ secrets.PROD_SERVICE_ACCOUNT_KEY }}

# .github/actions/setup-kubectl/action.yaml
name: "setup-kubectl"
description: "Setup gcloud and kubenetes"
inputs:
  stag_service_account_key:
    description: "Manabie Staging service account key"
    required: true
  prod_service_account_key:
    description: "Manabie Production service account key"
    required: true
```

Then, again, inside the `setup-kubectl`, we can choose which secrets to use
with a simple script:

```yaml
# .github/actions/setup-kubectl/action.yaml
  steps:
    - name: Get configuration
      uses: actions/github-script@v5
      with:
        script: |
          const env = '${{ inputs.environment }};
          switch (env) {
            case 'stag':
              core.exportVariable('SERVICE_ACCOUNT_KEY', '${{ inputs.stag_service_account_key }}');
            case 'prod':
              core.exportVariable('SERVICE_ACCOUNT_KEY', '${{ inputs.prod_service_account_key }}');
          }
```

**However**, if you run this javascript script in Github Action, you'll encounter an error:

```sh
Error: Unhandled error: SyntaxError: Invalid or unexpected token
```

Remember that the service account key is a JSON. There are a lof of special characters in the key.
We have not escape those characters yet, so we cannot directly input it into the javascript.

Luckily, the action [google-github-actions/auth@v0](https://github.com/google-github-actions/auth)
allows us to input the key in a base64 format instead. After encoded, the key will look like this:

```sh
echo '{ <service_account_key> }' | base64
ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAibWFuYWJpZS1z
dGFnIiwKICAicHJpdmF0ZV9rZXlfaWQiOiAiLi4uIiwKICAicHJpdmF0ZV9rZXkiOiAiLS0tLS1C
RUdJTiBQUklWQVRFIEtFWS0tLS0tXG4uLi5cbi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS1cbiIs
CiAgImNsaWVudF9lbWFpbCI6ICJtYW5hYmllQG1hbmFiaWUtc3RhZy5pYW0uZ3NlcnZpY2VhY2Nv
dW50LmNvbSIsCiAgImNsaWVudF9pZCI6ICIuLi4iLAogICJhdXRoX3VyaSI6ICJodHRwczovLy4u
LiIsCiAgInRva2VuX3VyaSI6ICJodHRwczovLy4uLiIsCiAgImF1dGhfcHJvdmlkZXJfeDUwOV9j
ZXJ0X3VybCI6ICJodHRwczovLy4uLiIsCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBz
Oi8vLi4uIgp9Cg==
```

No more special characters! Let's use these base64 encoded key instead for our secrets.

##### **Combining everything**

We combine all the previous ideas:
- Adding inputs for workflow
- Adding script to handle different configs for different environments
- Base64 encode the service account key

With that in mind, we update the `setup-kubectl` action to:

```yaml
# .github/actions/setup-kubectl/action.yaml
name: "setup-kubectl"
description: "Setup gcloud and kubenetes"
inputs:
  environment:
    description: "Environment (stag or prod)"
    required: true
  stag_service_account_key:
    description: "Manabie Staging service account key"
    required: true
  prod_service_account_key:
    description: "Manabie Production service account key"
    required: true
runs:
  using: "composite"
  steps:
    - name: Get configuration
      uses: actions/github-script@v5
      with:
        script: |
          const env = '${{ inputs.environment }};
          switch (env) {
            case 'stag':
              core.exportVariable('PROJECT_ID', 'manabie-stag');
              core.exportVariable('CLUSTER', 'manabie-stag-cluster');
              core.exportVariable('REGION', 'asia-southeast1');
              core.exportVariable('SERVICE_ACCOUNT_KEY', '${{ inputs.stag_service_account_key }}');
            case 'prod':
              core.exportVariable('PROJECT_ID', 'manabie-prod');
              core.exportVariable('CLUSTER', 'manabie-prod-cluster');
              core.exportVariable('REGION', 'asia-northeast1');
              core.exportVariable('SERVICE_ACCOUNT_KEY', '${{ inputs.prod_service_account_key }}');
          }

    - name: Authenticate to Google Cloud
      uses: google-github-actions/auth@v0
      with:
        credentials_json: ${{ env.SERVICE_ACCOUNT_KEY }}
        project_id: ${{ env.PROJECT_ID }}
        create_credentials_file: true

    - name: Setup Cloud SDK
      uses: google-github-actions/setup-gcloud@v0
      with:
        project_id: ${{ env.PROJECT_ID }}

    - name: Setup kubectl
      shell: bash
      run: gcloud container clusters get-credentials ${{ env.CLUSTER }} --region ${{ env.REGION }}
```

and in the caller workflow `gke`:

```yaml
# .github/workflows/gke.yaml
name: gke
on:
  workflow_dispatch:
    inputs:
      environment:
        type: choice
        description: "Environment"
        require: true
        options:
          - stag
          - prod
jobs:
  manabie-stag:
    runs-on: ubuntu-20.04
    steps:
      - name: Checkout source
        uses: actions/checkout@v3

      - name: Setup kubectl
        uses: ./.github/actions/setup-kubectl
        with:
          environment: ${{ inputs.environment }}
          stag_service_account_key: ${{ secrets.STAG_SERVICE_ACCOUNT_KEY }}
          prod_service_account_key: ${{ secrets.PROD_SERVICE_ACCOUNT_KEY }}

      - name: Run some kubectl commands
        run: |
          kubectl get namespaces
          kubectl get nodes
```

Time to commit and push the update

```bash
git add .
git commit -m "Allow to setup kubectl for different environments"
```

Finally, we can trigger the workflow `gke` using the Github UI. If using the CLI, the command would be:

```bash
gh workflow run gke -f environment=stag
# or
gh workflow run gke -f environment=prod
```


### **3. Conclusion**

In this blog, we have:
- Discussed how to reuse a Github Action workflow using a composite action
- Extended and improved the composite action to include more features
- Run the created workflow on Github Action
