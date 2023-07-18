---
subcategory: "SES (Simple Email)"
layout: "aws"
page_title: "AWS: aws_ses_receipt_filter"
description: |-
  Provides an SES receipt filter
---

# Resource: aws_ses_receipt_filter

Provides an SES receipt filter resource

## Example Usage

```terraform
resource "aws_ses_receipt_filter" "filter" {
  name   = "block-spammer"
  cidr   = "10.10.10.10"
  policy = "Block"
}
```

## Argument Reference

This resource supports the following arguments:

* `name` - (Required) The name of the filter
* `cidr` - (Required) The IP address or address range to filter, in CIDR notation
* `policy` - (Required) Block or Allow

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - The SES receipt filter name.
* `arn` - The SES receipt filter ARN.

## Import

SES Receipt Filter can be imported using their `name`, e.g.,

```
$ terraform import aws_ses_receipt_filter.test some-filter
```
