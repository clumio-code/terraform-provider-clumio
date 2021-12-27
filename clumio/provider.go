// Copyright 2021. Clumio, Inc.

// This file contains the functions related to provider definition and initialization.

package clumio

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	awsCredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	awsAccountIDRegexpInternalPattern = `(aws|\d{12})`
	awsPartitionRegexpInternalPattern = `aws(-[a-z]+)*`
	awsRegionRegexpInternalPattern    = `[a-z]{2}(-[a-z]+)+-\d`
	awsAccountIDRegexpPattern         = "^" + awsAccountIDRegexpInternalPattern + "$"
	awsPartitionRegexpPattern         = "^" + awsPartitionRegexpInternalPattern + "$"
	awsRegionRegexpPattern            = "^" + awsRegionRegexpInternalPattern + "$"
	awsProfile = "AWS_PROFILE"
	awsSharedCredsFile = "AWS_SHARED_CREDENTIALS_FILE"
	clumioApiToken = "CLUMIO_API_TOKEN"
	clumioApiBaseUrl = "CLUMIO_API_BASE_URL"
)

var awsAccountIDRegexp = regexp.MustCompile(awsAccountIDRegexpPattern)
var awsPartitionRegexp = regexp.MustCompile(awsPartitionRegexpPattern)
var awsRegionRegexp = regexp.MustCompile(awsRegionRegexpPattern)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

// New is the factory method that returns a function which, when called,
// creates a new Provider instance.
func New(isUnitTest bool) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"clumio_callback_resource": clumioCallback(),
				"clumio_aws_connection": clumioAWSConnection(),
				"clumio_post_process_aws_connection": clumioPostProcessAWSConnection(),
			},
		}
		p.ConfigureContextFunc = configure(p, isUnitTest)
		p.Schema = map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AWS Access Key.",
			},
			"assume_role": assumeRoleSchema(),
			"clumio_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Clumio Control Plane AWS Region.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Region.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AWS Secret Key.",
			},
			"session_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AWS Session Token.",
			},
			"profile": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "The profile for API operations. If not set, the default profile\n" +
					"created with `aws configure` will be used.",
			},
			"shared_credentials_file": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "The path to the shared credentials file. If not set\n" +
					"this defaults to ~/.aws/credentials.",
			},
			"clumio_api_token": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "The API token required to invoke Clumio APIs.",
			},
			"clumio_api_base_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://api.clumio.com",
				Description: "The base URL for Clumio APIs.",
			},
		}
		return p
	}
}

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"duration_seconds": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Seconds to restrict the assume role session duration. Defaults to 15 minutes if not set.",
				},
				"external_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Unique identifier that might be required for assuming a role in another account.",
				},
				"role_arn": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Amazon Resource Name of an IAM Role to assume prior to making API calls.",
					ValidateFunc: validateArn,
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Identifier for the assumed role session.",
				},
			},
		},
	}
}

func validateArn(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return ws, errors
	}

	parsedARN, err := arn.Parse(value)

	if err != nil {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: %s", k, value, err))
		return ws, errors
	}

	if parsedARN.Partition == "" {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: missing partition value", k, value))
	} else if !awsPartitionRegexp.MatchString(parsedARN.Partition) {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: invalid partition value (expecting to match regular expression: %s)", k, value, awsPartitionRegexpPattern))
	}

	if parsedARN.Region != "" && !awsRegionRegexp.MatchString(parsedARN.Region) {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: invalid region value (expecting to match regular expression: %s)", k, value, awsRegionRegexpPattern))
	}

	if parsedARN.AccountID != "" && !awsAccountIDRegexp.MatchString(parsedARN.AccountID) {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: invalid account ID value (expecting to match regular expression: %s)", k, value, awsAccountIDRegexpPattern))
	}

	if parsedARN.Resource == "" {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: missing resource value", k, value))
	}

	return ws, errors
}

// apiClient defines the APIs/connections required by the resources.
type apiClient struct {
	snsAPI SNSAPI
	s3API  S3API
	clumioConfig clumioConfig.Config
}

