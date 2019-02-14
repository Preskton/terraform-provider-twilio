[![Build Status](https://travis-ci.com/Preskton/terraform-provider-twilio.svg?branch=master)](https://travis-ci.com/Preskton/terraform-provider-twilio)

# Twilio Terraform Provider

The goal of this Terraform provider plugin is to make manging your Twilio account easier.

Current features:

- `twilio_phone_number`
  - Search
    - US & International
    - Number prefix (or place * wherever you'd like!)
  - Purchase (`terraform apply`)
  - Delete/release (`terraform destroy`)

More coming soon.

## Example

```hcl
provider "twilio" {
    account_sid = "<your account sid here>"
    auth_token = "<your auth token here>"
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