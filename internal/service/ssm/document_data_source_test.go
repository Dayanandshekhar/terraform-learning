// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ssm_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSSMDocumentDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_ssm_document.test"
	resourceName := "aws_ssm_document.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDocumentDataSourceConfig_basic(rName, "JSON"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "document_format", resourceName, "document_format"),
					resource.TestCheckResourceAttr(dataSourceName, "document_version", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "document_type", resourceName, "document_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "content", resourceName, "content"),
				),
			},
			{
				Config: testAccDocumentDataSourceConfig_basic(rName, "YAML"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttr(dataSourceName, "document_format", "YAML"),
					resource.TestCheckResourceAttr(dataSourceName, "document_version", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "document_type", resourceName, "document_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "content"),
				),
			},
		},
	})
}

func TestAccSSMDocumentDataSource_managed(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_ssm_document.test"
	nameWithAWSPrefix := "AWS-StartEC2Instance"
	nameWithAWSSSOPrefix := "AWSSSO-CreateSSOUser"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDocumentDataSourceConfig_managed(nameWithAWSPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", nameWithAWSPrefix),
					acctest.CheckResourceAttrRegionalARNNoAccount(dataSourceName, "arn", "ssm", fmt.Sprintf("document/%s", nameWithAWSPrefix)),
				),
			},
			{
				Config: testAccDocumentDataSourceConfig_managed(nameWithAWSSSOPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", nameWithAWSSSOPrefix),
					acctest.CheckResourceAttrRegionalARNNoAccount(dataSourceName, "arn", "ssm", fmt.Sprintf("document/%s", nameWithAWSSSOPrefix)),
				),
			},
		},
	})
}

func testAccDocumentDataSourceConfig_basic(rName, documentFormat string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "test" {
  name          = %[1]q
  document_type = "Command"

  content = <<DOC
{
  "schemaVersion": "1.2",
  "description": "Check ip configuration of a Linux instance.",
  "parameters": {},
  "runtimeConfig": {
    "aws:runShellScript": {
      "properties": [
        {
          "id": "0.aws:runShellScript",
          "runCommand": [
            "ifconfig"
          ]
        }
      ]
    }
  }
}
DOC
}

data "aws_ssm_document" "test" {
  name            = aws_ssm_document.test.name
  document_format = %[2]q
}
`, rName, documentFormat)
}

func testAccDocumentDataSourceConfig_managed(name string) string {
	return fmt.Sprintf(`
data "aws_ssm_document" "test" {
  name = %[1]q
}
`, name)
}
