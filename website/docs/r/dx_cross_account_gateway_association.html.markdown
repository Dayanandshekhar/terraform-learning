---
layout: "aws"
page_title: "AWS: aws_dx_cross_account_gateway_association"
sidebar_current: "docs-aws-resource-dx-cross-account-gateway-association"
description: |-
  Associates a Direct Connect Gateway with a VGW or transit gateway in another AWS Account.
---

# Resource: aws_dx_cross_account_gateway_association

Associates a Direct Connect Gateway with a VGW or transit gateway in another AWS Account by accepting a Direct Connect Gateway Association Proposal.
For single account associations, see the [`aws_dx_gateway_association` resource](/docs/providers/aws/r/dx_gateway_association.html).


To create a cross-account association, create an [`aws_dx_gateway_association_proposal` resource](/docs/providers/aws/r/dx_gateway_association_proposal.html)
in the AWS account that owns the VGW or transit gateway and then accept the proposal in the AWS account that owns the Direct Connect Gateway
by creating an `aws_dx_cross_account_gateway_association` resource.

## Example Usage

```hcl
provider "aws" {
  # Creator's credentials.
}

provider "aws" {
  alias = "accepter"

  # Accepter's credentials.
}

# Creator's side of the proposal.
data "aws_caller_identity" "creator" {}

resource "aws_vpc" "example" {
  cidr_block = "10.255.255.0/28"
}

resource "aws_vpn_gateway" "example" {
  vpc_id = "${aws_vpc.example.id}"
}

resource "aws_dx_gateway_association_proposal" "example" {
  dx_gateway_id               = "${aws_dx_gateway.example.id}"
  dx_gateway_owner_account_id = "${aws_dx_gateway.example.owner_account_id}"
  associated_gateway_id       = "${aws_vpn_gateway.example.id}"
}

# Accepter's side of the proposal.
resource "aws_dx_gateway" "example" {
  provider = "aws.accepter"

  name            = "example"
  amazon_side_asn = "64512"
}

resource "aws_dx_cross_account_gateway_association" "example" {
  provider = "aws.accepter"

  proposal_id                         = "${aws_dx_gateway_association_proposal.example.id}"
  dx_gateway_id                       = "${aws_dx_gateway.example.id}"
  associated_gateway_owner_account_id = "${data.aws_caller_identity.creator.account_id}"
}
```

## Argument Reference

The following arguments are supported:

* `associated_gateway_owner_account_id` - (Required) The ID of the AWS account that owns the VGW or transit with which to associate the Direct Connect gateway.
* `dx_gateway_id` - (Required) The ID of the Direct Connect gateway.
* `proposal_id` - (Required) The ID of the Direct Connect gateway association proposal.
* `allowed_prefixes` - (Optional) VPC prefixes (CIDRs) to advertise to the Direct Connect gateway. Defaults to the CIDR block of the VPC associated with the VGW. To enable drift detection, must be configured.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Direct Connect gateway association resource.
* `associated_gateway_id` - The ID of the VGW or transit gateway with which the gateway is associated.
* `associated_gateway_type` - The type of the associated gateway, `transitGateway` or `virtualPrivateGateway`.
* `dx_gateway_association_id` - The ID of the Direct Connect gateway association.
* `dx_gateway_owner_account_id` - The ID of the AWS account that owns the Direct Connect gateway.

## Timeouts

`aws_dx_cross_account_gateway_association` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `15 minutes`) Used for creating the association
- `update` - (Default `10 minutes`) Used for updating the association
- `delete` - (Default `15 minutes`) Used for destroying the association
