package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
)

func resourceAwsAcmCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsAcmCertificateCreate,
		Read:   resourceAwsAcmCertificateRead,
		Delete: resourceAwsAcmCertificateDelete,

		Schema: map[string]*schema.Schema{
			"domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subject_alternative_names": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"validation_method": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"domain_validation_options": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_record_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_record_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_record_value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAwsAcmCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	acmconn := meta.(*AWSClient).acmconn
	params := &acm.RequestCertificateInput{
		DomainName:       aws.String(d.Get("domain_name").(string)),
		ValidationMethod: aws.String("DNS"),
	}

	// TODO: check that validation method is DNS, nothing else supported at the moment

	sans, ok := d.GetOk("subject_alternative_names")
	if ok {
		sanStrings := sans.([]interface{})
		params.SubjectAlternativeNames = expandStringList(sanStrings)
	}

	log.Printf("[DEBUG] ACM Certificate Request: %#v", params)
	resp, err := acmconn.RequestCertificate(params)

	if err != nil {
		return fmt.Errorf("Error requesting certificate: %s", err)
	}

	d.SetId(*resp.CertificateArn)
	d.Set("certificate_arn", *resp.CertificateArn)

	return resourceAwsAcmCertificateRead(d, meta)
}

func resourceAwsAcmCertificateRead(d *schema.ResourceData, meta interface{}) error {
	acmconn := meta.(*AWSClient).acmconn

	params := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(d.Id()),
	}

	return resource.Retry(time.Duration(1)*time.Minute, func() *resource.RetryError {
		resp, err := acmconn.DescribeCertificate(params)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error describing certificate: %s", err))
		}

		if err := d.Set("domain_name", resp.Certificate.DomainName); err != nil {
			return resource.NonRetryableError(err)
		}
		if err := d.Set("subject_alternative_names", cleanUpSubjectAlternativeNames(resp.Certificate)); err != nil {
			return resource.NonRetryableError(err)
		}

		domainValidationOptions, err := convertDomainValidationOptions(resp.Certificate.DomainValidationOptions)

		if err != nil {
			return resource.RetryableError(err)
		}

		if err := d.Set("domain_validation_options", domainValidationOptions); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

}

func cleanUpSubjectAlternativeNames(cert *acm.CertificateDetail) []string {
	sans := cert.SubjectAlternativeNames
	vs := make([]string, 0, len(sans)-1)
	for _, v := range sans {
		if *v != *cert.DomainName {
			vs = append(vs, *v)
		}
	}
	return vs

}

func convertDomainValidationOptions(validations []*acm.DomainValidation) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(validations))

	for _, o := range validations {
		validationOption := make(map[string]interface{})
		validationOption["domain_name"] = *o.DomainName
		if o.ResourceRecord != nil {
			validationOption["resource_record_name"] = *o.ResourceRecord.Name
			validationOption["resource_record_type"] = *o.ResourceRecord.Type
			validationOption["resource_record_value"] = *o.ResourceRecord.Value
		} else {
			log.Printf("[DEBUG] No resource record found in validation options, need to retry: %#v", o)
			return nil, fmt.Errorf("No resource record found in DNS DomainValidationOptions: %v", o)
		}

		result = append(result, validationOption)
	}

	return result, nil
}

func resourceAwsAcmCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	acmconn := meta.(*AWSClient).acmconn

	if err := resourceAwsAcmCertificateRead(d, meta); err != nil {
		return err
	}
	if d.Id() == "" {
		// This might happen from the read
		return nil
	}

	params := &acm.DeleteCertificateInput{
		CertificateArn: aws.String(d.Id()),
	}

	_, err := acmconn.DeleteCertificate(params)

	if err != nil {
		return fmt.Errorf("Error deleting certificate: %s", err)
	}

	d.SetId("")
	return nil
}
