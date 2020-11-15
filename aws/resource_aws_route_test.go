package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAWSRoute_basic(t *testing.T) {
	var route ec2.Route
	var routeTable ec2.RouteTable
	resourceName := "aws_route.test"
	igwResourceName := "aws_internet_gateway.test"
	rtResourceName := "aws_route_table.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4InternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteGatewayRoute(igwResourceName, &route),
					testAccCheckRouteTableExists(rtResourceName, &routeTable),
					testAccCheckAWSRouteNumberOfRoutes(&routeTable, 2),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_disappears(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4InternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsRoute(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_EgressOnlyInternetGateway(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	eoigwResourceName := "aws_egress_only_internet_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6EgressOnlyInternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteEgressOnlyInternetGatewayRoute(eoigwResourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				// Verify that expanded form of the destination CIDR causes no diff.
				Config:   testAccAWSRouteConfigIpv6EgressOnlyInternetGateway(rName, "::0/0"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_InternetGateway(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	igwResourceName := "aws_internet_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6InternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteGatewayRoute(igwResourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_Instance(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	instanceResourceName := "aws_instance.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6Instance(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteInstanceRoute(instanceResourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_NetworkInterface(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	eniResourceName := "aws_network_interface.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6NetworkInterface(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteNetworkInterfaceRoute(eniResourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_VpcPeeringConnection(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	pcxResourceName := "aws_vpc_peering_connection.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6VpcPeeringConnection(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteVpcPeeringConnectionRoute(pcxResourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_InternetGateway(t *testing.T) {
	var route ec2.Route
	var routeTable ec2.RouteTable
	resourceName := "aws_route.test"
	igwResourceName := "aws_internet_gateway.test"
	rtResourceName := "aws_route_table.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"
	changedDestinationCidr := "10.2.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4InternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteGatewayRoute(igwResourceName, &route),
					testAccCheckRouteTableExists(rtResourceName, &routeTable),
					testAccCheckAWSRouteNumberOfRoutes(&routeTable, 2),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
				),
			},
			{
				Config: testAccAWSRouteConfigIpv4InternetGateway(rName, changedDestinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteGatewayRoute(igwResourceName, &route),
					testAccCheckRouteTableExists(rtResourceName, &routeTable),
					testAccCheckAWSRouteNumberOfRoutes(&routeTable, 2),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", changedDestinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_Instance(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	instanceResourceName := "aws_instance.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4Instance(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteInstanceRoute(instanceResourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_DoesNotCrashWithVpcEndpoint(t *testing.T) {
	var route ec2.Route
	var routeTable ec2.RouteTable
	resourceName := "aws_route.test"
	rtResourceName := "aws_route_table.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigWithVpcEndpoint(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckRouteTableExists(rtResourceName, &routeTable),
					testAccCheckAWSRouteNumberOfRoutes(&routeTable, 3),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_TransitGateway(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4TransitGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					testAccCheckAWSRouteTransitGatewayRoute(transitGatewayResourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_id", transitGatewayResourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_LocalGatewayID(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	localGatewayDataSourceName := "data.aws_ec2_local_gateway.first"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSOutpostsOutposts(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteResourceConfigLocalGatewayID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					resource.TestCheckResourceAttrPair(resourceName, "local_gateway_id", localGatewayDataSourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_ConditionalCidrBlock(t *testing.T) {
	var route ec2.Route
	resourceName := "aws_route.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigConditionalIpv4Ipv6(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
				),
			},
			{
				Config: testAccAWSRouteConfigConditionalIpv4Ipv6(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", "::/0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_VpcEndpointId(t *testing.T) {
	var route ec2.Route
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route.test"
	vpcEndpointResourceName := "aws_vpc_endpoint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteResourceConfigVpcEndpointId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRouteExists(resourceName, &route),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_endpoint_id", vpcEndpointResourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAWSRouteExists(n string, res *ec2.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s\n", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn
		r, err := resourceAwsRouteFindRoute(
			conn,
			rs.Primary.Attributes["route_table_id"],
			rs.Primary.Attributes["destination_cidr_block"],
			rs.Primary.Attributes["destination_ipv6_cidr_block"],
		)

		if err != nil {
			return err
		}

		if r == nil {
			return fmt.Errorf("Route not found")
		}

		*res = *r

		return nil
	}
}

func testAccCheckAWSRouteDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_route" {
			continue
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn
		route, err := resourceAwsRouteFindRoute(
			conn,
			rs.Primary.Attributes["route_table_id"],
			rs.Primary.Attributes["destination_cidr_block"],
			rs.Primary.Attributes["destination_ipv6_cidr_block"],
		)

		if route == nil && err == nil {
			return nil
		}
	}

	return nil
}

func testAccAWSRouteImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		destination := rs.Primary.Attributes["destination_cidr_block"]
		if v, ok := rs.Primary.Attributes["destination_ipv6_cidr_block"]; ok && v != "" {
			destination = v
		}

		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["route_table_id"], destination), nil
	}
}

func testAccCheckAWSRouteEgressOnlyInternetGatewayRoute(n string, route *ec2.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s\n", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if gwId := aws.StringValue(route.EgressOnlyInternetGatewayId); gwId != rs.Primary.ID {
			return fmt.Errorf("Egress Only Internet Gateway ID (Expected=%s, Actual=%s)\n", rs.Primary.ID, gwId)
		}

		return nil
	}
}

func testAccCheckAWSRouteGatewayRoute(n string, route *ec2.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s\n", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if gwId := aws.StringValue(route.GatewayId); gwId != rs.Primary.ID {
			return fmt.Errorf("Internet Gateway ID (Expected=%s, Actual=%s)\n", rs.Primary.ID, gwId)
		}

		return nil
	}
}

func testAccCheckAWSRouteInstanceRoute(n string, route *ec2.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s\n", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if gwId := aws.StringValue(route.InstanceId); gwId != rs.Primary.ID {
			return fmt.Errorf("Instance ID (Expected=%s, Actual=%s)\n", rs.Primary.ID, gwId)
		}

		return nil
	}
}

func testAccCheckAWSRouteNetworkInterfaceRoute(n string, route *ec2.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s\n", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if gwId := aws.StringValue(route.NetworkInterfaceId); gwId != rs.Primary.ID {
			return fmt.Errorf("Network Interface ID (Expected=%s, Actual=%s)\n", rs.Primary.ID, gwId)
		}

		return nil
	}
}

func testAccCheckAWSRouteTransitGatewayRoute(n string, route *ec2.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s\n", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if gwId := aws.StringValue(route.TransitGatewayId); gwId != rs.Primary.ID {
			return fmt.Errorf("Transit Gateway ID (Expected=%s, Actual=%s)\n", rs.Primary.ID, gwId)
		}

		return nil
	}
}

func testAccCheckAWSRouteVpcPeeringConnectionRoute(n string, route *ec2.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s\n", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if gwId := aws.StringValue(route.VpcPeeringConnectionId); gwId != rs.Primary.ID {
			return fmt.Errorf("VPC Peering Connection ID (Expected=%s, Actual=%s)\n", rs.Primary.ID, gwId)
		}

		return nil
	}
}

