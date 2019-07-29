package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costandusagereportservice"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAwsCurReportDefinition_basic(t *testing.T) {
	oldvar := os.Getenv("AWS_DEFAULT_REGION")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	defer os.Setenv("AWS_DEFAULT_REGION", oldvar)

	resourceName := "aws_cur_report_definition.test"

	redshiftQuicksightReportName := acctest.RandomWithPrefix("tf_acc_test")
	redshiftQuicksightBucketName := fmt.Sprintf("tf-test-bucket-%d-%s", acctest.RandInt(), redshiftQuicksightReportName)
	athenaReportName := acctest.RandomWithPrefix("tf_acc_test")
	athenaBucketName := fmt.Sprintf("tf-test-bucket-%d-%s", acctest.RandInt(), athenaReportName)
	bucketRegion := "us-east-1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSCur(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsCurReportDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsCurReportDefinitionConfig_basic_redshift_quicksight(redshiftQuicksightReportName, redshiftQuicksightBucketName, bucketRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsCurReportDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "report_name", redshiftQuicksightReportName),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "format", "textORcsv"),
					resource.TestCheckResourceAttr(resourceName, "compression", "GZIP"),
					resource.TestCheckResourceAttr(resourceName, "additional_schema_elements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket", redshiftQuicksightBucketName),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "s3_region", bucketRegion),
					resource.TestCheckResourceAttr(resourceName, "additional_artifacts.#", "2"),
				),
			},
			{
				Config: testAccAwsCurReportDefinitionConfig_basic_athena(athenaReportName, athenaBucketName, bucketRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsCurReportDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "report_name", athenaReportName),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "format", "Parquet"),
					resource.TestCheckResourceAttr(resourceName, "compression", "Parquet"),
					resource.TestCheckResourceAttr(resourceName, "additional_schema_elements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket", athenaBucketName),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "s3_region", bucketRegion),
					resource.TestCheckResourceAttr(resourceName, "additional_artifacts.#", "1"),
				),
			},
		},
	})
}

func testAccCheckAwsCurReportDefinitionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).costandusagereportconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cur_report_definition" {
			continue
		}
		reportName := rs.Primary.ID
		matchingReportDefinition, err := describeCurReportDefinition(conn, reportName)
		if err != nil {
			return err
		}
		if matchingReportDefinition != nil {
			return fmt.Errorf("Report Definition still exists: %q", rs.Primary.ID)
		}
	}
	return nil

}

func testAccCheckAwsCurReportDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).costandusagereportconn

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}
		reportName := rs.Primary.ID
		matchingReportDefinition, err := describeCurReportDefinition(conn, reportName)
		if err != nil {
			return err
		}
		if matchingReportDefinition == nil {
			return fmt.Errorf("Report Definition does not exist: %q", rs.Primary.ID)
		}
		return nil
	}
}

func testAccPreCheckAWSCur(t *testing.T) {
	conn := testAccProvider.Meta().(*AWSClient).costandusagereportconn

	input := &costandusagereportservice.DescribeReportDefinitionsInput{
		MaxResults: aws.Int64(5),
	}

	_, err := conn.DescribeReportDefinitions(input)

	if testAccPreCheckSkipError(err) || isAWSErr(err, "AccessDeniedException", "linked account is not allowed to modify report preference") {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

// note: cur report definitions are currently only supported in us-east-1
func testAccAwsCurReportDefinitionConfig_basic(reportName string, bucketName string, bucketRegion string) string {
	return fmt.Sprintf(`
provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "test" {
  bucket        = "%[2]s"
  acl           = "private"
  force_destroy = true
  region        = "%[3]s"
}

resource "aws_s3_bucket_policy" "test" {
  bucket = "${aws_s3_bucket.test.id}"

  policy = <<POLICY
{
    "Version": "2008-10-17",
    "Id": "s3policy",
    "Statement": [
        {
            "Sid": "AllowCURBillingACLPolicy",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::386209384616:root"
            },
            "Action": [
                "s3:GetBucketAcl",
                "s3:GetBucketPolicy"
            ],
            "Resource": "arn:aws:s3:::${aws_s3_bucket.test.id}"
        },
        {
            "Sid": "AllowCURPutObject",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::386209384616:root"
            },
            "Action": "s3:PutObject",
            "Resource": "arn:aws:s3:::${aws_s3_bucket.test.id}/*"
        }
    ]
}
POLICY
}
`, reportName, bucketName, bucketRegion)
}

func testAccAwsCurReportDefinitionConfig_basic_redshift_quicksight(reportName string, bucketName string, bucketRegion string) string {
	reportDefinition := fmt.Sprintf(`
resource "aws_cur_report_definition" "redshift_quicksight" {
  report_name                = "%[1]s"
  time_unit                  = "DAILY"
  format                     = "textORcsv"
  compression                = "GZIP"
  additional_schema_elements = ["RESOURCES"]
  s3_bucket                  = "${aws_s3_bucket.test.id}"
  s3_prefix                  = ""
  s3_region                  = "${aws_s3_bucket.test.region}"
  additional_artifacts       = ["ATHENA"]
}
`, reportName)

	return testAccAwsCurReportDefinitionConfig_basic(reportName, bucketName, bucketRegion) + reportDefinition
}

//``
func testAccAwsCurReportDefinitionConfig_basic_athena(reportName string, bucketName string, bucketRegion string) string {
	reportDefinition := fmt.Sprintf(`
resource "aws_cur_report_definition" "athena" {
  report_name                = "%[1]s"
  time_unit                  = "DAILY"
  format                     = "Parquet"
  compression                = "Parquet"
  additional_schema_elements = ["RESOURCES"]
  s3_bucket                  = "${aws_s3_bucket.test.id}"
  s3_prefix                  = ""
  s3_region                  = "${aws_s3_bucket.test.region}"
  additional_artifacts       = ["REDSHIFT", "QUICKSIGHT"]
}
`, reportName)

	return testAccAwsCurReportDefinitionConfig_basic(reportName, bucketName, bucketRegion) + reportDefinition
}
