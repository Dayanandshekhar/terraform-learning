package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/chime"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAwsChimeVoiceConnectorGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsChimeVoiceConnectorGroupCreate,
		ReadContext:   resourceAwsChimeVoiceConnectorGroupRead,
		UpdateContext: resourceAwsChimeVoiceConnectorGroupUpdate,
		DeleteContext: resourceAwsChimeVoiceConnectorGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"connector": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 3,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connector_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 99),
						},
					},
				},
			},
		},
	}
}

func resourceAwsChimeVoiceConnectorGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).chimeconn

	input := &chime.CreateVoiceConnectorGroupInput{
		Name: aws.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("connector"); ok && len(v.([]interface{})) > 0 {
		input.VoiceConnectorItems = expandVoiceConnectorItems(v.([]interface{}))
	}

	resp, err := conn.CreateVoiceConnectorGroupWithContext(ctx, input)
	if err != nil {
		return diag.Errorf("error creating Chime Voice connector group: %s", err)
	}

	d.SetId(aws.StringValue(resp.VoiceConnectorGroup.VoiceConnectorGroupId))

	return resourceAwsChimeVoiceConnectorGroupRead(ctx, d, meta)
}

func resourceAwsChimeVoiceConnectorGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).chimeconn

	getInput := &chime.GetVoiceConnectorGroupInput{
		VoiceConnectorGroupId: aws.String(d.Id()),
	}

	resp, err := conn.GetVoiceConnectorGroupWithContext(ctx, getInput)
	if isAWSErr(err, chime.ErrCodeNotFoundException, "") {
		log.Printf("[WARN] Chime Voice conector group %s not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("error getting Voice connector group (%s): %s", d.Id(), err)
	}

	d.Set("name", resp.VoiceConnectorGroup.Name)

	return nil
}

func resourceAwsChimeVoiceConnectorGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).chimeconn

	input := &chime.UpdateVoiceConnectorGroupInput{
		Name:                  aws.String(d.Get("name").(string)),
		VoiceConnectorGroupId: aws.String(d.Id()),
	}

	if d.HasChange("connector") {
		if v, ok := d.GetOk("connector"); ok {
			input.VoiceConnectorItems = expandVoiceConnectorItems(v.([]interface{}))
		}
	} else if !d.IsNewResource() {
		input.VoiceConnectorItems = make([]*chime.VoiceConnectorItem, 0)
	}

	if _, err := conn.UpdateVoiceConnectorGroupWithContext(ctx, input); err != nil {
		if isAWSErr(err, chime.ErrCodeNotFoundException, "") {
			log.Printf("[WARN] Chime Voice conector group %s not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error updating Voice connector group (%s): %s", d.Id(), err)
	}

	return resourceAwsChimeVoiceConnectorGroupRead(ctx, d, meta)
}

func resourceAwsChimeVoiceConnectorGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).chimeconn

	if v, ok := d.GetOk("connector"); ok && len(v.([]interface{})) > 0 {
		if err := resourceAwsChimeVoiceConnectorGroupUpdate(ctx, d, meta); err != nil {
			return err
		}
	}

	input := &chime.DeleteVoiceConnectorGroupInput{
		VoiceConnectorGroupId: aws.String(d.Id()),
	}

	if _, err := conn.DeleteVoiceConnectorGroupWithContext(ctx, input); err != nil {
		if isAWSErr(err, chime.ErrCodeNotFoundException, "") {
			log.Printf("[WARN] Chime Voice conector group %s not found", d.Id())
			return nil
		}
		return diag.Errorf("error deleting Voice connector group (%s): %s", d.Id(), err)
	}

	return nil
}

func expandVoiceConnectorItems(data []interface{}) []*chime.VoiceConnectorItem {
	var connectorsItems []*chime.VoiceConnectorItem

	for _, rItem := range data {
		item := rItem.(map[string]interface{})
		connectorsItems = append(connectorsItems, &chime.VoiceConnectorItem{
			VoiceConnectorId: aws.String(item["connector_id"].(string)),
			Priority:         aws.Int64(int64(item["priority"].(int))),
		})
	}

	return connectorsItems
}
