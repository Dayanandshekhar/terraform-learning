package aws

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/cloudwatch/waiter"
)

func resourceAwsCloudWatchMetricStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsCloudWatchMetricStreamCreate,
		ReadContext:   resourceAwsCloudWatchMetricStreamRead,
		UpdateContext: resourceAwsCloudWatchMetricStreamCreate,
		DeleteContext: resourceAwsCloudWatchMetricStreamDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exclude_filter": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"include_filter"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
			},
			"firehose_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateArn,
			},
			"include_filter": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"exclude_filter"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
			},
			"last_update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validateCloudWatchMetricStreamName,
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validateCloudWatchMetricStreamName,
			},
			"output_format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateArn,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAwsCloudWatchMetricStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).cloudwatchconn

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	params := cloudwatch.PutMetricStreamInput{
		Name:         aws.String(name),
		FirehoseArn:  aws.String(d.Get("firehose_arn").(string)),
		RoleArn:      aws.String(d.Get("role_arn").(string)),
		OutputFormat: aws.String(d.Get("output_format").(string)),
		Tags:         keyvaluetags.New(d.Get("tags").(map[string]interface{})).IgnoreAws().CloudwatchTags(),
	}

	if v, ok := d.GetOk("include_filter"); ok && v.(*schema.Set).Len() > 0 {
		params.IncludeFilters = expandCloudWatchMetricStreamFilters(v.(*schema.Set))
	}

	if v, ok := d.GetOk("exclude_filter"); ok && v.(*schema.Set).Len() > 0 {
		params.ExcludeFilters = expandCloudWatchMetricStreamFilters(v.(*schema.Set))
	}

	log.Printf("[DEBUG] Putting CloudWatch MetricStream: %#v", params)
	_, err := conn.PutMetricStreamWithContext(ctx, &params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Putting metric_stream failed: %s", err))
	}
	d.SetId(name)
	log.Println("[INFO] CloudWatch MetricStream put finished")

	return resourceAwsCloudWatchMetricStreamRead(ctx, d, meta)
}

func resourceAwsCloudWatchMetricStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Reading CloudWatch MetricStream: %s", name)
	conn := meta.(*AWSClient).cloudwatchconn

	params := cloudwatch.GetMetricStreamInput{
		Name: aws.String(d.Id()),
	}

	var err error
	var resp *cloudwatch.GetMetricStreamOutput

	if d.IsNewResource() {
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			resp, err = conn.GetMetricStreamWithContext(ctx, &params)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, cloudwatch.ErrCodeResourceNotFoundException) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if isResourceTimeoutError(err) {
			resp, err = conn.GetMetricStreamWithContext(ctx, &params)
		}
	} else {
		resp, err = conn.GetMetricStreamWithContext(ctx, &params)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, cloudwatch.ErrCodeResourceNotFoundException) {
				log.Printf("[WARN] CloudWatch MetricStream (%s) not found, removing from state", d.Id())
				d.SetId("")
				return nil
			}
		}
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("Reading metric_stream failed: %s", err))
	}

	d.Set("arn", resp.Arn)
	d.Set("creation_date", resp.CreationDate.Format(time.RFC3339))
	d.Set("firehose_arn", resp.FirehoseArn)
	d.Set("last_update_date", resp.CreationDate.Format(time.RFC3339))
	d.Set("name", resp.Name)
	d.Set("output_format", resp.OutputFormat)
	d.Set("role_arn", resp.RoleArn)
	d.Set("state", resp.State)

	if resp.IncludeFilters != nil {
		if err := d.Set("include_filter", flattenCloudWatchMetricStreamFilter(resp.IncludeFilters)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting include_filter error: %w", err))
		}
	}

	if resp.ExcludeFilters != nil {
		if err := d.Set("exclude_filter", flattenCloudWatchMetricStreamFilter(resp.ExcludeFilters)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting exclude_filter error: %w", err))
		}
	}

	return nil
}

func resourceAwsCloudWatchMetricStreamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting CloudWatch MetricStream %s", d.Id())
	conn := meta.(*AWSClient).cloudwatchconn
	params := cloudwatch.DeleteMetricStreamInput{
		Name: aws.String(d.Id()),
	}

	if _, err := conn.DeleteMetricStreamWithContext(ctx, &params); err != nil {
		return diag.FromErr(fmt.Errorf("Error deleting CloudWatch MetricStream: %s", err))
	}

	if _, err := waiter.MetricStreamDeleted(ctx, conn, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("error while waiting for CloudWatch Metric Stream (%s) to become deleted: %w", d.Id(), err))
	}

	log.Printf("[INFO] CloudWatch MetricStream %s deleted", d.Id())

	return nil
}

func validateCloudWatchMetricStreamName(v interface{}, k string) (ws []string, errors []error) {
	return validation.All(
		validation.StringLenBetween(1, 255),
		validation.StringMatch(regexp.MustCompile(`^[\-_A-Za-z0-9]*$`), "must match [\\-_A-Za-z0-9]"),
	)(v, k)
}

func expandCloudWatchMetricStreamFilters(s *schema.Set) []*cloudwatch.MetricStreamFilter {
	var filters []*cloudwatch.MetricStreamFilter

	for _, filterRaw := range s.List() {
		filter := &cloudwatch.MetricStreamFilter{}
		mFilter := filterRaw.(map[string]interface{})

		if v, ok := mFilter["namespace"].(string); ok && v != "" {
			filter.Namespace = aws.String(v)
		}

		filters = append(filters, filter)
	}

	return filters
}

func flattenCloudWatchMetricStreamFilter(s []*cloudwatch.MetricStreamFilter) []map[string]interface{} {
	filters := make([]map[string]interface{}, 0)

	for _, bd := range s {
		if bd.Namespace != nil {
			stage := make(map[string]interface{})
			stage["namespace"] = aws.StringValue(bd.Namespace)

			filters = append(filters, stage)
		}
	}

	if len(filters) > 0 {
		return filters
	}

	return nil
}
