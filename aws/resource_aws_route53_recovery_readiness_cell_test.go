package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53recoveryreadiness"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAwsRoute53RecoveryReadinessCell_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53_recovery_readiness_cell.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAwsRoute53RecoveryReadiness(t) },
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessCellDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessCellConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessCellExists(resourceName),
					testAccMatchResourceAttrGlobalARN(resourceName, "cell_arn", "route53-recovery-readiness", regexp.MustCompile(`cell/.+`)),
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

func TestAccAwsRoute53RecoveryReadinessCell_nestedCell(t *testing.T) {
	rNameParent := acctest.RandomWithPrefix("tf-acc-test-parent")
	rNameChild := acctest.RandomWithPrefix("tf-acc-test-child")
	resourceNameParent := "aws_route53_recovery_readiness_cell.test_parent"
	resourceNameChild := "aws_route53_recovery_readiness_cell.test_child"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAwsRoute53RecoveryReadiness(t) },
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessCellDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessCellChildConfig(rNameChild),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessCellExists(resourceNameChild),
					testAccMatchResourceAttrGlobalARN(resourceNameChild, "cell_arn", "route53-recovery-readiness", regexp.MustCompile(`cell/.+`)),
				),
			},
			{
				Config: testAccAwsRoute53RecoveryReadinessCellParentConfig(rNameChild, rNameParent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessCellExists(resourceNameParent),
					testAccMatchResourceAttrGlobalARN(resourceNameParent, "cell_arn", "route53-recovery-readiness", regexp.MustCompile(`cell/.+`)),
					resource.TestCheckResourceAttr(resourceNameParent, "cells.#", "1"),
					resource.TestCheckResourceAttr(resourceNameParent, "parent_readiness_scopes.#", "0"),
					testAccCheckAwsRoute53RecoveryReadinessCellExists(resourceNameChild),
					testAccMatchResourceAttrGlobalARN(resourceNameChild, "cell_arn", "route53-recovery-readiness", regexp.MustCompile(`cell/.+`)),
					resource.TestCheckResourceAttr(resourceNameChild, "cells.#", "0"),
				),
			},
			{
				Config: testAccAwsRoute53RecoveryReadinessCellParentConfig(rNameChild, rNameParent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameChild, "parent_readiness_scopes.#", "1"),
				),
			},
			{
				ResourceName:      resourceNameParent,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceNameChild,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAwsRoute53RecoveryReadinessCellDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).route53recoveryreadinessconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_route53_recovery_readiness_cell" {
			continue
		}

		input := &route53recoveryreadiness.GetCellInput{
			CellName: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetCell(input)
		if err == nil {
			return fmt.Errorf("Route53RecoveryReadiness Channel (%s) not deleted", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAwsRoute53RecoveryReadinessCellExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := testAccProvider.Meta().(*AWSClient).route53recoveryreadinessconn

		input := &route53recoveryreadiness.GetCellInput{
			CellName: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetCell(input)

		return err
	}
}

func testAccPreCheckAwsRoute53RecoveryReadiness(t *testing.T) {
	conn := testAccProvider.Meta().(*AWSClient).route53recoveryreadinessconn

	input := &route53recoveryreadiness.ListCellsInput{}

	_, err := conn.ListCells(input)

	if testAccPreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccAwsRoute53RecoveryReadinessCellConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_route53_recovery_readiness_cell" "test" {
  cell_name = %q
}
`, rName)
}

func testAccAwsRoute53RecoveryReadinessCellChildConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_route53_recovery_readiness_cell" "test_child" {
  cell_name = %q
}
`, rName)
}

func testAccAwsRoute53RecoveryReadinessCellParentConfig(rName, rName2 string) string {
	return composeConfig(testAccAwsRoute53RecoveryReadinessCellChildConfig(rName), fmt.Sprintf(`
resource "aws_route53_recovery_readiness_cell" "test_parent" {
  cell_name = %q
  cells     = [aws_route53_recovery_readiness_cell.test_child.cell_arn]
}
`, rName2))
}
