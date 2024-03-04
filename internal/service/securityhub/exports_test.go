// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securityhub

// Exports for use in tests only.
var (
	ResourceAutomationRule = newAutomationRuleResource

	FindAutomationRuleByARN = findAutomationRuleByARN
)
