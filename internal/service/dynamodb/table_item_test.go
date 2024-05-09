// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynamodb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdynamodb "github.com/hashicorp/terraform-provider-aws/internal/service/dynamodb"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccDynamoDBTableItem_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var conf dynamodb.GetItemOutput

	tableName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	hashKey := "hashKey"
	itemContent := `{
	"hashKey": {"S": "something"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"},
	"four": {"N": "44444"}
}`

	checkFn := func(s []*terraform.InstanceState) error {
		return testAccDynamoDbTableItemCompareItemAttribute(itemContent, s)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemConfig_basic(tableName, hashKey, itemContent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test", &conf),
					testAccCheckTableItemCount(ctx, tableName, 1),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test", "item", itemContent),
				),
			},
			{
				ResourceName:            "aws_dynamodb_table_item.test",
				ImportStateId:           fmt.Sprintf("%s|%s", tableName, "something"),
				ImportStateCheck:        checkFn,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"item"},
			},
			{
				ResourceName:            "aws_dynamodb_table_item.test",
				ImportStateId:           fmt.Sprintf("[\"%s\", \"%s\"]", tableName, "something"),
				ImportStateCheck:        checkFn,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"item"},
			},
		},
	})
}

func TestAccDynamoDBTableItem_rangeKey(t *testing.T) {
	ctx := acctest.Context(t)
	var conf dynamodb.GetItemOutput

	tableName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	hashKey := "hashKey"
	rangeKey := "rangeKey"
	itemContent := `{
	"hashKey": {"S": "something"},
	"rangeKey": {"S": "something-else"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"},
	"four": {"N": "44444"}
}`

	checkFn := func(s []*terraform.InstanceState) error {
		return testAccDynamoDbTableItemCompareItemAttribute(itemContent, s)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemConfig_rangeKey(tableName, hashKey, rangeKey, itemContent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test", &conf),
					testAccCheckTableItemCount(ctx, tableName, 1),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "range_key", rangeKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test", "item", itemContent),
				),
			},
			{
				ResourceName:            "aws_dynamodb_table_item.test",
				ImportStateId:           fmt.Sprintf("%s|%s|%s", tableName, "something", "something-else"),
				ImportStateCheck:        checkFn,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"item"},
			},
			{
				ResourceName:            "aws_dynamodb_table_item.test",
				ImportStateId:           fmt.Sprintf("[\"%s\", \"%s\", \"%s\"]", tableName, "something", "something-else"),
				ImportStateCheck:        checkFn,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"item"},
			},
		},
	})
}

func TestAccDynamoDBTableItem_withMultipleItems(t *testing.T) {
	ctx := acctest.Context(t)
	var conf1 dynamodb.GetItemOutput
	var conf2 dynamodb.GetItemOutput

	tableName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	hashKey := "hashKey"
	rangeKey := "rangeKey"
	firstItem := `{
	"hashKey": {"S": "something"},
	"rangeKey": {"S": "first"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"}
}`
	secondItem := `{
	"hashKey": {"S": "something"},
	"rangeKey": {"S": "second"},
	"one": {"S": "one"},
	"two": {"S": "two"},
	"three": {"S": "three"},
	"four": {"S": "four"}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemConfig_multiple(tableName, hashKey, rangeKey, firstItem, secondItem),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test1", &conf1),
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test2", &conf2),
					testAccCheckTableItemCount(ctx, tableName, 2),

					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test1", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test1", "range_key", rangeKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test1", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test1", "item", firstItem),

					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test2", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test2", "range_key", rangeKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test2", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test2", "item", secondItem),
				),
			},
		},
	})
}

func TestAccDynamoDBTableItem_withDuplicateItemsSameRangeKey(t *testing.T) {
	ctx := acctest.Context(t)
	tableName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	hashKey := "hashKey"
	rangeKey := "rangeKey"
	firstItem := `{
	"hashKey": {"S": "something"},
	"rangeKey": {"S": "first"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config:      testAccTableItemConfig_multiple(tableName, hashKey, rangeKey, firstItem, firstItem),
				ExpectError: regexache.MustCompile(`ConditionalCheckFailedException: The conditional request failed`),
			},
		},
	})
}

