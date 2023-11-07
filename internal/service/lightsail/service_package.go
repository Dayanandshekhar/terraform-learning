// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lightsail

import (
	"context"
	"log"
	"strings"
	"time"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	retry_sdkv2 "github.com/aws/aws-sdk-go-v2/aws/retry"
	lightsail_sdkv2 "github.com/aws/aws-sdk-go-v2/service/lightsail"
)

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*lightsail_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return lightsail_sdkv2.NewFromConfig(cfg, func(o *lightsail_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)

			if o.EndpointOptions.UseFIPSEndpoint == aws_sdkv2.FIPSEndpointStateEnabled {
				// The SDK doesn't allow setting a custom non-FIPS endpoint *and* enabling UseFIPSEndpoint.
				// However there are a few cases where this is necessary; some services don't have FIPS endpoints,
				// and for some services (e.g. CloudFront) the SDK generates the wrong fips endpoint.
				// While forcing this to disabled may result in the end-user not using a FIPS endpoint as specified
				// by setting UseFIPSEndpoint=true, the user also explicitly changed the endpoint, so
				// here we need to assume the user knows what they're doing.
				log.Printf("[WARN] UseFIPSEndpoint is enabled but a custom endpoint (%s) is configured, ignoring UseFIPSEndpoint.", endpoint)
				o.EndpointOptions.UseFIPSEndpoint = aws_sdkv2.FIPSEndpointStateDisabled
			}
		}

		retryable := retry_sdkv2.IsErrorRetryableFunc(func(e error) aws_sdkv2.Ternary {
			if strings.Contains(e.Error(), "Please try again in a few minutes") || strings.Contains(e.Error(), "Please wait for it to complete before trying again") {
				return aws_sdkv2.TrueTernary
			}
			return aws_sdkv2.UnknownTernary
		})
		const (
			backoff = 10 * time.Second
		)
		o.Retryer = retry_sdkv2.NewStandard(func(o *retry_sdkv2.StandardOptions) {
			o.Retryables = append(o.Retryables, retryable)
			o.MaxAttempts = 18
			o.Backoff = retry_sdkv2.NewExponentialJitterBackoff(backoff)
		})
	}), nil
}
