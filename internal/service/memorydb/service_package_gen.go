// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package memorydb

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	memorydb_sdkv1 "github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceACL,
			TypeName: "aws_memorydb_acl",
		},
		{
			Factory:  DataSourceCluster,
			TypeName: "aws_memorydb_cluster",
		},
		{
			Factory:  DataSourceParameterGroup,
			TypeName: "aws_memorydb_parameter_group",
		},
		{
			Factory:  DataSourceSnapshot,
			TypeName: "aws_memorydb_snapshot",
		},
		{
			Factory:  DataSourceSubnetGroup,
			TypeName: "aws_memorydb_subnet_group",
		},
		{
			Factory:  DataSourceUser,
			TypeName: "aws_memorydb_user",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceACL,
			TypeName: "aws_memorydb_acl",
			Name:     "ACL",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  ResourceCluster,
			TypeName: "aws_memorydb_cluster",
			Name:     "Cluster",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  ResourceParameterGroup,
			TypeName: "aws_memorydb_parameter_group",
			Name:     "Parameter Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  ResourceSnapshot,
			TypeName: "aws_memorydb_snapshot",
			Name:     "Snapshot",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  ResourceSubnetGroup,
			TypeName: "aws_memorydb_subnet_group",
			Name:     "Subnet Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  ResourceUser,
			TypeName: "aws_memorydb_user",
			Name:     "User",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.MemoryDB
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*memorydb_sdkv1.MemoryDB, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return memorydb_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
