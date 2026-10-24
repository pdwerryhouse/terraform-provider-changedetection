<!--
 Copyright (C) 2026 Paul Dwerryhouse <paul@dwerryhouse.com.au>
 
 This file is part of terraform-provider-changedetection.
 
 terraform-provider-changedetection is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 
 terraform-provider-changedetection is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 
 You should have received a copy of the GNU General Public License
 along with terraform-provider-changedetection.  If not, see <https://www.gnu.org/licenses/>.
-->

# Terraform Provider for changedetection.io

## Repository

The main repository for this code is at [codeberg](https://codeberg.org/pdwerryhouse/terraform-provider-changedetection).

The github repository is only to allow access to the Hashicorp registry.

## Introduction

This is a Terraform provider for [changedetection.io](https://changedetection.io/), the website change detection system.

## Installation

```terraform
terraform {
  required_providers {
    changedetection = {
      source  = "pdwerryhouse/changedetection"
      version = "~> 0.0.1"
    }
  }
}
```

## Usage

```terraform

provider "changedetection" {
  endpoint = "http://localhost:5000/api/v1"
  apikey   = "some_api_key"
}

resource "changedetection_tag" "things" {
  title = "things"
  tag_colour = "#112233"
}

resource "changedetection_watch" "something" {
  title = "Something"
  url = "https://example.com/item"
  method = "GET"
  notification_muted = true
  processor = "restock_diff"
  tags = [ changedetection_tag.things.id ]
  webdriver_delay = 3
  paused = true
  ignore_status_codes = true
}
```

## Documentation

The documentation is available at the [Terraform Registry](https://registry.terraform.io/providers/pdwerryhouse/changedetection/latest/docs).

## Development

terraform-provider-changedetection uses [changedetection-client-go](https://codeberg.org/pdwerryhouse/changedetection-client-go), a Go client library for the [changedetection.io API](https://changedetection.io/docs/api_v1/index.html).

## Licence

GNU GENERAL PUBLIC [LICENSE](LICENSE) V3 or later
