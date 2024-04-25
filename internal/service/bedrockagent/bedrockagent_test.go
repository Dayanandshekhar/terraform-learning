// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bedrockagent_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccBedrockAgent_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"KnowledgeBase": {
			"basic":      testAccKnowledgeBase_basic,
			"disappears": testAccKnowledgeBase_disappears,
			"rds":        testAccKnowledgeBase_rds,
			"update":     testAccKnowledgeBase_update,
			"tags":       testAccKnowledgeBase_tags,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}
