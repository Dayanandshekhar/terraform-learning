// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ecs_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccECSClusterDataSource_ecsCluster(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_ecs_cluster.test"
	resourceName := "aws_ecs_cluster.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ECSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, names.AttrARN, resourceName, names.AttrARN),
					resource.TestCheckResourceAttr(dataSourceName, "pending_tasks_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "registered_container_instances_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "running_tasks_count", "0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "service_connect_defaults.#", resourceName, "service_connect_defaults.#"),
					resource.TestCheckResourceAttr(dataSourceName, names.AttrStatus, "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccECSClusterDataSource_ecsClusterContainerInsights(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_ecs_cluster.test"
	resourceName := "aws_ecs_cluster.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ECSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig_containerInsights(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, names.AttrARN, resourceName, names.AttrARN),
					resource.TestCheckResourceAttr(dataSourceName, "pending_tasks_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "registered_container_instances_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "running_tasks_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, names.AttrStatus, "ACTIVE"),
					resource.TestCheckResourceAttrPair(dataSourceName, "setting.#", resourceName, "setting.#"),
				),
			},
		},
	})
}

func TestAccECSClusterDataSource_tags(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_ecs_cluster.test"
	resourceName := "aws_ecs_cluster.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ECSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig_tags(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, names.AttrARN, resourceName, names.AttrARN),
					resource.TestCheckResourceAttr(dataSourceName, "pending_tasks_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "registered_container_instances_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "running_tasks_count", "0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "service_connect_defaults.#", resourceName, "service_connect_defaults.#"),
					resource.TestCheckResourceAttr(dataSourceName, names.AttrStatus, "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
		},
	})
}

func TestAccECSClusterDataSource_ecsClusterCapacityProvider(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_ecs_cluster.test"
	resourceName := "aws_ecs_cluster.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ECSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig_capacityProvider(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttr(dataSourceName, "pending_tasks_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "registered_container_instances_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "running_tasks_count", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(dataSourceName, "capacity_providers.#", resourceName, "capacity_providers.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "default_capacity_provider_strategy.#", resourceName, "default_capacity_provider_strategy.#"),
				),
			},
		},
	})
}

func testAccClusterDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_ecs_cluster" "test" {
  name = %[1]q
}

data "aws_ecs_cluster" "test" {
  cluster_name = aws_ecs_cluster.test.name
}
`, rName)
}

func testAccClusterDataSourceConfig_containerInsights(rName string) string {
	return fmt.Sprintf(`
resource "aws_ecs_cluster" "test" {
  name = %[1]q

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

data "aws_ecs_cluster" "test" {
  cluster_name = aws_ecs_cluster.test.name
}
`, rName)
}

func testAccClusterDataSourceConfig_tags(rName, tagKey, tagValue string) string {
	return fmt.Sprintf(`
resource "aws_ecs_cluster" "test" {
  name = %[1]q
  tags = {
    %[2]q = %[3]q
  }
}
data "aws_ecs_cluster" "test" {
  cluster_name = aws_ecs_cluster.test.name
}
`, rName, tagKey, tagValue)
}

func testAccClusterDataSourceConfig_capacityProvider(rName string) string {
	return fmt.Sprintf(`
resource "aws_ecs_cluster" "test" {
  name               = %[1]q
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    base              = 0
    weight            = 1
    capacity_provider = "FARGATE_SPOT"
  }
}

data "aws_ecs_cluster" "test" {
  cluster_name = aws_ecs_cluster.test.name
}
`, rName)
}
