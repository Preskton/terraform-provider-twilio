[![Build Status](https://travis-ci.com/Preskton/terraform-provider-twilio.svg?branch=master)](https://travis-ci.com/Preskton/terraform-provider-twilio) [![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/Preskton/terraform-provider-twilio?include_prereleases)](https://github.com/Preskton/terraform-provider-twilio/releases/latest)

# Twilio Terraform Provider

The goal of this Terraform provider plugin is to make managing your Twilio account easier.

Current features:

- Compatible with Terraform `v0.12.10`
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
