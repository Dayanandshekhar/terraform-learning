package aws

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsIotTopicRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsIotTopicRuleCreate,
		Read:   resourceAwsIotTopicRuleRead,
		Update: resourceAwsIotTopicRuleUpdate,
		Delete: resourceAwsIotTopicRuleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, s string) ([]string, []error) {
					name := v.(string)
					if len(name) < 1 || len(name) > 128 {
						return nil, []error{fmt.Errorf("Name must between 1 and 128 characters long")}
					}

					matched, err := regexp.MatchReader("^[a-zA-Z0-9_]+$", strings.NewReader(name))

					if err != nil {
						return nil, []error{err}
					}

					if !matched {
						return nil, []error{fmt.Errorf("Name must match the pattern ^[a-zA-Z0-9_]+$")}
					}

					return nil, nil
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sql": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sql_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloudwatch_alarm": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"state_reason": {
							Type:     schema.TypeString,
							Required: true,
						},
						"state_value": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, s string) ([]string, []error) {
								switch v.(string) {
								case
									"OK",
									"ALARM",
									"INSUFFICIENT_DATA":
									return nil, nil
								}

								return nil, []error{fmt.Errorf("State must be one of OK, ALARM, or INSUFFICIENT_DATA")}
							},
						},
					},
				},
			},
			"cloudwatch_metric": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metric_namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metric_timestamp": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, s string) ([]string, []error) {
								dateString := v.(string)
								if _, err := time.Parse(time.RFC3339, dateString); err != nil {
									return nil, []error{err}
								}
								return nil, nil
							},
						},
						"metric_unit": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metric_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"dynamodb": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hash_key_field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"hash_key_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"hash_key_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"payload_field": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"range_key_field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"range_key_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"range_key_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"elasticsearch": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"index": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"firehose": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delivery_stream_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"kinesis": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partition_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"stream_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"lambda": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"republish": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"s3": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"sns": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"sqs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"use_base64": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createTopicRulePayload(d *schema.ResourceData) *iot.TopicRulePayload {
	cloudwatchAlarmActions := d.Get("cloudwatch_alarm").(*schema.Set).List()
	cloudwatchMetricActions := d.Get("cloudwatch_metric").(*schema.Set).List()
	dynamoDbActions := d.Get("dynamodb").(*schema.Set).List()
	elasticsearchActions := d.Get("elasticsearch").(*schema.Set).List()
	firehoseActions := d.Get("firehose").(*schema.Set).List()
	kinesisActions := d.Get("kinesis").(*schema.Set).List()
	lambdaActions := d.Get("lambda").(*schema.Set).List()
	republishActions := d.Get("republish").(*schema.Set).List()
	s3Actions := d.Get("s3").(*schema.Set).List()
	snsActions := d.Get("sns").(*schema.Set).List()
	sqsActions := d.Get("sqs").(*schema.Set).List()

	numActions := len(cloudwatchAlarmActions) + len(cloudwatchMetricActions) +
		len(dynamoDbActions) + len(elasticsearchActions) + len(firehoseActions) +
		len(kinesisActions) + len(lambdaActions) + len(republishActions) +
		len(s3Actions) + len(snsActions) + len(sqsActions)
	actions := make([]*iot.Action, numActions)

	i := 0
	// Add Cloudwatch Alarm actions
	for _, a := range cloudwatchAlarmActions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			CloudwatchAlarm: &iot.CloudwatchAlarmAction{
				AlarmName:   aws.String(raw["alarm_name"].(string)),
				RoleArn:     aws.String(raw["role_arn"].(string)),
				StateReason: aws.String(raw["state_reason"].(string)),
				StateValue:  aws.String(raw["state_value"].(string)),
			},
		}
		i++
	}

	// Add Cloudwatch Metric actions
	for _, a := range cloudwatchMetricActions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			CloudwatchMetric: &iot.CloudwatchMetricAction{
				MetricName:      aws.String(raw["metric_name"].(string)),
				MetricNamespace: aws.String(raw["metric_namespace"].(string)),
				MetricUnit:      aws.String(raw["metric_unit"].(string)),
				MetricValue:     aws.String(raw["metric_value"].(string)),
				RoleArn:         aws.String(raw["role_arn"].(string)),
				MetricTimestamp: aws.String(raw["metric_timestamp"].(string)),
			},
		}
		i++
	}

	// Add DynamoDB actions
	for _, a := range dynamoDbActions {
		raw := a.(map[string]interface{})
		act := &iot.Action{
			DynamoDB: &iot.DynamoDBAction{
				HashKeyField:  aws.String(raw["hash_key_field"].(string)),
				HashKeyValue:  aws.String(raw["hash_key_value"].(string)),
				RangeKeyField: aws.String(raw["range_key_field"].(string)),
				RangeKeyValue: aws.String(raw["range_key_value"].(string)),
				RoleArn:       aws.String(raw["role_arn"].(string)),
				TableName:     aws.String(raw["table_name"].(string)),
			},
		}
		if hkt, ok := raw["hash_key_type"].(string); ok {
			act.DynamoDB.HashKeyType = aws.String(hkt)
		}
		if rkt, ok := raw["range_key_type"].(string); ok {
			act.DynamoDB.RangeKeyType = aws.String(rkt)
		}
		if plf, ok := raw["payload_field"].(string); ok {
			act.DynamoDB.PayloadField = aws.String(plf)
		}
		actions[i] = act
		i++
	}

	// Add Elasticsearch actions

	for _, a := range elasticsearchActions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			Elasticsearch: &iot.ElasticsearchAction{
				Endpoint: aws.String(raw["endpoint"].(string)),
				Id:       aws.String(raw["id"].(string)),
				Index:    aws.String(raw["index"].(string)),
				RoleArn:  aws.String(raw["role_arn"].(string)),
				Type:     aws.String(raw["type"].(string)),
			},
		}
		i++
	}

	// Add Firehose actions

	for _, a := range firehoseActions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			Firehose: &iot.FirehoseAction{
				DeliveryStreamName: aws.String(raw["delivery_stream_name"].(string)),
				RoleArn:            aws.String(raw["role_arn"].(string)),
			},
		}
		i++
	}

	// Add Kinesis actions

	for _, a := range kinesisActions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			Kinesis: &iot.KinesisAction{
				RoleArn:      aws.String(raw["role_arn"].(string)),
				StreamName:   aws.String(raw["stream_name"].(string)),
				PartitionKey: aws.String(raw["partition_key"].(string)),
			},
		}
		i++
	}

	// Add Lambda actions

	for _, a := range lambdaActions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			Lambda: &iot.LambdaAction{
				FunctionArn: aws.String(raw["function_arn"].(string)),
			},
		}
		i++
	}

	// Add Republish actions

	for _, a := range republishActions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			Republish: &iot.RepublishAction{
				RoleArn: aws.String(raw["role_arn"].(string)),
				Topic:   aws.String(raw["topic"].(string)),
			},
		}
		i++
	}

	// Add S3 actions

	for _, a := range s3Actions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			S3: &iot.S3Action{
				BucketName: aws.String(raw["bucket_name"].(string)),
				Key:        aws.String(raw["key"].(string)),
				RoleArn:    aws.String(raw["role_arn"].(string)),
			},
		}
		i++
	}

	// Add SNS actions

	for _, a := range snsActions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			Sns: &iot.SnsAction{
				RoleArn:       aws.String(raw["role_arn"].(string)),
				TargetArn:     aws.String(raw["target_arn"].(string)),
				MessageFormat: aws.String(raw["message_format"].(string)),
			},
		}
		i++
	}

	// Add SQS actions

	for _, a := range sqsActions {
		raw := a.(map[string]interface{})
		actions[i] = &iot.Action{
			Sqs: &iot.SqsAction{
				QueueUrl:  aws.String(raw["queue_url"].(string)),
				RoleArn:   aws.String(raw["role_arn"].(string)),
				UseBase64: aws.Bool(raw["use_base64"].(bool)),
			},
		}
		i++
	}

	return &iot.TopicRulePayload{
		Description:      aws.String(d.Get("description").(string)),
		RuleDisabled:     aws.Bool(!d.Get("enabled").(bool)),
		Sql:              aws.String(d.Get("sql").(string)),
		AwsIotSqlVersion: aws.String(d.Get("sql_version").(string)),
		Actions:          actions,
	}
}

func resourceAwsIotTopicRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn

	ruleName := d.Get("name").(string)

	_, err := conn.CreateTopicRule(&iot.CreateTopicRuleInput{
		RuleName:         aws.String(ruleName),
		TopicRulePayload: createTopicRulePayload(d),
	})

	if err != nil {
		return err
	}

	d.SetId(ruleName)

	return resourceAwsIotTopicRuleRead(d, meta)
}

func resourceAwsIotTopicRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn

	out, err := conn.GetTopicRule(&iot.GetTopicRuleInput{
		RuleName: aws.String(d.Id()),
	})

	if err != nil {
		return err
	}

	d.Set("arn", out.RuleArn)
	d.Set("name", out.Rule.RuleName)
	d.Set("enabled", !(*out.Rule.RuleDisabled))
	d.Set("sql", out.Rule.Sql)

	return nil
}

func resourceAwsIotTopicRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn

	_, err := conn.ReplaceTopicRule(&iot.ReplaceTopicRuleInput{
		RuleName:         aws.String(d.Get("name").(string)),
		TopicRulePayload: createTopicRulePayload(d),
	})

	if err != nil {
		return err
	}

	return resourceAwsIotTopicRuleRead(d, meta)
}

func resourceAwsIotTopicRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn

	_, err := conn.DeleteTopicRule(&iot.DeleteTopicRuleInput{
		RuleName: aws.String(d.Id()),
	})

	if err != nil {
		return err
	}

	return nil
}
