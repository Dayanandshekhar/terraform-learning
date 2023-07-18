---
subcategory: "SageMaker"
layout: "aws"
page_title: "AWS: aws_sagemaker_notebook_instance_lifecycle_configuration"
description: |-
  Provides a lifecycle configuration for SageMaker Notebook Instances.
---

# Resource: aws_sagemaker_notebook_instance_lifecycle_configuration

Provides a lifecycle configuration for SageMaker Notebook Instances.

## Example Usage

Usage:

```terraform
resource "aws_sagemaker_notebook_instance_lifecycle_configuration" "lc" {
  name      = "foo"
  on_create = base64encode("echo foo")
  on_start  = base64encode("echo bar")
}
```

## Argument Reference

This resource supports the following arguments:

* `name` - (Optional) The name of the lifecycle configuration (must be unique). If omitted, Terraform will assign a random, unique name.
* `on_create` - (Optional) A shell script (base64-encoded) that runs only once when the SageMaker Notebook Instance is created.
* `on_start` - (Optional) A shell script (base64-encoded) that runs every time the SageMaker Notebook Instance is started including the time it's created.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - The Amazon Resource Name (ARN) assigned by AWS to this lifecycle configuration.

## Import

Models can be imported using the `name`, e.g.,

```
$ terraform import aws_sagemaker_notebook_instance_lifecycle_configuration.lc foo
```
