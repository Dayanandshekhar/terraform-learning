// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package amplify_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

const serializeDelay = 4*time.Minute + 30*time.Second

// Serialize to limit API rate-limit exceeded errors.
func TestAccAmplify_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"App": {
			acctest.CtBasic:            testAccApp_basic,
			"disappears":               testAccApp_disappears,
			"tags":                     testAccAmplifyApp_tagsSerial,
			"AutoBranchCreationConfig": testAccApp_AutoBranchCreationConfig,
			"BasicAuthCredentials":     testAccApp_BasicAuthCredentials,
			"BuildSpec":                testAccApp_BuildSpec,
			"CustomRules":              testAccApp_CustomRules,
			"Description":              testAccApp_Description,
			"EnvironmentVariables":     testAccApp_EnvironmentVariables,
			"IamServiceRole":           testAccApp_IAMServiceRole,
			"Name":                     testAccApp_Name,
			"Repository":               testAccApp_Repository,
		},
		"BackendEnvironment": {
			acctest.CtBasic:                 testAccBackendEnvironment_basic,
			"disappears":                    testAccBackendEnvironment_disappears,
			"DeploymentArtifacts_StackName": testAccBackendEnvironment_DeploymentArtifacts_StackName,
		},
		"Branch": {
			acctest.CtBasic:        testAccBranch_basic,
			"disappears":           testAccBranch_disappears,
			"tags":                 testAccAmplifyBranch_tagsSerial,
			"BasicAuthCredentials": testAccBranch_BasicAuthCredentials,
			"EnvironmentVariables": testAccBranch_EnvironmentVariables,
			"OptionalArguments":    testAccBranch_OptionalArguments,
		},
		"DomainAssociation": {
			acctest.CtBasic: testAccDomainAssociation_basic,
			"disappears":    testAccDomainAssociation_disappears,
			"update":        testAccDomainAssociation_update,
		},
		"Webhook": {
			acctest.CtBasic: testAccWebhook_basic,
			"disappears":    testAccWebhook_disappears,
			"update":        testAccWebhook_update,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, serializeDelay)
}
