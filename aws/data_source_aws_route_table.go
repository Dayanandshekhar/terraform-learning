package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func dataSourceAwsRouteTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsRouteTableRead,

		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter": ec2CustomFiltersSchema(),
			"tags":   tagsSchemaComputed(),
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ipv6_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"egress_only_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"transit_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_peering_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_table_association_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"main": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getVpcFromSubnet(subnetId string, conn *ec2.EC2) (*string, error) {
	log.Printf("[DEBUG] Obtaning vpcId from Subnet: %s", subnetId)
	req := &ec2.DescribeSubnetsInput{}
	req.SubnetIds = []*string{aws.String(subnetId)}
	log.Printf("[DEBUG] Reading Subnet: %s", req)
	resp, err := conn.DescribeSubnets(req)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Subnets) == 0 {
		return nil, fmt.Errorf("no matching subnet found")
	}
	subnet := resp.Subnets[0]
	return subnet.VpcId, nil
}

func dataSourceAwsRouteTableReadFromVpc(subnetId string, conn *ec2.EC2) (*ec2.DescribeRouteTablesOutput, error) {
	log.Printf("[DEBUG] Obtaning default route of VPC from Subnet: %s", subnetId)
	vpcId, err := getVpcFromSubnet(subnetId, conn)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Obtained VpcID: %s", *vpcId)

	req := &ec2.DescribeRouteTablesInput{}
	req.Filters = buildEC2AttributeFilterList(
		map[string]string{
			"vpc-id":           *vpcId,
			"association.main": "true",
		},
	)
	log.Printf("[DEBUG] Reading Route Table: %s", req)
	resp, err := conn.DescribeRouteTables(req)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.RouteTables) == 0 {
		return nil, fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}
	return resp, nil
}

func dataSourceAwsRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn
	req := &ec2.DescribeRouteTablesInput{}
	vpcId, vpcIdOk := d.GetOk("vpc_id")
	subnetId, subnetIdOk := d.GetOk("subnet_id")
	gatewayId, gatewayIdOk := d.GetOk("gateway_id")
	rtbId, rtbOk := d.GetOk("route_table_id")
	tags, tagsOk := d.GetOk("tags")
	filter, filterOk := d.GetOk("filter")

	if !rtbOk && !vpcIdOk && !subnetIdOk && !gatewayIdOk && !filterOk && !tagsOk {
		return fmt.Errorf("One of route_table_id, vpc_id, subnet_id, gateway_id, filters, or tags must be assigned")
	}
	req.Filters = buildEC2AttributeFilterList(
		map[string]string{
			"route-table-id":         rtbId.(string),
			"vpc-id":                 vpcId.(string),
			"association.subnet-id":  subnetId.(string),
			"association.gateway-id": gatewayId.(string),
		},
	)
	req.Filters = append(req.Filters, buildEC2TagFilterList(
		keyvaluetags.New(tags.(map[string]interface{})).Ec2Tags(),
	)...)
	req.Filters = append(req.Filters, buildEC2CustomFilterList(
		filter.(*schema.Set),
	)...)

	log.Printf("[DEBUG] Reading Route Table: %s", req)
	resp, err := conn.DescribeRouteTables(req)
	if err != nil {
		return err
	}

	if resp == nil || len(resp.RouteTables) == 0 {
		if subnetIdOk && !rtbOk && !vpcIdOk && !gatewayIdOk && !filterOk && !tagsOk {
			// that means the user did send a subnet and nothing else

			// it also means that the AWS API returned an empty result
			// this edge case happens when a Subnet has no explicit route set
			// in this case, the default VPC route is used, but the route is considered implicit
			// and hence not returned. Somehow AWS thinks this is a good idea.
			//
			// source: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeRouteTables.html
			//
			// which means we have get the subnet's VPC and return its default route, for that's what the subnet uses.
			log.Printf("[DEBUG] Subnet %s has no explicit route. Returning the implicit route (VPC's default route).", subnetId)
			resp, err = dataSourceAwsRouteTableReadFromVpc(subnetId.(string), conn)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
		}
	}
	if len(resp.RouteTables) > 1 {
		return fmt.Errorf("Multiple Route Table matched; use additional constraints to reduce matches to a single Route Table")
	}

	rt := resp.RouteTables[0]

	d.SetId(aws.StringValue(rt.RouteTableId))
	d.Set("route_table_id", rt.RouteTableId)
	d.Set("vpc_id", rt.VpcId)

	if err := d.Set("tags", keyvaluetags.Ec2KeyValueTags(rt.Tags).IgnoreAws().Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	d.Set("owner_id", rt.OwnerId)
	if err := d.Set("routes", dataSourceRoutesRead(rt.Routes)); err != nil {
		return err
	}

	if err := d.Set("associations", dataSourceAssociationsRead(rt.Associations)); err != nil {
		return err
	}

	return nil
}

func dataSourceRoutesRead(ec2Routes []*ec2.Route) []map[string]interface{} {
	routes := make([]map[string]interface{}, 0, len(ec2Routes))
	// Loop through the routes and add them to the set
	for _, r := range ec2Routes {
		if r.GatewayId != nil && *r.GatewayId == "local" {
			continue
		}

		if r.Origin != nil && *r.Origin == "EnableVgwRoutePropagation" {
			continue
		}

		if r.DestinationPrefixListId != nil {
			// Skipping because VPC endpoint routes are handled separately
			// See aws_vpc_endpoint
			continue
		}

		m := make(map[string]interface{})

		if r.DestinationCidrBlock != nil {
			m["cidr_block"] = *r.DestinationCidrBlock
		}
		if r.DestinationIpv6CidrBlock != nil {
			m["ipv6_cidr_block"] = *r.DestinationIpv6CidrBlock
		}
		if r.EgressOnlyInternetGatewayId != nil {
			m["egress_only_gateway_id"] = *r.EgressOnlyInternetGatewayId
		}
		if r.GatewayId != nil {
			m["gateway_id"] = *r.GatewayId
		}
		if r.NatGatewayId != nil {
			m["nat_gateway_id"] = *r.NatGatewayId
		}
		if r.InstanceId != nil {
			m["instance_id"] = *r.InstanceId
		}
		if r.TransitGatewayId != nil {
			m["transit_gateway_id"] = *r.TransitGatewayId
		}
		if r.VpcPeeringConnectionId != nil {
			m["vpc_peering_connection_id"] = *r.VpcPeeringConnectionId
		}
		if r.NetworkInterfaceId != nil {
			m["network_interface_id"] = *r.NetworkInterfaceId
		}

		routes = append(routes, m)
	}
	return routes
}

func dataSourceAssociationsRead(ec2Assocations []*ec2.RouteTableAssociation) []map[string]interface{} {
	associations := make([]map[string]interface{}, 0, len(ec2Assocations))
	// Loop through the routes and add them to the set
	for _, a := range ec2Assocations {

		m := make(map[string]interface{})
		m["route_table_id"] = *a.RouteTableId
		m["route_table_association_id"] = *a.RouteTableAssociationId
		// GH[11134]
		if a.SubnetId != nil {
			m["subnet_id"] = *a.SubnetId
		}
		if a.GatewayId != nil {
			m["gateway_id"] = *a.GatewayId
		}
		m["main"] = *a.Main
		associations = append(associations, m)
	}
	return associations
}
