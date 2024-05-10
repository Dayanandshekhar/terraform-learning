// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mq_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/hashicorp/terraform-plugin-testing/config"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfmq "github.com/hashicorp/terraform-provider-aws/internal/service/mq"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccMQConfiguration_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mq_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.MQEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MQServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				//Config: testAccConfigurationConfig_basic(rName),
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"description":             config.StringVariable("TfAccTest MQ Configuration"),
					"random_name":             config.StringVariable(rName),
					"engine_type":             config.StringVariable("ActiveMQ"),
					"engine_version":          config.StringVariable("5.17.6"),
					"authentication_strategy": config.StringVariable("simple"),
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "mq", regexache.MustCompile(`configuration:+.`)),
					resource.TestCheckResourceAttr(resourceName, "authentication_strategy", "simple"),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, "TfAccTest MQ Configuration"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "ActiveMQ"),
					resource.TestCheckResourceAttr(resourceName, names.AttrEngineVersion, "5.17.6"),
					resource.TestCheckResourceAttr(resourceName, "latest_revision", "2"),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				//Config: testAccConfigurationConfig_descriptionUpdated(rName),
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"description":    config.StringVariable("TfAccTest MQ Configuration Updated"),
					"random_name":    config.StringVariable(rName),
					"engine_type":    config.StringVariable("ActiveMQ"),
					"engine_version": config.StringVariable("5.17.6"),
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "mq", regexache.MustCompile(`configuration:+.`)),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, "TfAccTest MQ Configuration Updated"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "ActiveMQ"),
					resource.TestCheckResourceAttr(resourceName, names.AttrEngineVersion, "5.17.6"),
					resource.TestCheckResourceAttr(resourceName, "latest_revision", "3"),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
				),
			},
		},
	})
}

func TestAccMQConfiguration_activeMQData(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mq_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.MQEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MQServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				//Config: testAccConfigurationConfig_activeData(rName),
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"description":    config.StringVariable("TfAccTest MQ Configuration"),
					"random_name":    config.StringVariable(rName),
					"engine_type":    config.StringVariable("ActiveMQ"),
					"engine_version": config.StringVariable("5.17.6"),
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "mq", regexache.MustCompile(`configuration:+.`)),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, "TfAccTest MQ Configuration"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "ActiveMQ"),
					resource.TestCheckResourceAttr(resourceName, names.AttrEngineVersion, "5.17.6"),
					resource.TestCheckResourceAttr(resourceName, "latest_revision", "2"),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMQConfiguration_activeMQLDAPData(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mq_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.MQEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MQServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				//Config: testAccConfigurationConfig_activeLdapData(rName),
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"description":             config.StringVariable("TfAccTest MQ Configuration"),
					"random_name":             config.StringVariable(rName),
					"engine_type":             config.StringVariable("ActiveMQ"),
					"engine_version":          config.StringVariable("5.17.6"),
					"authentication_strategy": config.StringVariable("ldap"),
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "mq", regexache.MustCompile(`configuration:+.`)),
					resource.TestCheckResourceAttr(resourceName, "authentication_strategy", "ldap"),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, "TfAccTest MQ Configuration"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "ActiveMQ"),
					resource.TestCheckResourceAttr(resourceName, names.AttrEngineVersion, "5.17.6"),
					resource.TestCheckResourceAttr(resourceName, "latest_revision", "2"),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMQConfiguration_rabbitMQData(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mq_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.MQEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MQServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				//Config: testAccConfigurationConfig_rabbitData(rName),
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"description":    config.StringVariable("TfAccTest MQ Configuration"),
					"random_name":    config.StringVariable(rName),
					"engine_type":    config.StringVariable("RabbitMQ"),
					"engine_version": config.StringVariable("3.11.16"),
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "mq", regexache.MustCompile(`configuration:+.`)),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, "TfAccTest MQ Configuration"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "RabbitMQ"),
					resource.TestCheckResourceAttr(resourceName, names.AttrEngineVersion, "3.11.16"),
					resource.TestCheckResourceAttr(resourceName, "latest_revision", "2"),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "data", "consumer_timeout = 60000\n"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMQConfiguration_tags(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_mq_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.MQEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MQServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				//Config: testAccConfigurationConfig_tags1(rName, "key1", "value1"),
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"description":             config.StringVariable("TfAccTest MQ Configuration"),
					"random_name":             config.StringVariable(rName),
					"engine_type":             config.StringVariable("ActiveMQ"),
					"engine_version":          config.StringVariable("5.17.6"),
					"authentication_strategy": config.StringVariable("simple"),
					"key1_value":              config.StringVariable("value1"),
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				//Config: testAccConfigurationConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"description":             config.StringVariable("TfAccTest MQ Configuration"),
					"random_name":             config.StringVariable(rName),
					"engine_type":             config.StringVariable("ActiveMQ"),
					"engine_version":          config.StringVariable("5.17.6"),
					"authentication_strategy": config.StringVariable("simple"),
					"key1_value":              config.StringVariable("value1updated"),
					"key2_value":              config.StringVariable("value2"),
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				//Config: testAccConfigurationConfig_tags1(rName, "key2", "value2"),
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"description":             config.StringVariable("TfAccTest MQ Configuration"),
					"random_name":             config.StringVariable(rName),
					"engine_type":             config.StringVariable("ActiveMQ"),
					"engine_version":          config.StringVariable("5.17.6"),
					"authentication_strategy": config.StringVariable("simple"),
					"key2_value":              config.StringVariable("value2"),
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckConfigurationExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).MQClient(ctx)

		_, err := tfmq.FindConfigurationByID(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccConfigurationConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_mq_configuration" "test" {
  description             = "TfAccTest MQ Configuration"
  name                    = %[1]q
  engine_type             = "ActiveMQ"
  engine_version          = "5.17.6"
  authentication_strategy = "simple"

  data = <<DATA
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<broker xmlns="http://activemq.apache.org/schema/core">
</broker>
DATA
}
`, rName)
}

func testAccConfigurationConfig_descriptionUpdated(rName string) string {
	return fmt.Sprintf(`
