// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package sagemaker

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	endpoints_sdkv1 "github.com/aws/aws-sdk-go/aws/endpoints"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	sagemaker_sdkv1 "github.com/aws/aws-sdk-go/service/sagemaker"
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
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourcePrebuiltECRImage,
			TypeName: "aws_sagemaker_prebuilt_ecr_image",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceApp,
			TypeName: "aws_sagemaker_app",
			Name:     "App",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceAppImageConfig,
			TypeName: "aws_sagemaker_app_image_config",
			Name:     "App Image Config",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceCodeRepository,
			TypeName: "aws_sagemaker_code_repository",
			Name:     "Code Repository",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceDataQualityJobDefinition,
			TypeName: "aws_sagemaker_data_quality_job_definition",
			Name:     "Data Quality Job Definition",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceDevice,
			TypeName: "aws_sagemaker_device",
		},
		{
			Factory:  ResourceDeviceFleet,
			TypeName: "aws_sagemaker_device_fleet",
			Name:     "Device Fleet",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceDomain,
			TypeName: "aws_sagemaker_domain",
			Name:     "Domain",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceEndpoint,
			TypeName: "aws_sagemaker_endpoint",
			Name:     "Endpoint",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceEndpointConfiguration,
			TypeName: "aws_sagemaker_endpoint_configuration",
			Name:     "Endpoint Configuration",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceFeatureGroup,
			TypeName: "aws_sagemaker_feature_group",
			Name:     "Feature Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceFlowDefinition,
			TypeName: "aws_sagemaker_flow_definition",
			Name:     "Flow Definition",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceHumanTaskUI,
			TypeName: "aws_sagemaker_human_task_ui",
			Name:     "Human Task UI",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceImage,
			TypeName: "aws_sagemaker_image",
			Name:     "Image",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceImageVersion,
			TypeName: "aws_sagemaker_image_version",
		},
		{
			Factory:  ResourceModel,
			TypeName: "aws_sagemaker_model",
			Name:     "Model",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceModelPackageGroup,
			TypeName: "aws_sagemaker_model_package_group",
			Name:     "Model Package Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceModelPackageGroupPolicy,
			TypeName: "aws_sagemaker_model_package_group_policy",
		},
		{
			Factory:  ResourceMonitoringSchedule,
			TypeName: "aws_sagemaker_monitoring_schedule",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceNotebookInstance,
			TypeName: "aws_sagemaker_notebook_instance",
			Name:     "Notebook Instance",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceNotebookInstanceLifeCycleConfiguration,
			TypeName: "aws_sagemaker_notebook_instance_lifecycle_configuration",
		},
		{
			Factory:  ResourcePipeline,
			TypeName: "aws_sagemaker_pipeline",
			Name:     "Pipeline",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceProject,
			TypeName: "aws_sagemaker_project",
			Name:     "Project",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceServicecatalogPortfolioStatus,
			TypeName: "aws_sagemaker_servicecatalog_portfolio_status",
		},
		{
			Factory:  ResourceSpace,
			TypeName: "aws_sagemaker_space",
			Name:     "Space",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceStudioLifecycleConfig,
			TypeName: "aws_sagemaker_studio_lifecycle_config",
			Name:     "Studio Lifecycle Config",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceUserProfile,
			TypeName: "aws_sagemaker_user_profile",
			Name:     "User Profile",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceWorkforce,
			TypeName: "aws_sagemaker_workforce",
		},
		{
			Factory:  ResourceWorkteam,
			TypeName: "aws_sagemaker_workteam",
			Name:     "Workteam",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.SageMaker
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*sagemaker_sdkv1.SageMaker, error) {
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

	return sagemaker_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
