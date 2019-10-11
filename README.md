[![Build Status](https://travis-ci.com/Preskton/terraform-provider-twilio.svg?branch=master)](https://travis-ci.com/Preskton/terraform-provider-twilio) ![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/Preskton/terraform-provider-twilio?include_prereleases)

# Twilio Terraform Provider

The goal of this Terraform provider plugin is to make manging your Twilio account easier.

Current features:

- `twilio_phone_number`
  - Search
    - US & International
    - Number prefix (or place * wherever you'd like!)
  - Purchase (`terraform apply`)
  - Delete/release (`terraform destroy`)
- `twilio_subaccount`
  - Create
  - Update
  - Delete

More coming soon.

## Getting Started

1. Start a trial account at twilio.com (if you don't have one already). Use the Console Dashboard to take note of your Account SID (a long string starts with `AC` and looks like a GUID) and Auth Token (also a long GUID-like string, hidden under the `View` link).
2. Download the latest release of the provider and place in your `~/.terraform.d/plugins` directory.
3. Use the example below, replacing `account_sid` and `auth_token` with the appropriate values.
4. `terraform apply` Note: this will cost you REAL MONEY (or at the very least trial credits).

## Example

Note: running and applying the below could cost you REAL MONEY! Please use this tool wisely!

```hcl
provider "twilio" {
    account_sid = "<your account sid here>"
    auth_token = "<your auth token here>"
}

resource "twilio_subaccount" "woomy" {
    friendly_name = "Woomy Subaccount #1"
}

resource "twilio_phone_number" "us_dallas_tx" {
    country_code = "US"
    search = "972"
    friendly_name = "Howdy from TX"
}

resource "twilio_phone_number" "japan_somewhere" {
    country_code = "JP"
    search = "503*"
    friendly_name = "日本"
}
```
