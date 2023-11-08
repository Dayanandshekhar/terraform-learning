// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package securityhub

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	endpoints_sdkv1 "github.com/aws/aws-sdk-go/aws/endpoints"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	securityhub_sdkv1 "github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
	"log"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceAccount,
			TypeName: "aws_securityhub_account",
		},
		{
			Factory:  ResourceActionTarget,
			TypeName: "aws_securityhub_action_target",
		},
		{
			Factory:  ResourceFindingAggregator,
			TypeName: "aws_securityhub_finding_aggregator",
		},
		{
			Factory:  ResourceInsight,
			TypeName: "aws_securityhub_insight",
		},
		{
			Factory:  ResourceInviteAccepter,
			TypeName: "aws_securityhub_invite_accepter",
		},
		{
			Factory:  ResourceMember,
			TypeName: "aws_securityhub_member",
		},
		{
			Factory:  ResourceOrganizationAdminAccount,
			TypeName: "aws_securityhub_organization_admin_account",
		},
		{
			Factory:  ResourceOrganizationConfiguration,
			TypeName: "aws_securityhub_organization_configuration",
		},
		{
			Factory:  ResourceProductSubscription,
			TypeName: "aws_securityhub_product_subscription",
		},
		{
			Factory:  ResourceStandardsControl,
			TypeName: "aws_securityhub_standards_control",
		},
		{
			Factory:  ResourceStandardsSubscription,
			TypeName: "aws_securityhub_standards_subscription",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.SecurityHub
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*securityhub_sdkv1.SecurityHub, error) {
	sess := config["session"].(*session_sdkv1.Session)

	if endpoint := config["endpoint"].(string); endpoint != "" && sess.Config.UseFIPSEndpoint == endpoints_sdkv1.FIPSEndpointStateEnabled {
		// The SDK doesn't allow setting a custom non-FIPS endpoint *and* enabling UseFIPSEndpoint.
		// However there are a few cases where this is necessary; some services don't have FIPS endpoints,
		// and for some services (e.g. CloudFront) the SDK generates the wrong fips endpoint.
		// While forcing this to disabled may result in the end-user not using a FIPS endpoint as specified
		// by setting UseFIPSEndpoint=true in the provider, the user also explicitly changed the endpoint, so
		// here we need to assume the user knows what they're doing.
		log.Printf("[WARN] UseFIPSEndpoint is enabled but a custom endpoint (%s) is configured, ignoring UseFIPSEndpoint.", endpoint)
		sess.Config.UseFIPSEndpoint = endpoints_sdkv1.FIPSEndpointStateDisabled
	}

	return securityhub_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
