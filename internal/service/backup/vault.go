// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backup

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	awstypes "github.com/aws/aws-sdk-go-v2/service/backup/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_backup_vault", name="Vault")
// @Tags(identifierAttribute="arn")
func ResourceVault() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVaultCreate,
		ReadWithoutTimeout:   resourceVaultRead,
		UpdateWithoutTimeout: resourceVaultUpdate,
		DeleteWithoutTimeout: resourceVaultDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"kms_key_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(2, 50),
					validation.StringMatch(regexache.MustCompile(`^[0-9A-Za-z_-]*$`), "must consist of letters, numbers, and hyphens."),
				),
			},
			"recovery_points": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceVaultCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BackupClient(ctx)

	name := d.Get("name").(string)
	input := &backup.CreateBackupVaultInput{
		BackupVaultName: aws.String(name),
		BackupVaultTags: getTagsIn(ctx),
	}

	if v, ok := d.GetOk("kms_key_arn"); ok {
		input.EncryptionKeyArn = aws.String(v.(string))
	}

	_, err := conn.CreateBackupVault(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Backup Vault (%s): %s", name, err)
	}

	d.SetId(name)

	return append(diags, resourceVaultRead(ctx, d, meta)...)
}

func resourceVaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BackupClient(ctx)

	output, err := FindVaultByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Backup Vault (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Backup Vault (%s): %s", d.Id(), err)
	}

	d.Set("arn", output.BackupVaultArn)
	d.Set("kms_key_arn", output.EncryptionKeyArn)
	d.Set("name", output.BackupVaultName)
	d.Set("recovery_points", output.NumberOfRecoveryPoints)

	return diags
}

func resourceVaultUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Tags only.

	return append(diags, resourceVaultRead(ctx, d, meta)...)
}

func resourceVaultDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BackupClient(ctx)

	if d.Get("force_destroy").(bool) {
		input := &backup.ListRecoveryPointsByBackupVaultInput{
			BackupVaultName: aws.String(d.Id()),
		}
		var errs []error

		pages := backup.NewListRecoveryPointsByBackupVaultPaginator(conn, input)

		for pages.HasMorePages() {
			page, err := pages.NextPage(ctx)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "listing Backup Vault (%s) recovery points: %s", d.Id(), err)
			}

			if err := errors.Join(errs...); err != nil {
				return sdkdiag.AppendErrorf(diags, "deleting Backup Vault (%s): %s", d.Id(), err)
			}

			for _, v := range page.RecoveryPoints {
				recoveryPointARN := aws.ToString(v.RecoveryPointArn)

				log.Printf("[DEBUG] Deleting Backup Vault recovery point: %s", recoveryPointARN)
				_, err := conn.DeleteRecoveryPoint(ctx, &backup.DeleteRecoveryPointInput{
					BackupVaultName:  aws.String(d.Id()),
					RecoveryPointArn: aws.String(recoveryPointARN),
				})

				if err != nil {
					errs = append(errs, fmt.Errorf("deleting recovery point (%s): %w", recoveryPointARN, err))

					continue
				}

				if _, err := waitRecoveryPointDeleted(ctx, conn, d.Id(), recoveryPointARN, d.Timeout(schema.TimeoutDelete)); err != nil {
					errs = append(errs, fmt.Errorf("waiting for recovery point (%s) delete: %w", recoveryPointARN, err))

					continue
				}
			}
		}
	}

	log.Printf("[DEBUG] Deleting Backup Vault: %s", d.Id())
	_, err := conn.DeleteBackupVault(ctx, &backup.DeleteBackupVaultInput{
		BackupVaultName: aws.String(d.Id()),
	})

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Backup Vault (%s): %s", d.Id(), err)
	}

	return diags
}
