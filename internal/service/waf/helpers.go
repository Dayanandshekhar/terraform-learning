// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package waf

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/waf/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func SizeConstraintSetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		names.AttrARN: {
			Type:     schema.TypeString,
			Computed: true,
		},
		names.AttrName: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"size_constraints": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"comparison_operator": {
						Type:     schema.TypeString,
						Required: true,
					},
					"field_to_match": {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"data": {
									Type:     schema.TypeString,
									Optional: true,
								},
								names.AttrType: {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"size": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"text_transformation": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func DiffSizeConstraints(oldS, newS []interface{}) []awstypes.SizeConstraintSetUpdate {
	updates := make([]awstypes.SizeConstraintSetUpdate, 0)

	for _, os := range oldS {
		constraint := os.(map[string]interface{})

		if idx, contains := sliceContainsMap(newS, constraint); contains {
			newS = append(newS[:idx], newS[idx+1:]...)
			continue
		}

		updates = append(updates, awstypes.SizeConstraintSetUpdate{
			Action: awstypes.ChangeActionDelete,
			SizeConstraint: &awstypes.SizeConstraint{
				FieldToMatch:       expandFieldToMatch(constraint["field_to_match"].([]interface{})[0].(map[string]interface{})),
				ComparisonOperator: awstypes.ComparisonOperator(constraint["comparison_operator"].(string)),
				Size:               int64(constraint["size"].(int)),
				TextTransformation: awstypes.TextTransformation(constraint["text_transformation"].(string)),
			},
		})
	}

	for _, ns := range newS {
		constraint := ns.(map[string]interface{})

		updates = append(updates, awstypes.SizeConstraintSetUpdate{
			Action: awstypes.ChangeActionInsert,
			SizeConstraint: &awstypes.SizeConstraint{
				FieldToMatch:       expandFieldToMatch(constraint["field_to_match"].([]interface{})[0].(map[string]interface{})),
				ComparisonOperator: awstypes.ComparisonOperator(constraint["comparison_operator"].(string)),
				Size:               int64(constraint["size"].(int)),
				TextTransformation: awstypes.TextTransformation(constraint["text_transformation"].(string)),
			},
		})
	}
	return updates
}

func FlattenSizeConstraints(sc []awstypes.SizeConstraint) []interface{} {
	out := make([]interface{}, len(sc))
	for i, c := range sc {
		m := make(map[string]interface{})
		m["comparison_operator"] = c.ComparisonOperator
		if c.FieldToMatch != nil {
			m["field_to_match"] = flattenFieldToMatch(c.FieldToMatch)
		}
		m["size"] = c.Size
		m["text_transformation"] = c.TextTransformation
		out[i] = m
	}
	return out
}

func DiffRulePredicates(oldP, newP []interface{}) []awstypes.RuleUpdate {
	updates := make([]awstypes.RuleUpdate, 0)

	for _, op := range oldP {
		predicate := op.(map[string]interface{})

		if idx, contains := sliceContainsMap(newP, predicate); contains {
			newP = append(newP[:idx], newP[idx+1:]...)
			continue
		}

		updates = append(updates, awstypes.RuleUpdate{
			Action: awstypes.ChangeActionDelete,
			Predicate: &awstypes.Predicate{
				Negated: aws.Bool(predicate["negated"].(bool)),
				Type:    awstypes.PredicateType(predicate[names.AttrType].(string)),
				DataId:  aws.String(predicate["data_id"].(string)),
			},
		})
	}

	for _, np := range newP {
		predicate := np.(map[string]interface{})

		updates = append(updates, awstypes.RuleUpdate{
			Action: awstypes.ChangeActionInsert,
			Predicate: &awstypes.Predicate{
				Negated: aws.Bool(predicate["negated"].(bool)),
				Type:    awstypes.PredicateType(predicate[names.AttrType].(string)),
				DataId:  aws.String(predicate["data_id"].(string)),
			},
		})
	}
	return updates
}

func DiffRuleGroupActivatedRules(oldRules, newRules []interface{}) []awstypes.RuleGroupUpdate {
	updates := make([]awstypes.RuleGroupUpdate, 0)

	for _, op := range oldRules {
		rule := op.(map[string]interface{})

		if idx, contains := sliceContainsMap(newRules, rule); contains {
			newRules = append(newRules[:idx], newRules[idx+1:]...)
			continue
		}

		updates = append(updates, awstypes.RuleGroupUpdate{
			Action:        awstypes.ChangeActionDelete,
			ActivatedRule: ExpandActivatedRule(rule),
		})
	}

	for _, np := range newRules {
		rule := np.(map[string]interface{})

		updates = append(updates, awstypes.RuleGroupUpdate{
			Action:        awstypes.ChangeActionInsert,
			ActivatedRule: ExpandActivatedRule(rule),
		})
	}
	return updates
}

func FlattenActivatedRules(activatedRules []awstypes.ActivatedRule) []interface{} {
	out := make([]interface{}, len(activatedRules))
	for i, ar := range activatedRules {
		rule := map[string]interface{}{
			"priority":     aws.ToInt32(ar.Priority),
			"rule_id":      aws.ToString(ar.RuleId),
			names.AttrType: ar.Type,
		}
		if ar.Action != nil {
			rule["action"] = []interface{}{
				map[string]interface{}{
					names.AttrType: string(ar.Action.Type),
				},
			}
		}
		out[i] = rule
	}
	return out
}

func ExpandActivatedRule(rule map[string]interface{}) *awstypes.ActivatedRule {
	r := &awstypes.ActivatedRule{
		Priority: aws.Int32(int32(rule["priority"].(int))),
		RuleId:   aws.String(rule["rule_id"].(string)),
		Type:     awstypes.WafRuleType(rule[names.AttrType].(string)),
	}

	if a, ok := rule["action"].([]interface{}); ok && len(a) > 0 {
		m := a[0].(map[string]interface{})
		r.Action = &awstypes.WafAction{
			Type: awstypes.WafActionType(m[names.AttrType].(string)),
		}
	}
	return r
}
