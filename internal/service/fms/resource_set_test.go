// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fms_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/fms"
	"github.com/aws/aws-sdk-go-v2/service/fms/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tffms "github.com/hashicorp/terraform-provider-aws/internal/service/fms"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccFMSResourceSet_basic(t *testing.T) {
	ctx := acctest.Context(t)

	var resourceset fms.GetResourceSetOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_fms_resource_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.FMSEndpointID)
			// testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.FMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSetConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceSetExists(ctx, resourceName, &resourceset),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_set_status"),
				),
			},
		},
	})
}

func TestAccFMSResourceSet_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var resourceset fms.GetResourceSetOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_fms_resource_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.FMSEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.FMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSetConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceSetExists(ctx, resourceName, &resourceset),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tffms.ResourceSet, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccFMSResourceSet_tags(t *testing.T) {
	ctx := acctest.Context(t)

	var resourceset fms.GetResourceSetOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_fms_resource_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.FMSEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.FMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSetConfig_tags1(rName, acctest.CtKey1, acctest.CtValue1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceSetExists(ctx, resourceName, &resourceset),
					resource.TestCheckResourceAttr(resourceName, "resource_set.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct1),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceSetConfig_tags2(rName, acctest.CtKey1, acctest.CtValue1Updated, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceSetExists(ctx, resourceName, &resourceset),
					resource.TestCheckResourceAttr(resourceName, "resource_set.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct2),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1Updated),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
			{
				Config: testAccResourceSetConfig_tags1(rName, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceSetExists(ctx, resourceName, &resourceset),
					resource.TestCheckResourceAttr(resourceName, "resource_set.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct1),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
		},
	})
}

func testAccCheckResourceSetDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).FMSClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_fms_resource_set" {
				continue
			}

			_, err := tffms.FindResourceSetByID(ctx, conn, rs.Primary.ID)
			if errs.IsA[*types.ResourceNotFoundException](err) {
				return nil
			}
			if err != nil {
				return err
			}

			return create.Error(names.BCMDataExports, create.ErrActionCheckingDestroyed, tffms.ResNameResourceSet, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckResourceSetExists(ctx context.Context, name string, resourceset *fms.GetResourceSetOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.FMS, create.ErrActionCheckingExistence, tffms.ResNameResourceSet, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.FMS, create.ErrActionCheckingExistence, tffms.ResNameResourceSet, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).FMSClient(ctx)
		resp, err := tffms.FindResourceSetByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return create.Error(names.FMS, create.ErrActionCheckingExistence, tffms.ResNameResourceSet, rs.Primary.ID, err)
		}

		*resourceset = *resp

		return nil
	}
}

func testAccPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).FMSClient(ctx)

	input := &fms.ListResourceSetsInput{}
	_, err := conn.ListResourceSets(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccResourceSetConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_fms_resource_set" "test" {
  resource_set {
    name = %[1]q
    resource_type_list = ["testing"]
    resource_set_status = "ACTIVE"
  }
}
`, rName)
}

func testAccResourceSetConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_fms_resource_set" "test" {
  resource_set {
    name = %[1]q
    resource_type_list = ["testing"]
    resource_set_status = "ACTIVE"
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccResourceSetConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_fms_resource_set" "test" {
  resource_set {
    name = %[1]q
    resource_type_list = ["testing"]
    resource_set_status = "ACTIVE"
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
