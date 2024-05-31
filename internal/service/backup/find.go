// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backup

import (
	"context"

	backupv2 "github.com/aws/aws-sdk-go-v2/service/backup"
	awstypes "github.com/aws/aws-sdk-go-v2/service/backup/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindJobByID(ctx context.Context, conn *backup.Backup, id string) (*backup.DescribeBackupJobOutput, error) {
	input := &backup.DescribeBackupJobInput{
		BackupJobId: aws.String(id),
	}

	output, err := conn.DescribeBackupJobWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, backup.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func FindRecoveryPointByTwoPartKey(ctx context.Context, conn *backup.Backup, backupVaultName, recoveryPointARN string) (*backup.DescribeRecoveryPointOutput, error) {
	input := &backup.DescribeRecoveryPointInput{
		BackupVaultName:  aws.String(backupVaultName),
		RecoveryPointArn: aws.String(recoveryPointARN),
	}

	output, err := conn.DescribeRecoveryPointWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, backup.ErrCodeResourceNotFoundException, errCodeAccessDeniedException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func FindVaultAccessPolicyByName(ctx context.Context, conn *backup.Backup, name string) (*backup.GetBackupVaultAccessPolicyOutput, error) {
	input := &backup.GetBackupVaultAccessPolicyInput{
		BackupVaultName: aws.String(name),
	}

	output, err := conn.GetBackupVaultAccessPolicyWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, backup.ErrCodeResourceNotFoundException, errCodeAccessDeniedException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func FindVaultByName(ctx context.Context, conn *backup.Backup, name string) (*backup.DescribeBackupVaultOutput, error) {
	input := &backup.DescribeBackupVaultInput{
		BackupVaultName: aws.String(name),
	}

	output, err := conn.DescribeBackupVaultWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, backup.ErrCodeResourceNotFoundException, errCodeAccessDeniedException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func FindRestoreTestingPlanByName(ctx context.Context, conn *backupv2.Client, name string) (*backupv2.GetRestoreTestingPlanOutput, error) {
	in := &backupv2.GetRestoreTestingPlanInput{
		RestoreTestingPlanName: aws.String(name),
	}

	out, err := conn.GetRestoreTestingPlan(ctx, in)
	if err != nil {
		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.RestoreTestingPlan == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func FindRestoreTestingSelectionByName(ctx context.Context, conn *backupv2.Client, name string, restoreTestingPlanName string) (*backupv2.GetRestoreTestingSelectionOutput, error) {
	in := &backupv2.GetRestoreTestingSelectionInput{
		RestoreTestingPlanName:      aws.String(restoreTestingPlanName),
		RestoreTestingSelectionName: aws.String(name),
	}

	out, err := conn.GetRestoreTestingSelection(ctx, in)
	if err != nil {
		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.RestoreTestingSelection == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}
