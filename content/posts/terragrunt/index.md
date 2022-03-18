+++
date = "2021-12-15T10:00:00+07:00"
author = "nguyenhoaibao"
description = "Terraform modules and Terragrunt to connect multiple modules together"
title = "Terraform modules and Terragrunt to connect multiple modules together"
categories = ["Infrastructure"]
tags = ["Terraform"]
slug = "terraform-modules-and-terragrunt" 
+++

Suppose we want to provision infrastructure that including:
  - A VPC network.
  - A Postgresql instance.
  - A GKE cluster.
  - A module to create databases, service accounts and RBAC definition to let the app run.

And we also have to deploy that to multiple environments: staging, uat and production.

One way to implement that is putting all of the infrastructure code into a single module,
alongside with all variables for each environment. That works but some best practices out there
state that [large modules considered harmful](https://blog.gruntwork.io/5-lessons-learned-from-writing-over-300-000-lines-of-infrastructure-code-36ba7fadeac1)
and [Small, focused workspaces make Terraform runs fast, limit the blast radius, and enable easier work separation by teams](https://github.com/mitchellh/terraform-provider-multispace#why)
so we will try to split that into multiple modules to make the code runs faster and easier to maintain.

#### Create Terraform modules

First we will have a module to provision a VPC network:
```
modules
└── vpc
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
```

Inside that we just need to create VPC and its subnetworks resources:
```
// in vpc/main.tf

resource "google_compute_network" "vpc" {
  ...
}

resource "google_compute_subnetwork" "subnetworks" {
  ...
}
```

We also need to output the VPC network name, so we can use it in other modules later:
```
// in vpc/outputs.tf

output "network_name" {
  value = ...
}

output "network_self_link" {
  value = ...
}
```

Next we need to create another module for the Cloud SQL instance:
```
modules
└── vpc
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
└── postgresql
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
```

First we need to define the VPC network that this instance will be deployed to:
```
// in postgresql/variables.tf

variable "network_name" {
  ...
}
```

Then we will provision a Postgresql instance with the network name comes from the variable:
```
// in postgresql/main.tf

resource "google_sql_database_instance" "postgresql" {
  ...
  settings {
    ip_configuration {
      private_network = var.network_name
    }
  }
  ...
}
```

Again, we need to output the Postgresql instance name so we can use it in other modules later:
```
// in postgresql/outputs.tf

output "postgresql_instance" {
  value = ...
}
```

Next is the module for GKE cluster, it's similar to the Postgresql module so we won't set it here, but you get the idea.

And finally is the module to create necessary databases and service accounts for our app, so we name it the `apps` module:
```
modules
└── vpc
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
└── postgresql
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
└── gke
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
└── apps
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
```

This module need the Postgresql instance to create databases and GKE cluster name to create RBAC definition, so we will
define those variables for that:
```
// in apps/variables.tf

// the Postgresql instance name
variable "postgresql_instance" {
  ... 
}

// the GKE instance
variable "gke_cluster" {
  ...
}
```

Then we will create databases to that `postgresql_instance`:
```
// in apps/main.tf

resource "google_sql_database" "databases" {
  ...
  instance = var.postgresql_instance
}

```

And setup the `kubernetes` provider to create the RBAC definition:
```
provider "kubernetes" {
  host                   = "https://${var.gke_cluster}"
  ...
}
```

And we're done for the modules. Next is using those modules to provision our infrastructure.

#### Provision infrastructure

First we need to create another `live` directory to contain our infrastructure configuration for each environment.

We will start with the `stag` environment first:
```
live
├── stag
│   └── vpc
│       └── terragrunt.hcl
modules
└── vpc
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
└── postgresql
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
└── gke
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
└── apps
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
```

In `stag/vpc/terragrunt.hcl` we just need to set the module source point to our `vpc` module:
```
terraform {
  source = "../../../modules/vpc"
}
```
then run `terragrunt apply` to create a new VPC network for `stag` environment.

Next is creating the Postgresql instance:
```
live
├── stag
│   ├── vpc
│   │   └── terragrunt.hcl
│   └── postgresql
│       └── terragrunt.hcl
modules
└── ...
```

Here is the content inside `stag/postgresql/terragrunt.hcl`:
```
terraform {
  source = "../../../modules/postgresql"
}

dependency "vpc" {
  config_path = "../vpc"
}

inputs = {
  network_name = dependency.vpc.outputs.network_name
}
```

Note that in this configuration we set the `dependency` to the `vpc` module, and read the network name
from the outputs of that dependency.

The configuration for our GKE cluster is similar, so I won't set it here.

And finally is configuration for the `apps` module:
```
terraform {
  source = "../../../modules/postgresql"
}

dependency "postgresql" {
  config_path = "../postgresql"
}

dependency "gke" {
  config_path = "../gke"
}

inputs = {
  postgresql_instance = dependency.postgresql.outputs.postgresql_instance
  gke_cluster         = dependency.gke.outputs.gke_cluster
}
```
note that this time we set two dependencies: the `postgresql` module and the `gke` module.

And finally we will have this structure for our infrastructure:
```
live
├── stag
│   ├── apps
│   │   └── terragrunt.hcl
│   ├── postgresql
│   │   └── terragrunt.hcl
│   ├── gke
│   │   └── terragrunt.hcl
│   └── vpc
│       └── terragrunt.hcl
├── uat
│   ├── apps
│   │   └── terragrunt.hcl
│   ├── postgresql
│   │   └── terragrunt.hcl
│   ├── gke
│   │   └── terragrunt.hcl
│   └── vpc
│       └── terragrunt.hcl
├── prod
│   ├── apps
│   │   └── terragrunt.hcl
│   ├── postgresql
│   │   └── terragrunt.hcl
│   ├── gke
│   │   └── terragrunt.hcl
│   └── vpc
│       └── terragrunt.hcl
modules
└── vpc
│   ├── main.tf
│   ├── outputs.tf
│   └── variables.tf
└── postgresql
│   ├── main.tf
│   ├── outputs.tf
│   └── variables.tf
└── gke
│   ├── main.tf
│   ├── outputs.tf
│   └── variables.tf
└── apps
    ├── main.tf
    ├── outputs.tf
    └── variables.tf
```

#### Summary

To summary, the benefit of spliting infrastructure code into multiple modules, and wire those modules
together is it helps to make the module code smaller, runs faster and easier to maintain. For later if we
want to provision another Postgresql or GKE instance, we just need to create its directory in each
environment directory and reuse outputs from dependency module, without worrying about the dependency configuration.