resource "aws_mq_configuration" "test" {
  description    = "TfAccTest MQ Configuration Updated"
  name           = %[1]q
  engine_type    = "ActiveMQ"
  engine_version = "5.17.6"

  data = <<DATA
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<broker xmlns="http://activemq.apache.org/schema/core">
</broker>
DATA
}
`, rName)
}

func testAccConfigurationConfig_activeData(rName string) string {
	return fmt.Sprintf(`
resource "aws_mq_configuration" "test" {
  description    = "TfAccTest MQ Configuration"
  name           = %[1]q
  engine_type    = "ActiveMQ"
  engine_version = "5.17.6"

  data = <<DATA
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<broker xmlns="http://activemq.apache.org/schema/core">
  <plugins>
    <authorizationPlugin>
      <map>
        <authorizationMap>
          <authorizationEntries>
            <authorizationEntry admin="guests,users" queue="GUEST.&gt;" read="guests" write="guests,users"/>
            <authorizationEntry admin="guests,users" read="guests,users" topic="ActiveMQ.Advisory.&gt;" write="guests,users"/>
          </authorizationEntries>
          <tempDestinationAuthorizationEntry>
            <tempDestinationAuthorizationEntry admin="tempDestinationAdmins" read="tempDestinationAdmins" write="tempDestinationAdmins"/>
          </tempDestinationAuthorizationEntry>
        </authorizationMap>
      </map>
    </authorizationPlugin>
    <forcePersistencyModeBrokerPlugin persistenceFlag="true"/>
    <statisticsBrokerPlugin/>
    <timeStampingBrokerPlugin ttlCeiling="86400000" zeroExpirationOverride="86400000"/>
  </plugins>
</broker>
DATA
}
`, rName)
}

func testAccConfigurationConfig_activeLdapData(rName string) string {
	return fmt.Sprintf(`
resource "aws_mq_configuration" "test" {
  description             = "TfAccTest MQ Configuration"
  name                    = %[1]q
  engine_type             = "ActiveMQ"
  engine_version          = "5.17.6"
  authentication_strategy = "ldap"

  data = <<DATA
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<broker xmlns="http://activemq.apache.org/schema/core">
  <plugins>
    <authorizationPlugin>
      <map>
        <cachedLDAPAuthorizationMap legacyGroupMapping="false" queueSearchBase="ou=Queue,ou=Destination,ou=ActiveMQ,dc=example,dc=org" refreshInterval="0" tempSearchBase="ou=Temp,ou=Destination,ou=ActiveMQ,dc=example,dc=org" topicSearchBase="ou=Topic,ou=Destination,ou=ActiveMQ,dc=example,dc=org"/>
      </map>
    </authorizationPlugin>
    <forcePersistencyModeBrokerPlugin persistenceFlag="true"/>
    <statisticsBrokerPlugin/>
    <timeStampingBrokerPlugin ttlCeiling="86400000" zeroExpirationOverride="86400000"/>
  </plugins>
</broker>
DATA
}
`, rName)
}

func testAccConfigurationConfig_rabbitData(rName string) string {
	return fmt.Sprintf(`
resource "aws_mq_configuration" "test" {
  description    = "TfAccTest MQ Configuration"
  name           = %[1]q
  engine_type    = "RabbitMQ"
  engine_version = "3.11.16"

  data = <<DATA
consumer_timeout = 60000
DATA
}
`, rName)
}

func testAccConfigurationConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_mq_configuration" "test" {
  description             = "TfAccTest MQ Configuration"
  name                    = %[1]q
  engine_type             = "ActiveMQ"
  engine_version          = "5.17.6"
  authentication_strategy = "simple"

  tags = {
    %[2]q = %[3]q
  }

  data = <<DATA
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<broker xmlns="http://activemq.apache.org/schema/core">
</broker>
DATA
}
`, rName, tagKey1, tagValue1)
}

func testAccConfigurationConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_mq_configuration" "test" {
  description             = "TfAccTest MQ Configuration"
  name                    = %[1]q
  engine_type             = "ActiveMQ"
  engine_version          = "5.17.6"
  authentication_strategy = "simple"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }

  data = <<DATA
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<broker xmlns="http://activemq.apache.org/schema/core">
</broker>
DATA
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