func TestAccDynamoDBTableItem_withDuplicateItemsDifferentRangeKey(t *testing.T) {
	ctx := acctest.Context(t)
	var conf1 dynamodb.GetItemOutput
	var conf2 dynamodb.GetItemOutput

	tableName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	hashKey := "hashKey"
	rangeKey := "rangeKey"
	firstItem := `{
	"hashKey": {"S": "something"},
	"rangeKey": {"S": "first"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"}
}`
	secondItem := `{
	"hashKey": {"S": "something"},
	"rangeKey": {"S": "second"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemConfig_multiple(tableName, hashKey, rangeKey, firstItem, secondItem),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test1", &conf1),
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test2", &conf2),
					testAccCheckTableItemCount(ctx, tableName, 2),

					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test1", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test1", "range_key", rangeKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test1", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test1", "item", firstItem),

					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test2", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test2", "range_key", rangeKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test2", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test2", "item", secondItem),
				),
			},
		},
	})
}

func TestAccDynamoDBTableItem_wonkyItems(t *testing.T) {
	ctx := acctest.Context(t)
	var conf1 dynamodb.GetItemOutput

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	hashKey := "hash.Key"
	rangeKey := "range-Key"
	item := `{
	"hash.Key": {"S": "something"},
	"range-Key": {"S": "first"},
	"one1": {"N": "11111"},
	"two2": {"N": "22222"},
	"three3": {"N": "33333"}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemConfig_wonky(rName, hashKey, rangeKey, item),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test1", &conf1),
					testAccCheckTableItemCount(ctx, rName, 1),

					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test1", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test1", "range_key", rangeKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test1", "table_name", rName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test1", "item", item),
				),
			},
		},
	})
}

func TestAccDynamoDBTableItem_update(t *testing.T) {
	ctx := acctest.Context(t)
	var conf dynamodb.GetItemOutput

	tableName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	hashKey := "hashKey"

	itemBefore := `{
	"hashKey": {"S": "before"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"},
	"four": {"N": "44444"}
}`
	itemAfter := `{
	"hashKey": {"S": "before"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"new": {"S": "shiny new one"}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemConfig_basic(tableName, hashKey, itemBefore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test", &conf),
					testAccCheckTableItemCount(ctx, tableName, 1),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test", "item", itemBefore),
				),
			},
			{
				Config: testAccTableItemConfig_basic(tableName, hashKey, itemAfter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test", &conf),
					testAccCheckTableItemCount(ctx, tableName, 1),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test", "item", itemAfter),
				),
			},
		},
	})
}

func TestAccDynamoDBTableItem_updateWithRangeKey(t *testing.T) {
	ctx := acctest.Context(t)
	var conf dynamodb.GetItemOutput

	tableName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	hashKey := "hashKey"
	rangeKey := "rangeKey"

	itemBefore := `{
	"hashKey": {"S": "before"},
	"rangeKey": {"S": "rangeBefore"},
	"value": {"S": "valueBefore"}
}`
	itemAfter := `{
	"hashKey": {"S": "before"},
	"rangeKey": {"S": "rangeAfter"},
	"value": {"S": "valueAfter"}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemConfig_rangeKey(tableName, hashKey, rangeKey, itemBefore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test", &conf),
					testAccCheckTableItemCount(ctx, tableName, 1),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "range_key", rangeKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test", "item", itemBefore),
				),
			},
			{
				Config: testAccTableItemConfig_rangeKey(tableName, hashKey, rangeKey, itemAfter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test", &conf),
					testAccCheckTableItemCount(ctx, tableName, 1),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "range_key", rangeKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "table_name", tableName),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test", "item", itemAfter),
				),
			},
		},
	})
}

