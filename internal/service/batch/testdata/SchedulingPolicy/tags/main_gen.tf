# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "aws_batch_scheduling_policy" "test" {
  name = var.rName

  fair_share_policy {
    compute_reservation = 0
    share_decay_seconds = 0
  }

  tags = var.tags
}


variable "rName" {
  type     = string
  nullable = false
}

variable "tags" {
  type     = map(string)
  nullable = false
}


