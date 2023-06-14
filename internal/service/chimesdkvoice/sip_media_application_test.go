package chimesdkvoice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/chime"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfchimesdkvoice "github.com/hashicorp/terraform-provider-aws/internal/service/chimesdkvoice"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/chimesdkvoice"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccChimeSdkVoiceSipMediaApplication_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var chimeSipMediaApplication *chimesdkvoice.SipMediaApplication

	chimeSipMediaApplicationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_chimesdkvoice_sip_media_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckRegion(t, endpoints.UsEast1RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, chimesdkvoice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSipMediaApplicationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSipMediaApplicationConfig_basic(chimeSipMediaApplicationName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSipMediaApplicationExists(ctx, resourceName, chimeSipMediaApplication),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "aws_region", endpoints.UsEast1RegionID),
					resource.TestCheckResourceAttr(resourceName, "name", chimeSipMediaApplicationName),
					resource.TestCheckResourceAttrSet(resourceName, "endpoints.0.lambda_arn"),
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

func TestAccChimeSdkVoiceSipMediaApplication_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var chimeSipMediaApplication *chimesdkvoice.SipMediaApplication

	chimeSipMediaApplicationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_chimesdkvoice_sip_media_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckRegion(t, endpoints.UsEast1RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, chime.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSipMediaApplicationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSipMediaApplicationConfig_basic(chimeSipMediaApplicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSipMediaApplicationExists(ctx, resourceName, chimeSipMediaApplication),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfchimesdkvoice.ResourceSipMediaApplication(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccChimeSdkVoiceSipMediaApplication_update(t *testing.T) {
	ctx := acctest.Context(t)
	var chimeSipMediaApplication *chimesdkvoice.SipMediaApplication

	chimeSipMediaApplicationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	chimeSipMediaApplicationNameUpdated := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resourceName := "aws_chimesdkvoice_sip_media_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckRegion(t, endpoints.UsEast1RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, chimesdkvoice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSipMediaApplicationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSipMediaApplicationConfig_basic(chimeSipMediaApplicationName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSipMediaApplicationExists(ctx, resourceName, chimeSipMediaApplication),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "aws_region", endpoints.UsEast1RegionID),
					resource.TestCheckResourceAttr(resourceName, "name", chimeSipMediaApplicationName),
					resource.TestCheckResourceAttrSet(resourceName, "endpoints.0.lambda_arn"),
				),
			},
			{
				Config: testAccSipMediaApplicationConfig_basic(chimeSipMediaApplicationNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSipMediaApplicationExists(ctx, resourceName, chimeSipMediaApplication),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "aws_region", endpoints.UsEast1RegionID),
					resource.TestCheckResourceAttr(resourceName, "name", chimeSipMediaApplicationNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "endpoints.0.lambda_arn"),
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

func TestAccChimeSdkVoiceSipMediaApplication_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var sipMediaApplication *chimesdkvoice.SipMediaApplication

	sipMediaApplicationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_chimesdkvoice_sip_media_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckRegion(t, endpoints.UsEast1RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, chimesdkvoice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSipMediaApplicationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSipMediaApplicationConfig_tags1(sipMediaApplicationName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSipMediaApplicationExists(ctx, resourceName, sipMediaApplication),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "aws_region", endpoints.UsEast1RegionID),
					resource.TestCheckResourceAttr(resourceName, "name", sipMediaApplicationName),
					resource.TestCheckResourceAttrSet(resourceName, "endpoints.0.lambda_arn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSipMediaApplicationConfig_tags1(sipMediaApplicationName, "key1", "value1updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSipMediaApplicationExists(ctx, resourceName, sipMediaApplication),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "aws_region", endpoints.UsEast1RegionID),
					resource.TestCheckResourceAttr(resourceName, "name", sipMediaApplicationName),
					resource.TestCheckResourceAttrSet(resourceName, "endpoints.0.lambda_arn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
				),
			},
		},
	})
}

func testAccCheckSipMediaApplicationExists(ctx context.Context, name string, vc *chimesdkvoice.SipMediaApplication) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Chime voice connector ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ChimeSDKVoiceConn(ctx)
		input := &chimesdkvoice.GetSipMediaApplicationInput{
			SipMediaApplicationId: aws.String(rs.Primary.ID),
		}
		resp, err := conn.GetSipMediaApplicationWithContext(ctx, input)
		if err != nil {
			return err
		}

		vc = resp.SipMediaApplication

		return nil
	}
}

func testAccCheckSipMediaApplicationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_chimesdkvoice_chime_sip_media_application" {
				continue
			}
			conn := acctest.Provider.Meta().(*conns.AWSClient).ChimeSDKVoiceConn(ctx)
			input := &chimesdkvoice.GetSipMediaApplicationInput{
				SipMediaApplicationId: aws.String(rs.Primary.ID),
			}
			resp, err := conn.GetSipMediaApplicationWithContext(ctx, input)
			if err == nil {
				if resp.SipMediaApplication != nil && aws.StringValue(resp.SipMediaApplication.Name) != "" {
					return fmt.Errorf("error ChimeSdkVoice Sip Media Application still exists")
				}
			}
			return nil
		}
		return nil
	}
}

func testAccSipMediaApplicationConfigBase(rName string) string {
	return fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "test" {
  filename         = "test-fixtures/lambdatest.zip"
  source_code_hash = filebase64sha256("test-fixtures/lambdatest.zip")
  function_name    = %[1]q
  role             = aws_iam_role.test.arn
  runtime          = "nodejs16.x"
  handler          = "index.handler"
}
`, rName)
}

func testAccSipMediaApplicationConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccSipMediaApplicationConfigBase(rName),
		fmt.Sprintf(`
resource "aws_chimesdkvoice_sip_media_application" "test" {
  name       = %[1]q
  aws_region = data.aws_region.current.name
  endpoints {
    lambda_arn = aws_lambda_function.test.arn
  }
}
`, rName))
}

func testAccSipMediaApplicationConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(
		testAccSipMediaApplicationConfigBase(rName),
		fmt.Sprintf(`
resource "aws_chimesdkvoice_sip_media_application" "test" {
  name       = %[1]q
  aws_region = data.aws_region.current.name
  endpoints {
    lambda_arn = aws_lambda_function.test.arn
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccSipMediaApplicationConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(
		testAccSipMediaApplicationConfigBase(rName),
		fmt.Sprintf(`
resource "aws_chimesdkvoice_sip_media_application" "test" {
  name       = %[1]q
  aws_region = data.aws_region.current.name
  endpoints {
    lambda_arn = aws_lambda_function.test.arn
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}
