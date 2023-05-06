// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package route53

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newResourceCIDRCollection,
		},
		{
			Factory: newResourceCIDRLocation,
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceDelegationSet,
			TypeName: "aws_route53_delegation_set",
		},
		{
			Factory:  DataSourceTrafficPolicyDocument,
			TypeName: "aws_route53_traffic_policy_document",
		},
		{
			Factory:  DataSourceZone,
			TypeName: "aws_route53_zone",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceDelegationSet,
			TypeName: "aws_route53_delegation_set",
		},
		{
			Factory:  ResourceHealthCheck,
			TypeName: "aws_route53_health_check",
			Name:     "Health Check",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "id",
				ResourceType:        "healthcheck",
			},
		},
		{
			Factory:  ResourceHostedZoneDNSSEC,
			TypeName: "aws_route53_hosted_zone_dnssec",
		},
		{
			Factory:  ResourceKeySigningKey,
			TypeName: "aws_route53_key_signing_key",
		},
		{
			Factory:  ResourceQueryLog,
			TypeName: "aws_route53_query_log",
		},
		{
			Factory:  ResourceRecord,
			TypeName: "aws_route53_record",
		},
		{
			Factory:  ResourceTrafficPolicy,
			TypeName: "aws_route53_traffic_policy",
		},
		{
			Factory:  ResourceTrafficPolicyInstance,
			TypeName: "aws_route53_traffic_policy_instance",
		},
		{
			Factory:  ResourceVPCAssociationAuthorization,
			TypeName: "aws_route53_vpc_association_authorization",
		},
		{
			Factory:  ResourceZone,
			TypeName: "aws_route53_zone",
			Name:     "Hosted Zone",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "id",
				ResourceType:        "hostedzone",
			},
		},
		{
			Factory:  ResourceZoneAssociation,
			TypeName: "aws_route53_zone_association",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.Route53
}

var ServicePackage = &servicePackage{}
