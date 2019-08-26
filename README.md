[![Build Status](https://travis-ci.com/Preskton/terraform-provider-twilio.svg?branch=master)](https://travis-ci.com/Preskton/terraform-provider-twilio)

# Twilio Terraform Provider

The goal of this Terraform provider plugin is to make manging your Twilio account easier.

Current features:

- `twilio_subaccount`
  - Create
  - Update
  - Delete
- `twilio_application`
  - Create
  - Update
  - Delete

More coming soon.

## Build
Run:
```
make plugin
```
to build and move the plugin to `~/.terraform.d/plugins` which is where terraform will look for all 3rd party plugins

## Getting Started

1. Start a trial account at twilio.com (if you don't have one already). Use the Console Dashboard to take note of your Account SID (a long string starts with `AC` and looks like a GUID) and Auth Token (also a long GUID-like string, hidden under the `View` link).
2. Download the latest release of the provider and place in your `~/.terraform.d/plugins` directory.
3. Use the example below, replacing `account_sid` and `auth_token` with the appropriate values.
4. `terraform apply` Note: this will cost you REAL MONEY (or at the very least trial credits).

## Debugging

Adding the following lines to your bash profile will enable additional logging.
```
export TF_LOG=TRACE
export TF_LOG_PATH=./terraform.log
export DEBUG_HTTP_TRAFFIC=true
```
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

resource "twilio_application" "new_twiml_app" {
    friendly_name = "My new TwiML application"
}
```
