// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transfer

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transfer"
	awstypes "github.com/aws/aws-sdk-go-v2/service/transfer/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_transfer_ssh_key")
func ResourceSSHKey() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSSHKeyCreate,
		ReadWithoutTimeout:   resourceSSHKeyRead,
		DeleteWithoutTimeout: resourceSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"body": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					old = cleanSSHKey(old)
					new = cleanSSHKey(new)
					return strings.Trim(old, "\n") == strings.Trim(new, "\n")
				},
			},

			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validServerID,
			},

			names.AttrUserName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validUserName,
			},
		},
	}
}

func resourceSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).TransferClient(ctx)
	userName := d.Get(names.AttrUserName).(string)
	serverID := d.Get("server_id").(string)

	createOpts := &transfer.ImportSshPublicKeyInput{
		ServerId:         aws.String(serverID),
		UserName:         aws.String(userName),
		SshPublicKeyBody: aws.String(d.Get("body").(string)),
	}

	log.Printf("[DEBUG] Create Transfer SSH Public Key Option: %#v", createOpts)

	resp, err := conn.ImportSshPublicKey(ctx, createOpts)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "importing ssh public key: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", serverID, userName, aws.ToString(resp.SshPublicKeyId)))

	return diags
}

func resourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).TransferClient(ctx)
	serverID, userName, sshKeyID, err := DecodeSSHKeyID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parsing Transfer SSH Public Key ID: %s", err)
	}

	descOpts := &transfer.DescribeUserInput{
		UserName: aws.String(userName),
		ServerId: aws.String(serverID),
	}

	resp, err := conn.DescribeUser(ctx, descOpts)
	if err != nil {
		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			log.Printf("[WARN] Transfer User (%s) for Server (%s) not found, removing ssh public key (%s) from state", userName, serverID, sshKeyID)
			d.SetId("")
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "reading Transfer SSH Key (%s): %s", d.Id(), err)
	}

	var body string
	for _, s := range resp.User.SshPublicKeys {
		if sshKeyID == aws.ToString(s.SshPublicKeyId) {
			body = aws.ToString(s.SshPublicKeyBody)
		}
	}

	if body == "" {
		log.Printf("[WARN] No such ssh public key found for User (%s) in Server (%s)", userName, serverID)
		d.SetId("")
	}

	d.Set("server_id", resp.ServerId)
	d.Set(names.AttrUserName, resp.User.UserName)
	d.Set("body", body)

	return diags
}

func resourceSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).TransferClient(ctx)
	serverID, userName, sshKeyID, err := DecodeSSHKeyID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parsing Transfer SSH Public Key ID: %s", err)
	}

	delOpts := &transfer.DeleteSshPublicKeyInput{
		UserName:       aws.String(userName),
		ServerId:       aws.String(serverID),
		SshPublicKeyId: aws.String(sshKeyID),
	}

	_, err = conn.DeleteSshPublicKey(ctx, delOpts)
	if err != nil {
		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "deleting Transfer User SSH Key (%s): %s", d.Id(), err)
	}

	return diags
}

func DecodeSSHKeyID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected SERVERID/USERNAME/SSHKEYID", id)
	}
	return idParts[0], idParts[1], idParts[2], nil
}
func cleanSSHKey(key string) string {
	// Remove comments from SSH Keys
	// Comments are anything after "ssh-rsa XXXX" where XXXX is the key.
	parts := strings.Split(key, " ")
	if len(parts) > 2 {
		parts = parts[0:2]
	}
	return strings.Join(parts, " ")
}
