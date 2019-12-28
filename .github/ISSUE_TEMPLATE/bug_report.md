---
name: Bug report
about: Create a report to help us improve
title: ''
labels: bug
assignees: ''

---

### Description

A clear and concise description of what the bug is.

### Steps to reproduce

Steps to reproduce the behavior:

1. Using the sample configuration, run `terraform plan`
2. Run `terraform apply`

#### Sample Terraform Configuration

Please include a sample `.tf` that causes the issue you're encountering. PLEASE DO NOT INCLUDE ACCOUNT SIDs, SECRETS, OR OTHER SENSITIVE  INFORMATION!

```hcl
resource "twilio_phone_number" "mine" {
   ...
}
```

### Expected behavior

A clear and concise description of what you expected to happen.

### Environment:

- OS: [e.g. Windows 10]
- Terraform Version: [e.g. 0.12]
- `terraform-provider-twilio` version: [e.g 0.1.4]
 
### Additional context

Add any other context about the problem here.
