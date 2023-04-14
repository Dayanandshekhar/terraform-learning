package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	IPAMPreviewNewCIDRIdPartsCount = 2
	ResNameIPAMPreviewNewCIDR      = "IPAM Preview New CIDR"
)

// @SDKResource("aws_vpc_ipam_preview_next_cidr")
func ResourceIPAMPreviewNextCIDR() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIPAMPreviewNextCIDRCreate,
		ReadWithoutTimeout:   resourceIPAMPreviewNextCIDRRead,
		DeleteWithoutTimeout: schema.NoopContext,
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disallowed_cidrs": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						verify.ValidIPv4CIDRNetworkAddress,
						// Follow the numbers used for netmask_length
						validation.IsCIDRNetwork(0, 32),
					),
				},
			},
			"ipam_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"netmask_length": {
				// Possible netmask lengths for IPv4 addresses are 0 - 32.
				// AllocateIpamPoolCidr API
				//   - If there is no DefaultNetmaskLength allocation rule set on the pool,
				//   you must specify either the NetmaskLength or the CIDR.
				//   - If the DefaultNetmaskLength allocation rule is set on the pool,
				//   you can specify either the NetmaskLength or the CIDR and the
				//   DefaultNetmaskLength allocation rule will be ignored.
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 32),
			},
		},
	}
}

func resourceIPAMPreviewNextCIDRCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn()
	poolId := d.Get("ipam_pool_id").(string)

	input := &ec2.AllocateIpamPoolCidrInput{
		ClientToken:     aws.String(id.UniqueId()),
		IpamPoolId:      aws.String(poolId),
		PreviewNextCidr: aws.Bool(true),
	}

	if v, ok := d.GetOk("disallowed_cidrs"); ok && v.(*schema.Set).Len() > 0 {
		input.DisallowedCidrs = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("netmask_length"); ok {
		input.NetmaskLength = aws.Int64(int64(v.(int)))
	}

	output, err := conn.AllocateIpamPoolCidrWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "Error allocating cidr from IPAM pool (%s): %s", d.Get("ipam_pool_id").(string), err)
	}

	if output == nil || output.IpamPoolAllocation == nil {
		return sdkdiag.AppendErrorf(diags, "allocating from ipam pool (%s): empty response", poolId)
	}

	cidr := output.IpamPoolAllocation.Cidr

	d.Set("cidr", cidr)

	// Generate an ID
	idParts := []string{
		aws.StringValue(cidr),
		poolId,
	}

	id, err := flex.FlattenResourceId(idParts, IPAMPreviewNewCIDRIdPartsCount)

	if err != nil {
		return create.DiagError(names.Lightsail, create.ErrActionFlatteningResourceId, ResNameIPAMPreviewNewCIDR, poolId, err)
	}

	d.SetId(id)

	return append(diags, resourceIPAMPreviewNextCIDRRead(ctx, d, meta)...)
}

func resourceIPAMPreviewNextCIDRRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cidr, poolId, err := decodeIPAMPreviewNextCIDRID(d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 IPAM Preview Next CIDR: %s", err)
	}

	d.Set("cidr", cidr)
	d.Set("ipam_pool_id", poolId)

	return diags
}

func decodeIPAMPreviewNextCIDRID(id string) (string, string, error) {
	idParts, err := flex.ExpandResourceId(id, IPAMPreviewNewCIDRIdPartsCount)

	if err != nil {
		return "", "", err
	}

	return idParts[0], idParts[1], nil
}
