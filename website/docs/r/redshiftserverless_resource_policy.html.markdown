---
subcategory: "Redshift Serverless"
layout: "aws"
page_title: "AWS: aws_redshiftserverless_resource_policy"
description: |-
  Provides a Redshift Serverless Resource Policy resource.
---

# Resource: aws_redshiftserverless_resource_policy

Creates a new Amazon Redshift Serverless Resource Policy.

## Example Usage

```terraform
resource "aws_redshiftserverless_resource_policy" "example" {
  resource_arn = aws_redshiftserverless_snapshot.example.arn
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        AWS = ["12345678901"]
      }
      Action = [
        "redshift-serverless:RestoreFromSnapshot",
      ]
      Sid = ""
    }]
  })
}
```

## Argument Reference

This resource supports the following arguments:

* `resource_arn` - (Required) The Amazon Resource Name (ARN) of the account to create or update a resource policy for.
* `policy` - (Required) The policy to create or update. For example, the following policy grants a user authorization to restore a snapshot.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - The Amazon Resource Name (ARN) of the account to create or update a resource policy for.

## Import

Redshift Serverless Resource Policies can be imported using the `resource_arn`, e.g.,

```
$ terraform import aws_redshiftserverless_resource_policy.example example
```
