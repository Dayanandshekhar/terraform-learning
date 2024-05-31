// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transfer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transfer"
	awstypes "github.com/aws/aws-sdk-go-v2/service/transfer/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tftransfer "github.com/hashicorp/terraform-provider-aws/internal/service/transfer"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccSSHKey_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var conf awstypes.SshPublicKey
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resourceName := "aws_transfer_ssh_key.test"

	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.TransferServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSSHKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSSHKeyConfig_basic(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSHKeyExists(ctx, resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "server_id", "aws_transfer_server.test", names.AttrID),
					resource.TestCheckResourceAttrPair(resourceName, names.AttrUserName, "aws_transfer_user.test", names.AttrUserName),
					resource.TestCheckResourceAttr(resourceName, "body", publicKey),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSSHKeyExists(ctx context.Context, n string, res *awstypes.SshPublicKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Transfer SSH Public Key ID not set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).TransferClient(ctx)
		serverID, userName, sshKeyID, err := tftransfer.DecodeSSHKeyID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing Transfer SSH Public Key ID: %s", err)
		}

		describe, err := conn.DescribeUser(ctx, &transfer.DescribeUserInput{
			ServerId: aws.String(serverID),
			UserName: aws.String(userName),
		})

		if err != nil {
			return err
		}

		for _, sshPublicKey := range describe.User.SshPublicKeys {
			if sshKeyID == *sshPublicKey.SshPublicKeyId {
				*res = sshPublicKey
				return nil
			}
		}

		return fmt.Errorf("Transfer SSH Public Key does not exist")
	}
}

func testAccCheckSSHKeyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).TransferClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_transfer_ssh_key" {
				continue
			}
			serverID, userName, sshKeyID, err := tftransfer.DecodeSSHKeyID(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("error parsing Transfer SSH Public Key ID: %w", err)
			}

			describe, err := conn.DescribeUser(ctx, &transfer.DescribeUserInput{
				UserName: aws.String(userName),
				ServerId: aws.String(serverID),
			})

			if errs.IsA[*awstypes.ResourceNotFoundException](err) {
				continue
			}

			if err != nil {
				return err
			}

			for _, sshPublicKey := range describe.User.SshPublicKeys {
				if sshKeyID == *sshPublicKey.SshPublicKeyId {
					return fmt.Errorf("Transfer SSH Public Key still exists")
				}
			}
		}

		return nil
	}
}

func testAccSSHKeyConfig_basic(rName, publicKey string) string {
	return fmt.Sprintf(`
resource "aws_transfer_server" "test" {
  identity_provider_type = "SERVICE_MANAGED"
}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "transfer.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test" {
  name = %[1]q
  role = aws_iam_role.test.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowFullAccesstoS3",
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_transfer_user" "test" {
  server_id = aws_transfer_server.test.id
  user_name = "tftestuser"
  role      = aws_iam_role.test.arn
}

resource "aws_transfer_ssh_key" "test" {
  server_id = aws_transfer_server.test.id
  user_name = aws_transfer_user.test.user_name
  body      = "%[2]s"
}
`, rName, publicKey)
}
