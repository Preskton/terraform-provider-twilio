# 🎉 Twilio has launched an [OFFICIAL Terraform provider](https://github.com/twilio/terraform-provider-twilio). I **strongly** encourage you to check it out and, when possible, use it instead, as this version of the provider will no longer be maintained. 🎉

[![Build Status](https://github.com/Preskton/terraform-provider-twilio/workflows/build/badge.svg)](https://github.com/Preskton/terraform-provider-twilio/actions?query=workflow%3Abuild) [![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/Preskton/terraform-provider-twilio?include_prereleases)](https://github.com/Preskton/terraform-provider-twilio/releases/latest) [![Terraform Registry](https://img.shields.io/badge/registry-twilio-green?logo=terraform&style=flat)](https://registry.terraform.io/providers/Preskton/twilio/latest)

# Twilio Terraform Provider

The goal of this Terraform provider plugin is to make managing your Twilio account easier.

Current features:

- Compatible with Terraform `v0.12+`
- `twilio_phone_number`
  - Search
    - Country code
    - Area Code
    - Number prefix (or place * wherever you'd like!)
  - Create/Purchase
  - Update
  - Delete/Release
- `twilio_subaccount`
  - Create
  - Update
  - Delete
- `twilio_api_key`
  - Create
  - Delete

More coming eventually!

## Getting Started

### Installing the provider

#### Terraform .13+

Starting in Terraform versions `.13` and greater, you can reference this provider via the [Terraform Registry](https://registry.terraform.io/providers/Preskton/twilio/latest). You don't have to download a thing - Terraform does the work for you. Simply include this as part of your `.tf` file:

```hcl
terraform {
  required_providers {
    twilio = {
      source = "Preskton/twilio"
      version = "0.1.5"
    }
  }
}
```

#### Prior versions of Terraform

Download the [latest release of the provider](https://github.com/Preskton/terraform-provider-twilio/releases/latest) for your operating system/processor architecture, unzip, and place in your `~/.terraform.d/plugins` directory.

### The basics

1. Start a trial account at twilio.com (if you don't have one already). Use the Console Dashboard to take note of your Account SID (a long string starts with `AC` and looks like a GUID) and Auth Token (also a long GUID-like string, hidden under the `View` link).
2. Use the example below, replacing `account_sid` and `auth_token` with the appropriate values.
3. `terraform apply` Note: this will cost you REAL MONEY (or at the very least trial credits).

## Example

Note: running and applying the below could cost you REAL MONEY! Please use this tool wisely!

~> **Important:** You should never check the values for your `account_sid` and `auth_token` into source control, as this would allow others to modify your Twilio account. Instead, source these variables from the environment using [Terraform variables](https://www.terraform.io/docs/configuration/variables.html) sourced from envrionment variables or passed as arguments to your `plan`/`apply`.

```hcl
terraform {
  required_providers {
    twilio = {
      source = "Preskton/twilio"
      version = "0.1.5"
    }
  }
}

provider "twilio" {
    account_sid = "<your account sid here>"
    auth_token = "<your auth token here>"
}

resource "twilio_subaccount" "woomy" {
    friendly_name = "Woomy Subaccount #1"
}

resource "twilio_api_key" "woomy" {
    friendly_name = "Woomy Key #1"
}

resource "twilio_phone_number" "area_code_test" {
    // Find a number
    country_code = "US"
    area_code = "972"
    friendly_name = "terraform-provider-twilio area code test number"

    // Configure your number

    address_sid = "ADXXXX"          // Certain countries may require a validated address!
    identity_sid = "IDXXXX"         // Certain countries may require a validated identity!
    trunk_sid = "XXXXX"

    voice {
        primary_url = "https://genoq.com/handlers/voice-primary"
        primary_http_method = "POST"
        fallback_url = "https://genoq.com/handlers/voice-fallback"
        fallback_http_method = "GET"
        caller_id_enabled = "true"
        receive_mode = "voice"
    }

    sms {
        primary_url = "https://genoq.com/handlers/sms-primary"
        primary_http_method = "POST"
        fallback_url = "https://genoq.com/handlers/sms-fallback"
        fallback_http_method = "GeT"
    }

    status_callback {
        url = "https://genoq.com/handlers/status-callback"
        http_method = "GET"
    }

    // Note: Emergency calling requires a validated address
    emergency {
        address_sid = "ADXXXXX"
        status = "active"
    }
}

resource "twilio_phone_number" "search_test" {
    country_code = "US"
    search = "972*"
    friendly_name = "terraform-provider-twilio by-search test number"

    sms {
        primary_url = "https://genoq.com/handlers/sms-primary"
        primary_http_method = "POST"
        fallback_url = "https://genoq.com/handlers/sms-fallback"
        fallback_http_method = "GET"
    }

    voice {
        receive_mode = "fax"
        application_sid = "APXXXXX"
    }
}
```

## Disclaimer

This is NOT an official Twilio project and is maintained in [my](https://www.github.com/Preskton) free time.
