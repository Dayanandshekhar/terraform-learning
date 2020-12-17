package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/lakeformation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAWSLakeFormationDataLakeSettingsDataSource_serial(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic": testAccAWSLakeFormationDataLakeSettingsDataSource_basic,
		// if more tests are added, they should be serial (data catalog is account-shared resource)
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccAWSLakeFormationDataLakeSettingsDataSource_basic(t *testing.T) {
	callerIdentityName := "data.aws_caller_identity.current"
	resourceName := "data.aws_lakeformation_data_lake_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPartitionHasServicePreCheck(lakeformation.EndpointsID, t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSLakeFormationDataLakeSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLakeFormationDataLakeSettingsDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "catalog_id", callerIdentityName, "account_id"),
					resource.TestCheckResourceAttr(resourceName, "data_lake_admins.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "data_lake_admins.0", callerIdentityName, "arn"),
				),
			},
		},
	})
}

const testAccAWSLakeFormationDataLakeSettingsDataSourceConfig_basic = `
data "aws_caller_identity" "current" {}

resource "aws_lakeformation_data_lake_settings" "test" {
  catalog_id       = data.aws_caller_identity.current.account_id
  data_lake_admins = [data.aws_caller_identity.current.arn]
}

data "aws_lakeformation_data_lake_settings" "test" {
  catalog_id = aws_lakeformation_data_lake_settings.test.catalog_id
}
`