func TestAccDynamoDBTableItem_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var conf dynamodb.GetItemOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_dynamodb_table_item.test"

	hashKey := "hashKey"
	itemContent := `{
	"hashKey": {"S": "something"},
	"one": {"N": "11111"},
	"two": {"N": "22222"},
	"three": {"N": "33333"},
	"four": {"N": "44444"}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemConfig_basic(rName, hashKey, itemContent),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTableItemExists(ctx, resourceName, &conf),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdynamodb.ResourceTableItem(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDynamoDBTableItem_importWithHashKey(t *testing.T) {
	ctx := acctest.Context(t)
	var conf dynamodb.GetItemOutput

	tableName := fmt.Sprintf("tf-acc-test-%v", acctest.RandInt(t))
	hashKey := "hashKey"
	hashKeyValue := "importHashKey"

	item := fmt.Sprintf(`{
	"%s": {"S": "%s"},
	"otherAttrS": {"S": "otherStringValue"},
	"otherAttrN": {"N": "123"},
	"otherAttrSS": {"SS": ["a", "b", "c"]},
	"otherAttrNS": {"NS": ["0", "1.1", "-2.22"]},
	"otherAttrBool": {"BOOL": false},
	"otherAttrNull": {"NULL": true}
}`, hashKey, hashKeyValue)

	checkFn := func(s []*terraform.InstanceState) error {
		return testAccDynamoDbTableItemCompareItemAttribute(item, s)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, dynamodb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemImport_hashKey(tableName, hashKey, item),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test", &conf),
					testAccCheckTableItemCount(ctx, tableName, 1),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "table_name", tableName),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "item", item+"\n"),
				),
			},
			{
				ResourceName:            "aws_dynamodb_table_item.test",
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("%s|%s", tableName, hashKeyValue),
				ImportStateCheck:        checkFn,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"item"},
			},
			{
				ResourceName:            "aws_dynamodb_table_item.test",
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("[ \"%s\", \"%s\" ]", tableName, hashKeyValue),
				ImportStateCheck:        checkFn,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"item"},
			},
		},
	})
}

func TestAccDynamoDBTableItem_importWithRangeKey(t *testing.T) {
	ctx := acctest.Context(t)
	var conf dynamodb.GetItemOutput

	tableName := fmt.Sprintf("tf-acc-test-%v", acctest.RandInt(t))
	hashKey := "hashKey"
	hashKeyValue := "dGhpcyB0ZXh0IGlzIGJhc2U2NC1lbmNvZGVk"
	rangeKey := "rangeKey"
	rangeKeyValue := "100"

	item := fmt.Sprintf(`{
	"%s": {"B": "%s"},
	"%s": {"N": "%s"},
	"otherAttrS": {"S": "otherStringValue"},
	"otherAttrN": {"N": "123"},
	"otherAttrSS": {"SS": ["a", "b", "c"]},
	"otherAttrNS": {"NS": ["0", "1.1", "-2.22"]},
	"otherAttrBool": {"BOOL": false},
	"otherAttrNull": {"NULL": true}
}`, hashKey, hashKeyValue, rangeKey, rangeKeyValue)

	checkFn := func(s []*terraform.InstanceState) error {
		return testAccDynamoDbTableItemCompareItemAttribute(item, s)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, dynamodb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemImport_rangeKey(tableName, hashKey, rangeKey, item),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test", &conf),
					testAccCheckTableItemCount(ctx, tableName, 1),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "hash_key", hashKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "range_key", rangeKey),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "table_name", tableName),
					resource.TestCheckResourceAttr("aws_dynamodb_table_item.test", "item", item+"\n"),
				),
			},
			{
				ResourceName:            "aws_dynamodb_table_item.test",
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("%s|%s|%s", tableName, hashKeyValue, rangeKeyValue),
				ImportStateCheck:        checkFn,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"item"},
			},
			{
				ResourceName:            "aws_dynamodb_table_item.test",
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("[ \"%s\", \"%s\", \"%s\" ]", tableName, hashKeyValue, rangeKeyValue),
				ImportStateCheck:        checkFn,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"item"},
			},
		},
	})
}

func testAccDynamoDbTableItemCompareItemAttribute(item string, s []*terraform.InstanceState) error {
	if len(s) != 1 {
		return fmt.Errorf("expected 1 state: %#v", s)
	}
	var a, b map[string]interface{}

	err := json.Unmarshal([]byte(item), &a)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(s[0].Attributes["item"]), &b)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(a, b) {
		return fmt.Errorf("item attributes not equal:\n\texpected: %#v\n\tactual: %#v", a, b)
	}

	return nil
}

func TestAccDynamoDBTableItem_mapOutOfBandUpdate(t *testing.T) {
	ctx := acctest.Context(t)
	var conf dynamodb.GetItemOutput

	tableName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	hashKey := names.AttrKey
	tmpl := `{
	"key": {"S": "something"},
	"value": {
		"M": {
			"valid_after": {
				"N": %[1]q
			}
		}
	},
	"other": {
		"N": %[1]q
	}
}`

	oldValue := "300"
	newValue := "400"

	oldItem := fmt.Sprintf(tmpl, oldValue)
	newItem := fmt.Sprintf(tmpl, newValue)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableItemDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableItemConfig_map(tableName, names.AttrKey, oldItem),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTableItemExists(ctx, "aws_dynamodb_table_item.test", &conf),
					testAccCheckTableItemCount(ctx, tableName, 1),
					acctest.CheckResourceAttrEquivalentJSON("aws_dynamodb_table_item.test", "item", oldItem),
					acctest.CheckResourceAttrJMES("aws_dynamodb_table_item.test", "item", "value.M.valid_after.N", oldValue),
				),
			},
			{
				PreConfig: func() {
					conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBConn(ctx)

					attributes, err := tfdynamodb.ExpandTableItemAttributes(newItem)
					if err != nil {
						t.Fatalf("making out-of-band change: %s", err)
					}

					updates := map[string]*dynamodb.AttributeValueUpdate{}
					for key, value := range attributes {
						if key == hashKey {
							continue
						}
						updates[key] = &dynamodb.AttributeValueUpdate{
							Action: aws.String(dynamodb.AttributeActionPut),
							Value:  value,
						}
					}

					newQueryKey := tfdynamodb.BuildTableItemQueryKey(attributes, hashKey, "")
					_, err = conn.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
						AttributeUpdates: updates,
						TableName:        aws.String(tableName),
						Key:              newQueryKey,
					})
					if err != nil {
						t.Fatalf("making out-of-band change: %s", err)
					}
				},
				Config:   testAccTableItemConfig_map(tableName, names.AttrKey, newItem),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckTableItemDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_dynamodb_table_item" {
				continue
			}

			attrs := rs.Primary.Attributes
			attributes, err := tfdynamodb.ExpandTableItemAttributes(attrs["item"])
			if err != nil {
				return err
			}

			key := tfdynamodb.BuildTableItemQueryKey(attributes, attrs["hash_key"], attrs["range_key"])

			_, err = tfdynamodb.FindTableItem(ctx, conn, attrs["table_name"], key)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DynamoDB table item %s still exists.", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckTableItemExists(ctx context.Context, n string, item *dynamodb.GetItemOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DynamoDB table item ID specified!")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBConn(ctx)

		attrs := rs.Primary.Attributes
		attributes, err := tfdynamodb.ExpandTableItemAttributes(attrs["item"])
		if err != nil {
			return err
		}

		key := tfdynamodb.BuildTableItemQueryKey(attributes, attrs["hash_key"], attrs["range_key"])

		result, err := tfdynamodb.FindTableItem(ctx, conn, attrs["table_name"], key)

		if err != nil {
			return err
		}

		*item = *result

		return nil
	}
}

func testAccCheckTableItemCount(ctx context.Context, tableName string, count int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBConn(ctx)
		out, err := conn.ScanWithContext(ctx, &dynamodb.ScanInput{
			ConsistentRead: aws.Bool(true),
			TableName:      aws.String(tableName),
			Select:         aws.String(dynamodb.SelectCount),
		})
		if err != nil {
			return err
		}
		expectedCount := count
		if *out.Count != expectedCount {
			return fmt.Errorf("Expected %d items, got %d", expectedCount, *out.Count)
		}
		return nil
	}
}

func testAccTableItemImport_hashKey(tableName, hashKey, item string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = "%s"
  read_capacity  = 5
  write_capacity = 5
	hash_key       = "%s"
  attribute {
    name = "%s"
    type = "S"
  }
}
resource "aws_dynamodb_table_item" "test" {
  table_name = "${aws_dynamodb_table.test.name}"
	hash_key   = "${aws_dynamodb_table.test.hash_key}"
  item = <<ITEM
%s
ITEM
}
`, tableName, hashKey, hashKey, item)
}