func testAccCheckAWSRouteNumberOfRoutes(routeTable *ec2.RouteTable, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len := len(routeTable.Routes); len != n {
			return fmt.Errorf("Route Table has incorrect number of routes (Expected=%d, Actual=%d)\n", n, len)
		}

		return nil
	}
}

func testAccAWSRouteConfigIpv4InternetGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  gateway_id             = aws_internet_gateway.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv6InternetGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_egress_only_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  gateway_id                  = aws_internet_gateway.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv6NetworkInterface(rName, destinationCidr string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "current" {
  # Exclude usw2-az4 (us-west-2d) as it has limited instance types.
  blacklisted_zone_ids = ["usw2-az4"]
  state                = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.current.names[0]
  ipv6_cidr_block   = cidrsubnet(aws_vpc.test.ipv6_cidr_block, 8, 1)

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test" {
  subnet_id = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}


resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  network_interface_id        = aws_network_interface.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv6Instance(rName, destinationCidr string) string {
	return composeConfig(
		testAccLatestAmazonLinuxHvmEbsAmiConfig(),
		testAccAvailableEc2InstanceTypeForRegion("t3.micro", "t2.micro"),
		fmt.Sprintf(`
data "aws_availability_zones" "current" {
  # Exclude usw2-az4 (us-west-2d) as it has limited instance types.
  blacklisted_zone_ids = ["usw2-az4"]
  state                = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.current.names[0]
  ipv6_cidr_block   = cidrsubnet(aws_vpc.test.ipv6_cidr_block, 8, 1)

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type
  subnet_id     = aws_subnet.test.id

  ipv6_address_count = 1

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  instance_id                 = aws_instance.test.id
}
`, rName, destinationCidr))
}

func testAccAWSRouteConfigIpv6VpcPeeringConnection(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "source" {
  cidr_block                       = "10.0.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc" "target" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_peering_connection" "test" {
  vpc_id      = aws_vpc.source.id
  peer_vpc_id = aws_vpc.target.id
  auto_accept = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.source.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  vpc_peering_connection_id   = aws_vpc_peering_connection.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv6EgressOnlyInternetGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_egress_only_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  egress_only_gateway_id      = aws_egress_only_internet_gateway.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigWithVpcEndpoint(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = "10.3.0.0/16"
  gateway_id             = aws_internet_gateway.test.id

  # Forcing endpoint to create before route - without this the crash is a race.
  depends_on = [aws_vpc_endpoint.test]
}

data "aws_region" "current" {}

resource "aws_vpc_endpoint" "test" {
  vpc_id          = aws_vpc.test.id
  service_name    = "com.amazonaws.${data.aws_region.current.name}.s3"
  route_table_ids = [aws_route_table.test.id]
}
`, rName)
}

func testAccAWSRouteConfigIpv4TransitGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  # IncorrectState: Transit Gateway is not available in availability zone us-west-2d
  exclude_zone_ids = ["usw2-az4"]
  state            = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  destination_cidr_block = %[2]q
  route_table_id         = aws_route_table.test.id
  transit_gateway_id     = aws_ec2_transit_gateway_vpc_attachment.test.transit_gateway_id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigConditionalIpv4Ipv6(rName string, ipv6Route bool) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_egress_only_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

locals {
  ipv6             = %[2]t
  destination      = "0.0.0.0/0"
  destination_ipv6 = "::/0"
}

resource "aws_route" "test" {
  route_table_id = aws_route_table.test.id
  gateway_id     = aws_internet_gateway.test.id

  destination_cidr_block      = local.ipv6 ? "" : local.destination
  destination_ipv6_cidr_block = local.ipv6 ? local.destination_ipv6 : ""
}
`, rName, ipv6Route)
}

func testAccAWSRouteConfigIpv4Instance(rName, destinationCidr string) string {
	return composeConfig(
		testAccLatestAmazonLinuxHvmEbsAmiConfig(),
		testAccAvailableEc2InstanceTypeForRegion("t3.micro", "t2.micro"),
		fmt.Sprintf(`
data "aws_availability_zones" "current" {
  # Exclude usw2-az4 (us-west-2d) as it has limited instance types.
  blacklisted_zone_ids = ["usw2-az4"]
  state                = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.current.names[0]

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type
  subnet_id     = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_cidr_block      = %[2]q
  instance_id                 = aws_instance.test.id
}
`, rName, destinationCidr))
}

func testAccAWSRouteResourceConfigLocalGatewayID() string {
	return fmt.Sprintf(`
data "aws_ec2_local_gateways" "all" {}

data "aws_ec2_local_gateway" "first" {
  id = tolist(data.aws_ec2_local_gateways.all.ids)[0]
}

data "aws_ec2_local_gateway_route_tables" "all" {}

data "aws_ec2_local_gateway_route_table" "first" {
  local_gateway_route_table_id = tolist(data.aws_ec2_local_gateway_route_tables.all.ids)[0]
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_ec2_local_gateway_route_table_vpc_association" "example" {
  local_gateway_route_table_id = data.aws_ec2_local_gateway_route_table.first.id
  vpc_id                       = aws_vpc.test.id
}

resource "aws_route_table" "test" {
  vpc_id     = aws_vpc.test.id
  depends_on = [aws_ec2_local_gateway_route_table_vpc_association.example]
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = "172.16.1.0/24"
  local_gateway_id       = data.aws_ec2_local_gateway.first.id
}
`)
}

func testAccAWSRouteResourceConfigVpcEndpointId(rName string) string {
	return composeConfig(
		testAccAvailableAZsNoOptInConfig(),
		fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.10.10.0/25"

  tags = {
    Name = "tf-acc-test-load-balancer"
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 2, 0)
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = "tf-acc-test-load-balancer"
  }
}

resource "aws_lb" "test" {
  load_balancer_type = "gateway"
  name               = %[1]q

  subnet_mapping {
    subnet_id = aws_subnet.test.id
  }
}

resource "aws_vpc_endpoint_service" "test" {
  acceptance_required        = false
  allowed_principals         = [data.aws_caller_identity.current.arn]
  gateway_load_balancer_arns = [aws_lb.test.arn]
}

resource "aws_vpc_endpoint" "test" {
  service_name      = aws_vpc_endpoint_service.test.service_name
  subnet_ids        = [aws_subnet.test.id]
  vpc_endpoint_type = aws_vpc_endpoint_service.test.service_type
  vpc_id            = aws_vpc.test.id
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = "172.16.1.0/24"
  vpc_endpoint_id        = aws_vpc_endpoint.test.id
}
`, rName))
}
