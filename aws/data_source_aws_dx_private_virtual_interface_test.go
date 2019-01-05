package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceAwsDxPrivateVirtualInterface_basic(t *testing.T) {
	key := "DX_CONNECTION_ID"
	connectionID := os.Getenv(key)
	if connectionID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}
	vifName := fmt.Sprintf("terraform-testacc-dxvif-%s", acctest.RandString(5))
	bgpAsn := randIntRange(64512, 65534)
	vlan := randIntRange(2049, 4094)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDxPrivateVirtualInterfaceConfigBasic(connectionID, vifName, bgpAsn, vlan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAwsDxPrivateVirtualInterfaceExists("data.aws_dx_private_virtual_interface.foo"),
					resource.TestCheckResourceAttr("data.aws_dx_private_virtual_interface.foo", "name", vifName),
					resource.TestCheckResourceAttr("data.aws_dx_private_virtual_interface.foo", "tags.%", "0"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAwsDxPrivateVirtualInterfaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		return nil
	}
}

func testAccDxPrivateVirtualInterfaceConfigBasic(cid, n string, bgpAsn, vlan int) string {
	return fmt.Sprintf(`
resource "aws_vpn_gateway" "foo" {
  tags = {
    Name = "%s"
  }
}

resource "aws_dx_private_virtual_interface" "foo" {
  connection_id    = "%s"

  vpn_gateway_id = "${aws_vpn_gateway.foo.id}"
  name           = "%s"
  vlan           = %d
  address_family = "ipv4"
  bgp_asn        = %d
}

data "aws_dx_private_virtual_interface" "foo" {
	virtual_interface_id    = "${aws_dx_private_virtual_interface.foo.id}"
}
`, n, cid, n, vlan, bgpAsn)
}
