package sts

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	sts_sdkv1 "github.com/aws/aws-sdk-go/service/sts"
)

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context) (*sts_sdkv1.STS, error) {
	sess := p.config["session"].(*session_sdkv1.Session)
	config := &aws_sdkv1.Config{Endpoint: aws_sdkv1.String(p.config["endpoint"].(string))}

	if stsRegion := p.config["sts_region"].(string); stsRegion != "" {
		config.Region = aws_sdkv1.String(stsRegion)
	}

	return sts_sdkv1.New(sess.Copy(config)), nil
}
