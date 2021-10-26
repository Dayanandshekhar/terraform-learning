package fsx_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/fsx"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tffsx "github.com/hashicorp/terraform-provider-aws/internal/service/fsx"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccAWSFsxStorageVirtualMachine_basic(t *testing.T) {
	var svm fsx.StorageVirtualMachine
	resourceName := "aws_fsx_storage_virtual_machine.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxStorageVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsFsxStorageVirtualMachineConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxStorageVirtualMachineExists(resourceName, &svm),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "fsx", regexp.MustCompile(`storage-virtual-machine/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "file_system_id", "aws_fsx_ontap_file_system.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccAWSFsxStorageVirtualMachine_rootVolumeSecurityStyle(t *testing.T) {
	var svm fsx.StorageVirtualMachine
	resourceName := "aws_fsx_storage_virtual_machine.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxStorageVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsFsxStorageVirtualMachineConfigRootVolumeSecurityStyle(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxStorageVirtualMachineExists(resourceName, &svm),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "root_volume_security_style", "UNIX"),
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

func TestAccAWSFsxStorageVirtualMachine_svmAdminPassword(t *testing.T) {
	var svm fsx.StorageVirtualMachine
	resourceName := "aws_fsx_storage_virtual_machine.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	passUpdated := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxStorageVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsFsxStorageVirtualMachineConfigSvmAdminPassword(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxStorageVirtualMachineExists(resourceName, &svm),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "svm_admin_password", rName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_group_ids"},
			},
			{
				Config: testAccAwsFsxStorageVirtualMachineConfigSvmAdminPassword(rName, passUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxStorageVirtualMachineExists(resourceName, &svm),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "svm_admin_password", passUpdated),
				),
			},
		},
	})
}

func TestAccAWSFsxStorageVirtualMachine_disappears(t *testing.T) {
	var svm fsx.StorageVirtualMachine
	resourceName := "aws_fsx_storage_virtual_machine.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxStorageVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsFsxStorageVirtualMachineConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxStorageVirtualMachineExists(resourceName, &svm),
					acctest.CheckResourceDisappears(acctest.Provider, tffsx.ResourceStorageVirtualMachine(), resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tffsx.ResourceStorageVirtualMachine(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSFsxStorageVirtualMachine_tags(t *testing.T) {
	var svm1, svm2, svm3 fsx.StorageVirtualMachine
	resourceName := "aws_fsx_storage_virtual_machine.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxStorageVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsFsxStorageVirtualMachineConfigTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxStorageVirtualMachineExists(resourceName, &svm1),
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
				Config: testAccAwsFsxStorageVirtualMachineConfigTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxStorageVirtualMachineExists(resourceName, &svm2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAwsFsxStorageVirtualMachineConfigTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxStorageVirtualMachineExists(resourceName, &svm3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckFsxStorageVirtualMachineExists(resourceName string, svm *fsx.StorageVirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).FSxConn

		resp, err := tffsx.FindStorageVirtualMachineByID(conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp == nil {
			return fmt.Errorf("FSx Storage Virtual Machine (%s) not found", rs.Primary.ID)
		}

		*svm = *resp

		return nil
	}
}

func testAccCheckFsxStorageVirtualMachineDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).FSxConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_fsx_storage_virtual_machine" {
			continue
		}

		svm, err := tffsx.FindStorageVirtualMachineByID(conn, rs.Primary.ID)
		if tfresource.NotFound(err) {
			continue
		}

		if svm != nil {
			return fmt.Errorf("FSx Storage Virtual Machine (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAwsFsxStorageVirtualMachineConfigBase() string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), `
data "aws_partition" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "test1" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]
}

resource "aws_subnet" "test2" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]
}

resource "aws_fsx_ontap_file_system" "test" {
  storage_capacity    = 1024
  subnet_ids          = [aws_subnet.test1.id, aws_subnet.test2.id]
  deployment_type     = "MULTI_AZ_1"
  throughput_capacity = 512
  preferred_subnet_id = aws_subnet.test1.id
}
`)
}

func testAccAwsFsxStorageVirtualMachineConfigBasic(rName string) string {
	return acctest.ConfigCompose(testAccAwsFsxStorageVirtualMachineConfigBase(), fmt.Sprintf(`
resource "aws_fsx_storage_virtual_machine" "test" {
  name           = %[1]q
  file_system_id = aws_fsx_ontap_file_system.test.id
}
`, rName))
}

func testAccAwsFsxStorageVirtualMachineConfigRootVolumeSecurityStyle(rName string) string {
	return acctest.ConfigCompose(testAccAwsFsxStorageVirtualMachineConfigBase(), fmt.Sprintf(`
resource "aws_fsx_storage_virtual_machine" "test" {
  name                       = %[1]q
  file_system_id             = aws_fsx_ontap_file_system.test.id
  root_volume_security_style = "UNIX"
}
`, rName))
}

func testAccAwsFsxStorageVirtualMachineConfigSvmAdminPassword(rName, pass string) string {
	return acctest.ConfigCompose(testAccAwsFsxStorageVirtualMachineConfigBase(), fmt.Sprintf(`
resource "aws_fsx_storage_virtual_machine" "test" {
  name               = %[1]q
  file_system_id     = aws_fsx_ontap_file_system.test.id
  svm_admin_password = %[2]q
}
`, rName, pass))
}

func testAccAwsFsxStorageVirtualMachineConfigTags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccAwsFsxStorageVirtualMachineConfigBase(), fmt.Sprintf(`
resource "aws_fsx_storage_virtual_machine" "test" {
  name           = %[1]q
  file_system_id = aws_fsx_ontap_file_system.test.id
  tags = {
    %[1]q = %[2]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccAwsFsxStorageVirtualMachineConfigTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccAwsFsxStorageVirtualMachineConfigBase(), fmt.Sprintf(`
resource "aws_fsx_storage_virtual_machine" "test" {
  name           = %[1]q
  file_system_id = aws_fsx_ontap_file_system.test.id
  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}
