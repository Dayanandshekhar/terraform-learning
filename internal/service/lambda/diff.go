// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lambda

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customizeDiffValidateInput validates that `input` is JSON object when
// `lifecycle_scope` is not "CREATE_ONLY"
func customizeDiffValidateInput(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	if diff.Get("lifecycle_scope") == lifecycleScopeCreateOnly {
		return nil
	}

	if !diff.GetRawPlan().GetAttr("input").IsWhollyKnown() {
		return nil
	}

	// input is validated to be valid JSON in the schema already.
	inputNoSpaces := strings.TrimSpace(diff.Get("input").(string))
	if strings.HasPrefix(inputNoSpaces, "{") && strings.HasSuffix(inputNoSpaces, "}") {
		return nil
	}

	return errors.New(`lifecycle_scope other than "CREATE_ONLY" requires input to be a JSON object`)
}

// customizeDiffInputChangeWithCreateOnlyScope forces a new resource when `input` has
// a change and `lifecycle_scope` is set to "CREATE_ONLY"
func customizeDiffInputChangeWithCreateOnlyScope(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	if diff.HasChange("input") && diff.Get("lifecycle_scope").(string) == lifecycleScopeCreateOnly {
		return diff.ForceNew("input")
	}
	return nil
}

// customizeDiffInputChangedWithCrudScope invalidates the result of the previous invocation for any changes
// to the input, causing any dependent resources to refresh their own state, when `lifecycle_scope` is set to "CRUD"
func customizeDiffInputChangeWithCrudScope(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	if diff.HasChange("input") && diff.Get("lifecycle_scope").(string) == lifecycleScopeCrud {
		return diff.SetNewComputed("result")
	}
	return nil
}
