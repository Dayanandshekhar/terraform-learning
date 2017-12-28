package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsApiGatewayStage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsApiGatewayStageCreate,
		Read:   resourceAwsApiGatewayStageRead,
		Update: resourceAwsApiGatewayStageUpdate,
		Delete: resourceAwsApiGatewayStageDelete,

		Schema: map[string]*schema.Schema{
			"cache_cluster_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cache_cluster_size": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"canary_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"percent_traffic": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0.0,
						},
						"stage_variable_overrides": {
							Type:     schema.TypeMap,
							Elem:     schema.TypeString,
							Optional: true,
						},
						"use_stage_cache": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"client_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"documentation_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rest_api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stage_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func appendCanarySettingsPatchOperations(operations []*apigateway.PatchOperation, oldCanarySettingsRaw, newCanarySettingsRaw []interface{}) []*apigateway.PatchOperation {
	if len(newCanarySettingsRaw) == 0 { // Schema guarantees either 0 or 1
		return append(operations, &apigateway.PatchOperation{
			Op:   aws.String("remove"),
			Path: aws.String("/canarySettings"),
		})
	}
	newSettings := newCanarySettingsRaw[0].(map[string]interface{})

	var oldSettings map[string]interface{}
	if len(oldCanarySettingsRaw) == 1 { // Schema guarantees either 0 or 1
		oldSettings = oldCanarySettingsRaw[0].(map[string]interface{})
	} else {
		oldSettings = map[string]interface{}{
			"percent_traffic":          0.0,
			"stage_variable_overrides": make(map[string]interface{}),
			"use_stage_cache":          false,
		}
	}

	oldOverrides := oldSettings["stage_variable_overrides"].(map[string]interface{})
	newOverrides := newSettings["stage_variable_overrides"].(map[string]interface{})
	operations = append(operations, diffVariablesOps("/canarySettings/stageVariableOverrides/", oldOverrides, newOverrides)...)

	oldPercentTraffic := oldSettings["percent_traffic"].(float64)
	newPercentTraffic := newSettings["percent_traffic"].(float64)
	if oldPercentTraffic != newPercentTraffic {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/canarySettings/percentTraffic"),
			Value: aws.String(fmt.Sprintf("%f", newPercentTraffic)),
		})
	}

	oldUseStageCache := oldSettings["use_stage_cache"].(bool)
	newUseStageCache := newSettings["use_stage_cache"].(bool)
	if oldUseStageCache != newUseStageCache {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/canarySettings/useStageCache"),
			Value: aws.String(fmt.Sprintf("%t", newUseStageCache)),
		})
	}

	return operations
}

func readStageVariableOverrides(overrides map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range overrides {
		result[k] = v.(string)
	}

	return result
}

func resourceAwsApiGatewayStageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).apigateway

	d.Partial(true)

	input := apigateway.CreateStageInput{
		RestApiId:    aws.String(d.Get("rest_api_id").(string)),
		StageName:    aws.String(d.Get("stage_name").(string)),
		DeploymentId: aws.String(d.Get("deployment_id").(string)),
	}

	waitForCache := false
	if v, ok := d.GetOk("cache_cluster_enabled"); ok {
		input.CacheClusterEnabled = aws.Bool(v.(bool))
		waitForCache = true
	}
	if v, ok := d.GetOk("cache_cluster_size"); ok {
		input.CacheClusterSize = aws.String(v.(string))
		waitForCache = true
	}
	if v, ok := d.GetOk("canary_settings"); ok {
		canarySettings := v.([]interface{})
		canarySetting, ok := canarySettings[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("At least one field is expected inside canary_settings")
		}

		var stageVariableOverrides map[string]string
		if canarySettingVariables, ok := canarySetting["stage_variable_overrides"]; ok {
			stageVariableOverrides = readStageVariableOverrides(canarySettingVariables.(map[string]interface{}))
		} else {
			return fmt.Errorf("stage_variable_overrides must be set in canary_settings")
		}

		percentTraffic, ok := canarySetting["percent_traffic"]
		if !ok {
			return fmt.Errorf("percent_traffic must be set in canary_settings")
		}

		useStageCache, ok := canarySetting["use_stage_cache"]
		if !ok {
			return fmt.Errorf("use_stage_cache must be set in canary_settings")
		}

		input.CanarySettings = &apigateway.CanarySettings{
			DeploymentId:           aws.String(d.Get("deployment_id").(string)),
			StageVariableOverrides: aws.StringMap(stageVariableOverrides),
			PercentTraffic:         aws.Float64(percentTraffic.(float64)),
			UseStageCache:          aws.Bool(useStageCache.(bool)),
		}
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}
	if v, ok := d.GetOk("documentation_version"); ok {
		input.DocumentationVersion = aws.String(v.(string))
	}
	if vars, ok := d.GetOk("variables"); ok {
		variables := make(map[string]string, 0)
		for k, v := range vars.(map[string]interface{}) {
			variables[k] = v.(string)
		}
		input.Variables = aws.StringMap(variables)
	}

	out, err := conn.CreateStage(&input)
	if err != nil {
		return fmt.Errorf("Error creating API Gateway Stage: %s", err)
	}

	d.SetId(fmt.Sprintf("ags-%s-%s", d.Get("rest_api_id").(string), d.Get("stage_name").(string)))

	d.SetPartial("rest_api_id")
	d.SetPartial("stage_name")
	d.SetPartial("deployment_id")
	d.SetPartial("description")
	d.SetPartial("variables")

	if waitForCache && *out.CacheClusterStatus != "NOT_AVAILABLE" {
		stateConf := &resource.StateChangeConf{
			Pending: []string{
				"CREATE_IN_PROGRESS",
				"DELETE_IN_PROGRESS",
				"FLUSH_IN_PROGRESS",
			},
			Target: []string{"AVAILABLE"},
			Refresh: apiGatewayStageCacheRefreshFunc(conn,
				d.Get("rest_api_id").(string),
				d.Get("stage_name").(string)),
			Timeout: 90 * time.Minute,
		}

		_, err := stateConf.WaitForState()
		if err != nil {
			return err
		}
	}

	d.SetPartial("canary_settings")
	d.SetPartial("cache_cluster_enabled")
	d.SetPartial("cache_cluster_size")
	d.Partial(false)

	if _, ok := d.GetOk("client_certificate_id"); ok {
		return resourceAwsApiGatewayStageUpdate(d, meta)
	}
	return resourceAwsApiGatewayStageRead(d, meta)
}

func resourceAwsApiGatewayStageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).apigateway

	log.Printf("[DEBUG] Reading API Gateway Stage %s", d.Id())
	input := apigateway.GetStageInput{
		RestApiId: aws.String(d.Get("rest_api_id").(string)),
		StageName: aws.String(d.Get("stage_name").(string)),
	}
	stage, err := conn.GetStage(&input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NotFoundException" {
			log.Printf("[WARN] API Gateway Stage (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	log.Printf("[DEBUG] Received API Gateway Stage: %s", stage)

	d.Set("client_certificate_id", stage.ClientCertificateId)

	if stage.CacheClusterStatus != nil && *stage.CacheClusterStatus == "DELETE_IN_PROGRESS" {
		d.Set("cache_cluster_enabled", false)
		d.Set("cache_cluster_size", nil)
	} else {
		d.Set("cache_cluster_enabled", stage.CacheClusterEnabled)
		d.Set("cache_cluster_size", stage.CacheClusterSize)
	}

	if err := d.Set("canary_settings", flattenApiGatewayStageCanarySettings(stage.CanarySettings)); err != nil {
		log.Printf("[ERR] Error setting canary settings for api gateway stage (%s): %s", d.Id(), err)
	}

	d.Set("deployment_id", stage.DeploymentId)
	d.Set("description", stage.Description)
	d.Set("documentation_version", stage.DocumentationVersion)
	d.Set("variables", aws.StringValueMap(stage.Variables))

	return nil
}

func resourceAwsApiGatewayStageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).apigateway

	d.Partial(true)
	operations := make([]*apigateway.PatchOperation, 0)
	waitForCache := false
	if d.HasChange("cache_cluster_enabled") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/cacheClusterEnabled"),
			Value: aws.String(fmt.Sprintf("%t", d.Get("cache_cluster_enabled").(bool))),
		})
		waitForCache = true
	}
	if d.HasChange("cache_cluster_size") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/cacheClusterSize"),
			Value: aws.String(d.Get("cache_cluster_size").(string)),
		})
		waitForCache = true
	}
	if d.HasChange("canary_settings") {
		oldCanarySettingsRaw, newCanarySettingsRaw := d.GetChange("canary_settings")
		operations = appendCanarySettingsPatchOperations(operations, oldCanarySettingsRaw.([]interface{}), newCanarySettingsRaw.([]interface{}))
	}
	if d.HasChange("client_certificate_id") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/clientCertificateId"),
			Value: aws.String(d.Get("client_certificate_id").(string)),
		})
	}
	if d.HasChange("deployment_id") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/deploymentId"),
			Value: aws.String(d.Get("deployment_id").(string)),
		})

		if _, ok := d.GetOk("canary_settings"); ok {
			operations = append(operations, &apigateway.PatchOperation{
				Op:    aws.String("replace"),
				Path:  aws.String("/canarySettings/deploymentId"),
				Value: aws.String(d.Get("deployment_id").(string)),
			})
		}
	}
	if d.HasChange("description") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/description"),
			Value: aws.String(d.Get("description").(string)),
		})
	}
	if d.HasChange("documentation_version") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/documentationVersion"),
			Value: aws.String(d.Get("documentation_version").(string)),
		})
	}
	if d.HasChange("variables") {
		o, n := d.GetChange("variables")
		oldV := o.(map[string]interface{})
		newV := n.(map[string]interface{})
		operations = append(operations, diffVariablesOps("/variables/", oldV, newV)...)
	}

	input := apigateway.UpdateStageInput{
		RestApiId:       aws.String(d.Get("rest_api_id").(string)),
		StageName:       aws.String(d.Get("stage_name").(string)),
		PatchOperations: operations,
	}
	log.Printf("[DEBUG] Updating API Gateway Stage: %s", input)
	out, err := conn.UpdateStage(&input)
	if err != nil {
		return fmt.Errorf("Updating API Gateway Stage failed: %s", err)
	}

	d.SetPartial("client_certificate_id")
	d.SetPartial("deployment_id")
	d.SetPartial("description")
	d.SetPartial("variables")

	if waitForCache && *out.CacheClusterStatus != "NOT_AVAILABLE" {
		stateConf := &resource.StateChangeConf{
			Pending: []string{
				"CREATE_IN_PROGRESS",
				"FLUSH_IN_PROGRESS",
			},
			Target: []string{
				"AVAILABLE",
				// There's an AWS API bug (raised & confirmed in Sep 2016 by support)
				// which causes the stage to remain in deletion state forever
				"DELETE_IN_PROGRESS",
			},
			Refresh: apiGatewayStageCacheRefreshFunc(conn,
				d.Get("rest_api_id").(string),
				d.Get("stage_name").(string)),
			Timeout: 30 * time.Minute,
		}

		_, err := stateConf.WaitForState()
		if err != nil {
			return err
		}
	}

	d.SetPartial("canary_settings")
	d.SetPartial("cache_cluster_enabled")
	d.SetPartial("cache_cluster_size")
	d.Partial(false)

	return resourceAwsApiGatewayStageRead(d, meta)
}

func diffVariablesOps(prefix string, oldVars, newVars map[string]interface{}) []*apigateway.PatchOperation {
	ops := make([]*apigateway.PatchOperation, 0)

	for k, _ := range oldVars {
		if _, ok := newVars[k]; !ok {
			ops = append(ops, &apigateway.PatchOperation{
				Op:   aws.String("remove"),
				Path: aws.String(prefix + k),
			})
		}
	}

	for k, v := range newVars {
		newValue := v.(string)

		if oldV, ok := oldVars[k]; ok {
			oldValue := oldV.(string)
			if oldValue == newValue {
				continue
			}
		}
		ops = append(ops, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String(prefix + k),
			Value: aws.String(newValue),
		})
	}

	return ops
}

func apiGatewayStageCacheRefreshFunc(conn *apigateway.APIGateway, apiId, stageName string) func() (interface{}, string, error) {
	return func() (interface{}, string, error) {
		input := apigateway.GetStageInput{
			RestApiId: aws.String(apiId),
			StageName: aws.String(stageName),
		}
		out, err := conn.GetStage(&input)
		if err != nil {
			return 42, "", err
		}

		return out, *out.CacheClusterStatus, nil
	}
}

func resourceAwsApiGatewayStageDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).apigateway
	log.Printf("[DEBUG] Deleting API Gateway Stage: %s", d.Id())
	input := apigateway.DeleteStageInput{
		RestApiId: aws.String(d.Get("rest_api_id").(string)),
		StageName: aws.String(d.Get("stage_name").(string)),
	}
	_, err := conn.DeleteStage(&input)
	if err != nil {
		return fmt.Errorf("Deleting API Gateway Stage failed: %s", err)
	}

	return nil
}