func testAccTableItemImport_rangeKey(tableName, hashKey, rangeKey, item string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = "%s"
  read_capacity  = 5
  write_capacity = 5
	hash_key       = "%s"
	range_key      = "%s"
  attribute {
    name = "%s"
    type = "B"
  }
  attribute {
    name = "%s"
    type = "N"
  }
}
resource "aws_dynamodb_table_item" "test" {
  table_name = "${aws_dynamodb_table.test.name}"
	hash_key   = "${aws_dynamodb_table.test.hash_key}"
	range_key  = "${aws_dynamodb_table.test.range_key}"
  item = <<ITEM
%s
ITEM
}
`, tableName, hashKey, rangeKey, hashKey, rangeKey, item)
}

func testAccTableItemConfig_basic(tableName, hashKey, item string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = "%s"
  read_capacity  = 10
  write_capacity = 10
  hash_key       = "%s"

  attribute {
    name = "%s"
    type = "S"
  }
}

resource "aws_dynamodb_table_item" "test" {
  table_name = aws_dynamodb_table.test.name
  hash_key   = aws_dynamodb_table.test.hash_key

  item = <<ITEM
%s
ITEM
}
`, tableName, hashKey, hashKey, item)
}

func testAccTableItemConfig_rangeKey(tableName, hashKey, rangeKey, item string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = "%s"
  read_capacity  = 10
  write_capacity = 10
  hash_key       = "%s"
  range_key      = "%s"

  attribute {
    name = "%s"
    type = "S"
  }

  attribute {
    name = "%s"
    type = "S"
  }
}

