package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAwsDxTransitVirtualInterface_basic(t *testing.T) {
	key := "DX_CONNECTION_ID"
	connectionId := os.Getenv(key)
	if connectionId == "" {
		t.Skipf("Environment variable %s is not set", key)
	}

	resourceName := "aws_dx_transit_virtual_interface.test"
	rName := fmt.Sprintf("tf-testacc-transit-vif-%s", acctest.RandString(9))
	amzAsn := randIntRange(64512, 65534)
	bgpAsn := randIntRange(64512, 65534)
	vlan := randIntRange(2049, 4094)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsDxTransitVirtualInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDxTransitVirtualInterfaceConfig_basic(connectionId, rName, amzAsn, bgpAsn, vlan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsDxTransitVirtualInterfaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "1500"),
					resource.TestCheckResourceAttr(resourceName, "jumbo_frame_capable", "true"),
				),
			},
			{
				Config: testAccDxTransitVirtualInterfaceConfig_updated(connectionId, rName, amzAsn, bgpAsn, vlan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsDxTransitVirtualInterfaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "9001"),
				),
			},
			// Test import.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAwsDxTransitVirtualInterfaceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).dxconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_dx_transit_virtual_interface" {
			continue
		}

		input := &directconnect.DescribeVirtualInterfacesInput{
			VirtualInterfaceId: aws.String(rs.Primary.ID),
		}

		resp, err := conn.DescribeVirtualInterfaces(input)
		if err != nil {
			return err
		}
		for _, v := range resp.VirtualInterfaces {
			if *v.VirtualInterfaceId == rs.Primary.ID && !(*v.VirtualInterfaceState == directconnect.VirtualInterfaceStateDeleted) {
				return fmt.Errorf("[DESTROY ERROR] Dx Transit VIF (%s) not deleted", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckAwsDxTransitVirtualInterfaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		return nil
	}
}

func testAccDxTransitVirtualInterfaceConfig_basic(cid, rName string, amzAsn, bgpAsn, vlan int) string {
	return fmt.Sprintf(`
resource "aws_dx_gateway" "test" {
  name            = %[2]q
  amazon_side_asn = %[3]d
}

resource "aws_dx_transit_virtual_interface" "test" {
  connection_id    = %[1]q

  dx_gateway_id  = "${aws_vpn_gateway.test.id}"
  name           = %[2]q
  vlan           = %[5]d
  address_family = "ipv4"
  bgp_asn        = %[4]d
}
`, cid, rName, amzAsn, bgpAsn, vlan)
}

func testAccDxTransitVirtualInterfaceConfig_updated(cid, rName string, amzAsn, bgpAsn, vlan int) string {
	return fmt.Sprintf(`
resource "aws_dx_gateway" "test" {
  name            = %[2]q
  amazon_side_asn = %[3]d
}

resource "aws_dx_transit_virtual_interface" "test" {
  connection_id    = %[1]q

  dx_gateway_id  = "${aws_vpn_gateway.test.id}"
  name           = %[2]q
  vlan           = %[5]d
  address_family = "ipv4"
  bgp_asn        = %[4]d
  mtu            = 9001

  tags = {
    Environment = "test"
  }
}
`, cid, rName, amzAsn, bgpAsn, vlan)
}
