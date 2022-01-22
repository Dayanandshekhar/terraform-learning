package imagebuilder

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceDistributionConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDistributionConfigurations,
		Schema: map[string]*schema.Schema{
			"arns": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": dataSourceFiltersSchema(),
			"names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDistributionConfigurations(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn

	input := &imagebuilder.ListDistributionConfigurationsInput{}

	if v, ok := d.GetOk("filter"); ok {
		input.Filters = buildFiltersDataSource(v.(*schema.Set))
	}

	var results []*imagebuilder.DistributionConfigurationSummary

	err := conn.ListDistributionConfigurationsPages(input, func(page *imagebuilder.ListDistributionConfigurationsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, distributionConfigurationSummary := range page.DistributionConfigurationSummaryList {
			if distributionConfigurationSummary == nil {
				continue
			}

			results = append(results, distributionConfigurationSummary)
		}

		return !lastPage
	})

	if err != nil {
		return fmt.Errorf("error reading Image Builder Distribution Configurations: %w", err)
	}

	var arns, names []string

	for _, r := range results {
		arns = append(arns, aws.StringValue(r.Arn))
		names = append(names, aws.StringValue(r.Name))
	}

	d.SetId(meta.(*conns.AWSClient).Region)
	d.Set("arns", arns)
	d.Set("names", names)

	return nil
}