resource "aws_dynamodb_table_item" "test" {
  table_name = aws_dynamodb_table.test.name
  hash_key   = aws_dynamodb_table.test.hash_key
  range_key  = aws_dynamodb_table.test.range_key

  item = <<ITEM
%s
ITEM
}
`, tableName, hashKey, rangeKey, hashKey, rangeKey, item)
}

func testAccTableItemConfig_multiple(tableName, hashKey, rangeKey, firstItem, secondItem string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = "%s"
  read_capacity  = 10
  write_capacity = 10
  hash_key       = "%s"
  range_key      = "%s"

  attribute {
    name = "%s"
    type = "S"
  }

  attribute {
    name = "%s"
    type = "S"
  }
}

resource "aws_dynamodb_table_item" "test1" {
  table_name = aws_dynamodb_table.test.name
  hash_key   = aws_dynamodb_table.test.hash_key
  range_key  = aws_dynamodb_table.test.range_key

  item = <<ITEM
%s
ITEM
}

resource "aws_dynamodb_table_item" "test2" {
  table_name = aws_dynamodb_table.test.name
  hash_key   = aws_dynamodb_table.test.hash_key
  range_key  = aws_dynamodb_table.test.range_key

  item = <<ITEM
%s
ITEM
}
`, tableName, hashKey, rangeKey, hashKey, rangeKey, firstItem, secondItem)
}

func testAccTableItemConfig_wonky(tableName, hashKey, rangeKey, item string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = %[1]q
  read_capacity  = 10
  write_capacity = 10
  hash_key       = %[2]q
  range_key      = %[3]q

  attribute {
    name = %[2]q
    type = "S"
  }

  attribute {
    name = %[3]q
    type = "S"
  }
}

resource "aws_dynamodb_table_item" "test1" {
  table_name = aws_dynamodb_table.test.name
  hash_key   = aws_dynamodb_table.test.hash_key
  range_key  = aws_dynamodb_table.test.range_key

  item = <<ITEM
%[4]s
ITEM
}
`, tableName, hashKey, rangeKey, item)
}

func testAccADynamoDbItemCompareItemAttribute(item string, s []*terraform.InstanceState) error {
	if len(s) != 1 {
		return fmt.Errorf("expected 1 state: %#v", s)
	}
	var a, b map[string]interface{}

	err := json.Unmarshal([]byte(item), &a)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(s[0].Attributes["item"]), &b)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(a, b) {
		return fmt.Errorf("item attributes not equal:\n\texpected: %#v\n\tactual: %#v", a, b)
	}

	return nil
}
func testAccTableItemConfig_map(tableName, hashKey, content string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = %[1]q
  read_capacity  = 10
  write_capacity = 10
  hash_key       = %[2]q

  attribute {
    name = %[2]q
    type = "S"
  }
}

resource "aws_dynamodb_table_item" "test" {
  table_name = aws_dynamodb_table.test.name
  hash_key   = aws_dynamodb_table.test.hash_key

  item = <<ITEM
%[3]s
ITEM
}
`, tableName, hashKey, content)
}
