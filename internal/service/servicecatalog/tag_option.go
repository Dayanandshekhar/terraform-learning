// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicecatalog

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_servicecatalog_tag_option")
func ResourceTagOption() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTagOptionCreate,
		ReadWithoutTimeout:   resourceTagOptionRead,
		UpdateWithoutTimeout: resourceTagOptionUpdate,
		DeleteWithoutTimeout: resourceTagOptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(TagOptionReadyTimeout),
			Read:   schema.DefaultTimeout(TagOptionReadTimeout),
			Update: schema.DefaultTimeout(TagOptionUpdateTimeout),
			Delete: schema.DefaultTimeout(TagOptionDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTagOptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogClient(ctx)

	input := &servicecatalog.CreateTagOptionInput{
		Key:   aws.String(d.Get("key").(string)),
		Value: aws.String(d.Get("value").(string)),
	}

	var output *servicecatalog.CreateTagOptionOutput
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var err error

		output, err = conn.CreateTagOption(ctx, input)

		if errs.Contains(err, "profile does not exist") {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateTagOption(ctx, input)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Service Catalog Tag Option: %s", err)
	}

	if output == nil || output.TagOptionDetail == nil {
		return sdkdiag.AppendErrorf(diags, "creating Service Catalog Tag Option: empty response")
	}

	d.SetId(aws.ToString(output.TagOptionDetail.Id))

	// Active is not a field of CreateTagOption but is a field of UpdateTagOption. In order to create an
	// inactive Tag Option, you must create an active one and then update it (but calling this resource's
	// Update will error with ErrCodeDuplicateResourceException because Value is unchanged).
	if v, ok := d.GetOk("active"); !ok {
		_, err = conn.UpdateTagOption(ctx, &servicecatalog.UpdateTagOptionInput{
			Id:     aws.String(d.Id()),
			Active: aws.Bool(v.(bool)),
		})

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "creating Service Catalog Tag Option, updating active (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceTagOptionRead(ctx, d, meta)...)
}

func resourceTagOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogClient(ctx)

	output, err := WaitTagOptionReady(ctx, conn, d.Id(), d.Timeout(schema.TimeoutRead))

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Service Catalog Tag Option (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "describing Service Catalog Tag Option (%s): %s", d.Id(), err)
	}

	if output == nil {
		return sdkdiag.AppendErrorf(diags, "getting Service Catalog Tag Option (%s): empty response", d.Id())
	}

	d.Set("active", output.Active)
	d.Set("key", output.Key)
	d.Set("owner", output.Owner)
	d.Set("value", output.Value)

	return diags
}

func resourceTagOptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogClient(ctx)

	input := &servicecatalog.UpdateTagOptionInput{
		Id: aws.String(d.Id()),
	}

	// UpdateTagOption() is very particular about what it receives. Only fields that change should
	// be included or it will throw servicecatalog.ErrCodeDuplicateResourceException, "already exists"

	if d.HasChange("active") {
		input.Active = aws.Bool(d.Get("active").(bool))
	}

	if d.HasChange("value") {
		input.Value = aws.String(d.Get("value").(string))
	}

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, err := conn.UpdateTagOption(ctx, input)

		if errs.Contains(err, "profile does not exist") {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.UpdateTagOption(ctx, input)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Service Catalog Tag Option (%s): %s", d.Id(), err)
	}

	return append(diags, resourceTagOptionRead(ctx, d, meta)...)
}

func resourceTagOptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogClient(ctx)

	input := &servicecatalog.DeleteTagOptionInput{
		Id: aws.String(d.Id()),
	}

	_, err := conn.DeleteTagOption(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Service Catalog Tag Option (%s): %s", d.Id(), err)
	}

	if err := WaitTagOptionDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Service Catalog Tag Option (%s) to be deleted: %s", d.Id(), err)
	}

	return diags
}
