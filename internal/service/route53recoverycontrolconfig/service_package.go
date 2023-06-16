package route53recoverycontrolconfig

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	endpoints_sdkv1 "github.com/aws/aws-sdk-go/aws/endpoints"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	route53recoverycontrolconfig_sdkv1 "github.com/aws/aws-sdk-go/service/route53recoverycontrolconfig"
)

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context) (*route53recoverycontrolconfig_sdkv1.Route53RecoveryControlConfig, error) {
	sess := p.config["session"].(*session_sdkv1.Session)
	config := &aws_sdkv1.Config{Endpoint: aws_sdkv1.String(p.config["endpoint"].(string))}

	// Force "global" services to correct Regions.
	if p.config["partition"].(string) == endpoints_sdkv1.AwsPartitionID {
		config.Region = aws_sdkv1.String(endpoints_sdkv1.UsWest2RegionID)
	}

	return route53recoverycontrolconfig_sdkv1.New(sess.Copy(config)), nil
}