// configure is a factory method to configure the Provider.
func configure(_ *schema.Provider, isUnitTest bool) func(context.Context,
	*schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{},
		diag.Diagnostics) {
		if isUnitTest {
			return &apiClient{
				snsAPI: SNSClient{},
				s3API:  S3Client{},
			}, nil
		}

		apiToken := getStringValue(d, "clumio_api_token")
		baseUrl := getStringValue(d, "clumio_api_base_url")
		if apiToken == ""{
			apiToken = os.Getenv(clumioApiToken)
		}
		if baseUrl == ""{
			baseUrl = os.Getenv(clumioApiBaseUrl)
		}

		accessKey := getStringValue(d, "access_key")
		secretKey := getStringValue(d, "secret_key")
		sessionToken := getStringValue(d, "session_token")
		profile := getStringValue(d, "profile")
		if profile == "" {
			profile = os.Getenv(awsProfile)
		}
		sharedCredsFile := getStringValue(d, "shared_credentials_file")
		if sharedCredsFile == "" {
			sharedCredsFile = os.Getenv(awsSharedCredsFile)
		}
		loadOpts := []func(options *config.LoadOptions) error{
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), kMaxRetries)
			}),
		}
		if profile != "" {
			loadOpts = append(loadOpts, config.WithSharedConfigProfile(profile))
		}
		if sharedCredsFile != "" {
			loadOpts = append(
				loadOpts, config.WithSharedCredentialsFiles([]string{sharedCredsFile}))
		}

		cfg, err := config.LoadDefaultConfig(context.TODO(), loadOpts...)
		if err != nil {
			return nil, diag.Errorf(
				"Error loading default config for AWS Provider: %v", err)
		}

		region := getStringValue(d, "region")
		if region != "" {
			cfg.Region = region
		}

		var assumeRoleOptions *stscreds.AssumeRoleOptions
		var diagErr diag.Diagnostics
		if assumeRoleList, ok := d.Get("assume_role").([]interface{}); ok && len(assumeRoleList) > 0 && assumeRoleList[0] != nil {
			assumeRoleOptions, diagErr = getAssumeRoleOptions(assumeRoleList[0])
			if diagErr != nil {
				return nil, diagErr
			}
			assumeRoleOptionsFunc := func(options *stscreds.AssumeRoleOptions) {
				options.Duration = assumeRoleOptions.Duration
				options.ExternalID = assumeRoleOptions.ExternalID
				options.RoleARN = assumeRoleOptions.RoleARN
				options.RoleSessionName = assumeRoleOptions.RoleSessionName
			}
			if assumeRoleOptions != nil {
				client := sts.NewFromConfig(cfg)
				assumeRoleProvider := stscreds.NewAssumeRoleProvider(
					client, assumeRoleOptions.RoleARN, assumeRoleOptionsFunc)
				cfg.Credentials = assumeRoleProvider
			}
		}

		if accessKey != "" && (secretKey != "" || sessionToken != "") {
			cfg.Credentials = awsCredentials.NewStaticCredentialsProvider(accessKey,
				secretKey, sessionToken)
		}
		regionalSns := sns.NewFromConfig(cfg)
		clumioRegion := getStringValue(d, "clumio_region")
		if clumioRegion != "" {
			cfg.Region = clumioRegion
		}
		s3obj := s3.NewFromConfig(cfg)
		return &apiClient{
			snsAPI: regionalSns,
			s3API:  s3obj,
			clumioConfig: clumioConfig.Config{
				Token:   apiToken,
				BaseUrl: baseUrl,
			},
		}, nil
	}
}

// Utility function to construct AssumeRoleOptions struct from the assumeRole parameter
func getAssumeRoleOptions(
	assumeRole interface{}) (*stscreds.AssumeRoleOptions, diag.Diagnostics) {
	assumeRoleMap, ok := assumeRole.(map[string]interface{})
	if !ok {
		return nil, diag.Errorf("Invalid format for assume_role")
	}
	roleArn := getStringValueFromMap(assumeRoleMap, "role_arn")
	var duration time.Duration
	if v, ok := assumeRoleMap["duration_seconds"].(int); ok && v != 0 {
		duration = time.Duration(v) * time.Second
	}
	externalId := getStringValueFromMap(assumeRoleMap, "external_id")
	sessionName := getStringValueFromMap(assumeRoleMap, "session_name")
	if sessionName == nil {
		empty := ""
		sessionName = &empty
	}
	return &stscreds.AssumeRoleOptions{
		RoleARN:         *roleArn,
		Duration:        duration,
		ExternalID:      externalId,
		RoleSessionName: *sessionName,
	}, nil
}
