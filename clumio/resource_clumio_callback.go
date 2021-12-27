// Copyright 2021. Clumio, Inc.

// clumio_callback_resource definition and CRUD implementation.

package clumio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keyRegion           = "Region"
	keyToken            = "Token"
	keyType             = "Type"
	keyAccountID        = "AccountId"
	keyRoleID           = "RoleId"
	keyRoleArn          = "RoleArn"
	keyExternalID       = "RoleExternalId"
	keyClumioEventPubID = "ClumioEventPubId"
	keyCanonicalUser    = "CanonicalUser"
	keyTemplateConfig   = "TemplateConfiguration"
	keyEventPublishTime = "EventPublishTime"

	// Number of retries that we will perform before giving up a AWS request.
	kMaxRetries       = 8
	requestTypeCreate = "Create"
	requestTypeDelete = "Delete"
	requestTypeUpdate = "Update"

	//Status strings
	statusFailed  = "FAILED"
	statusSuccess = "SUCCESS"

	//bucket key format
	bucketKeyFormat = "acmtfstatus/%s/%s/%s/clumio-status-%s.json"

	//Error Strings
	movedPermanently = "MovedPermanently"
	snsPublishError  = "operation error SNS: Publish, https response error StatusCode: 403"
)

var (
	// protectInfoMap is the mapping of the the datasource to the resource parameter and
	// if a config section is required, then isConfig will be true.
	protectInfoMap = map[string]sourceConfigInfo{
		"ebs": {
			sourceKey: "protect_ebs_version",
			isConfig:  false,
		},
		"rds": {
			sourceKey: "protect_rds_version",
			isConfig:  false,
		},
		"ec2_mssql": {
			sourceKey: "protect_ec2_mssql_version",
			isConfig:  false,
		},
		"warm_tier": {
			sourceKey: "protect_warm_tier_version",
			isConfig:  true,
		},
		"s3": {
			sourceKey: "protect_s3_version",
			isConfig:  false,
		},
		"dynamodb": {
			sourceKey: "protect_dynamodb_version",
			isConfig:  false,
		},
	}
	// warmtierInfoMap is the mapping of the the warm tier datasource to the resource
	// parameter and if a config section is required, then isConfig will be true.
	warmtierInfoMap = map[string]sourceConfigInfo{
		"dynamodb": {
			sourceKey: "protect_warm_tier_dynamodb_version",
			isConfig:  false,
		},
	}
)

// SNSEvent is the event payload to be sent to the topic
type SNSEvent struct {
	RequestType        string                 `json:"RequestType"`
	ServiceToken       string                 `json:"ServiceToken"`
	ResponseURL        string                 `json:"ResponseURL"`
	StackID            string                 `json:"StackId"`
	RequestID          string                 `json:"RequestId"`
	LogicalResourceID  string                 `json:"LogicalResourceId"`
	ResourceType       string                 `json:"ResourceType"`
	ResourceProperties map[string]interface{} `json:"ResourceProperties"`
}

// The payload in the status file read from S3.
type StatusObject struct {
	Status string            `json:"Status"`
	Reason *string           `json:"Reason,omitempty"`
	Data   map[string]string `json:"Data,omitempty"`
}

// clumioCallback returns the resource for Clumio Callback. This resource is similar to
// the cloud formation custom resource. It will publish an event to the specified SNS
// topic and then wait for the status payload in the given S3 bucket.
func clumioCallback() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Callback Resource used while on-boarding AWS clients." +
			" The purpose of this resource is to send a SNS event with the necessary" +
			" details of the AWS connection configuration done on the client AWS" +
			" account so that necessary connection post processing can be done in Clumio.",

		CreateContext: clumioCallbackCreate,
		ReadContext:   clumioCallbackRead,
		UpdateContext: clumioCallbackUpdate,
		DeleteContext: clumioCallbackDelete,

		Schema: map[string]*schema.Schema{
			"sns_topic": {
				Description: "SNS Topic to publish event.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"token": {
				Type:        schema.TypeString,
				Description: "The AWS integration ID token.",
				Required:    true,
			},
			"role_external_id": {
				Type: schema.TypeString,
				Description: "A key that must be used by Clumio to assume the service role" +
					" in your account. This should be a secure string, like a password," +
					" but it does not need to be remembered (random characters are best).",
				Required: true,
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The AWS Customer Account ID.",
				Required:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The AWS Region.",
				Required:    true,
			},
			"role_id": {
				Type:        schema.TypeString,
				Description: "Clumio IAM Role ID.",
				Required:    true,
			},
			"role_arn": {
				Type:        schema.TypeString,
				Description: "Clumio IAM Role Arn.",
				Required:    true,
			},
			"clumio_event_pub_id": {
				Type:        schema.TypeString,
				Description: "Clumio Event Pub SNS topic ID.",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Registration Type.",
				Required:    true,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Description: "S3 bucket name where the status file is written.",
				Required:    true,
			},
			"canonical_user": {
				Type:        schema.TypeString,
				Description: "Canonical User ID of the account.",
				Required:    true,
			},
			"config_version": {
				Type:        schema.TypeString,
				Description: "Clumio Config version.",
				Required:    true,
			},
			"discover_enabled": {
				Type:        schema.TypeBool,
				Description: "Is Clumio Discover enabled.",
				Optional:    true,
			},
			"discover_version": {
				Type:        schema.TypeString,
				Description: "Clumio Discover version.",
				Required:    true,
			},
			"protect_enabled": {
				Type:        schema.TypeBool,
				Description: "Is Clumio Protect enabled.",
				Optional:    true,
			},
			"protect_config_version": {
				Type:        schema.TypeString,
				Description: "Clumio Protect Config version.",
				Optional:    true,
			},
			"protect_ebs_version": {
				Type:        schema.TypeString,
				Description: "Clumio EBS Protect version.",
				Optional:    true,
			},
			"protect_rds_version": {
				Type:        schema.TypeString,
				Description: "Clumio RDS Protect version.",
				Optional:    true,
			},
			"protect_ec2_mssql_version": {
				Type:        schema.TypeString,
				Description: "Clumio EC2 MSSQL Protect version.",
				Optional:    true,
			},
			"protect_s3_version": {
				Type:        schema.TypeString,
				Description: "Clumio S3 Protect version.",
				Optional:    true,
			},
			"protect_dynamodb_version": {
				Type:        schema.TypeString,
				Description: "Clumio DynamoDB Protect version.",
				Optional:    true,
			},
			"protect_warm_tier_version": {
				Type:        schema.TypeString,
				Description: "Clumio Warmtier Protect version.",
				Optional:    true,
			},
			"protect_warm_tier_dynamodb_version": {
				Type:        schema.TypeString,
				Description: "Clumio DynamoDB Warmtier Protect version.",
				Optional:    true,
			},
			"properties": {
				Type:        schema.TypeMap,
				Description: "Properties to be passed in the SNS event.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// clumioCallbackCreate handles the Create action for the Clumio Callback Resource.
func clumioCallbackCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioCallbackCommon(ctx, d, meta, requestTypeCreate)
}

// clumioCallbackCreate handles the Read action for the Clumio Callback Resource.
func clumioCallbackRead(_ context.Context, _ *schema.ResourceData,
	_ interface{}) diag.Diagnostics {
	// Nothing to Read for this resource
	return nil
}

// clumioCallbackCreate handles the Update action for the Clumio Callback Resource.
func clumioCallbackUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioCallbackCommon(ctx, d, meta, requestTypeUpdate)
}

// clumioCallbackCreate handles the Delete action for the Clumio Callback Resource.
func clumioCallbackDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioCallbackCommon(ctx, d, meta, requestTypeDelete)
}

// clumioCallbackCommon will construct the event payload from the resource properties and
// publish the event to the SNS topic and then wait for the status payload in the given
// S3 bucket.
func clumioCallbackCommon(ctx context.Context, d *schema.ResourceData, meta interface{},
	eventType string) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)
	bucketName := d.Get("bucket_name").(string)
	accountId := fmt.Sprintf("%v", d.Get("account_id"))

	regionalSns := client.snsAPI
	event := SNSEvent{
		RequestType:        eventType,
		ServiceToken:       fmt.Sprintf("%v", d.Get("token")),
		ResourceProperties: nil,
	}
	token := fmt.Sprintf("%v", d.Get("token"))
	region := fmt.Sprintf("%v", d.Get("region"))
	resourceProperties := make(map[string]interface{})
	resourceProperties[keyAccountID] = fmt.Sprintf("%v", d.Get("account_id"))
	resourceProperties[keyToken] = token
	resourceProperties[keyType] = fmt.Sprintf("%v", d.Get("type"))
	resourceProperties[keyAccountID] = accountId
	resourceProperties[keyRegion] = region
	resourceProperties[keyRoleID] = fmt.Sprintf("%v", d.Get("role_id"))
	resourceProperties[keyRoleArn] = fmt.Sprintf("%v", d.Get("role_arn"))
	resourceProperties[keyExternalID] =
		fmt.Sprintf("%v", d.Get("role_external_id"))
	resourceProperties[keyClumioEventPubID] =
		fmt.Sprintf("%v", d.Get("clumio_event_pub_id"))
	resourceProperties[keyCanonicalUser] = fmt.Sprintf("%v", d.Get("canonical_user"))

	templateConfigs, err := getTemplateConfiguration(d, false)
	if err != nil {
		return diag.Errorf("Error forming template configuration. Error: %v", err)
	}
	resourceProperties[keyTemplateConfig] = templateConfigs
	if val, ok := d.GetOk("properties"); ok && len(val.(map[string]interface{})) > 0 {
		properties := val.(map[string]interface{})
		for key, value := range properties {
			resourceProperties[key] = value.(string)
		}
	}
	startTime := time.Now()
	startTimeUnixStr := strconv.FormatInt(
		startTime.UnixNano()/int64(time.Millisecond), 10)
	resourceProperties[keyEventPublishTime] = startTimeUnixStr
	event.ResourceProperties = resourceProperties
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return diag.Errorf("Error occurred in marshalling event: %v", err)
	}
	// Publish event to SNS.
	publishInput := &sns.PublishInput{
		Message:  aws.String(string(eventBytes)),
		TopicArn: aws.String(fmt.Sprintf("%v", d.Get("sns_topic"))),
	}
	_, err = regionalSns.Publish(ctx, publishInput)
	if err != nil {
		return diag.Errorf("Error occurred in SNS Publish Request: %v",
			processErrorMessage(err.Error(), region, token))
	}
	if eventType == requestTypeCreate {
		d.SetId(uuid.New().String())
	}
	s3obj := client.s3API
	endTime := startTime.Add(5 * time.Minute)
	timeOut := false
	processingDone := false
	// Poll the s3 bucket for the clumio-status.json file. Keep retrying every 5 seconds
	// till the last modified time on the file is greater than the startTime and less than
	// the end time.
	for {
		if time.Now().After(endTime) {
			timeOut = true
			break
		}
		time.Sleep(5 * time.Second)
		objectKey := fmt.Sprintf(
			bucketKeyFormat, accountId, region, token, startTimeUnixStr)
		// HeadObject call to get the last modified time of the file.
		statusObj, err := s3obj.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			var aerr smithy.APIError
			if errors.As(err, &aerr) {
				// Checking for both forbidden as well as NoSuchKey as acceptance test
				// returns NoSuchKey while actual resource returns Forbidden error.
				_, ok := aerr.(*s3types.NoSuchKey)
				if aerr.ErrorCode() == "Forbidden" || ok {
					log.Println(aerr.Error())
					continue
				}
				return diag.Errorf("Error retrieving clumio-status.json. "+
					"Error Code : %v, Error message: %v, origError: %v",
					aerr.ErrorCode(), aerr.ErrorMessage(), err)
			}
			return diag.Errorf("Error retrieving clumio-status.json: %v", err)
		} else {
			var status StatusObject
			statusObjBytes := new(bytes.Buffer)
			_, err = statusObjBytes.ReadFrom(statusObj.Body)
			if err != nil {
				return diag.Errorf("Error reading status object: %v", err)
			}
			err = json.Unmarshal(statusObjBytes.Bytes(), &status)
			if err != nil {
				return diag.Errorf("Error unmarshalling status object: %v", err)
			}
			if status.Status == statusFailed {
				return diag.Errorf("Processing of Clumio Event failed. "+
					"Error Message : %s",
					processErrorMessage(*status.Reason, region, token))
			} else if status.Status == statusSuccess {
				processingDone = true
				break
			}
		}
	}
	if !processingDone && timeOut {
		return diag.Errorf("Timeout occurred waiting for status.")
	}
	return nil
}

type sourceConfigInfo struct {
	sourceKey string
	isConfig  bool
}

// getTemplateConfiguration returns the template configuration.
func getTemplateConfiguration(d *schema.ResourceData, isCamelCase bool) (
	map[string]interface{}, error) {
	templateConfigs := make(map[string]interface{})
	configMap, err := getConfigMapForKey(d, "config_version", false)
	if err != nil {
		return nil, err
	}
	if configMap == nil {
		return templateConfigs, nil
	}
	templateConfigs["config"] = configMap
	discoverMap, err := getConfigMapForKey(d, "discover_version", true)
	if err != nil {
		return nil, err
	}
	if discoverMap == nil {
		return templateConfigs, nil
	}
	templateConfigs["discover"] = discoverMap
	protectMap, err := getConfigMapForKey(d, "protect_config_version", true)
	if err != nil {
		return nil, err
	}
	if protectMap == nil {
		return templateConfigs, nil
	}
	err = populateConfigMap(d, protectInfoMap, protectMap, isCamelCase)
	if err != nil {
		return nil, err
	}
	warmTierKey := "warm_tier"
	if isCamelCase {
		warmTierKey = "warmTier"
	}
	if protectWarmtierMap, ok := protectMap[warmTierKey]; ok {
		err = populateConfigMap(
			d, warmtierInfoMap, protectWarmtierMap.(map[string]interface{}), isCamelCase)
		if err != nil {
			return nil, err
		}
	}
	templateConfigs["protect"] = protectMap
	return templateConfigs, nil
}

// populateConfigMap returns protect configuration information for the configs
// in the configInfoMap.
func populateConfigMap(d *schema.ResourceData, configInfoMap map[string]sourceConfigInfo,
	configMap map[string]interface{}, isCamelCase bool) error {
	for source, sourceInfo := range configInfoMap {
		configMapKey := source
		if isCamelCase {
			configMapKey = snakeCaseToCamelCase(source)
		}
		protectSourceMap, err := getConfigMapForKey(
			d, sourceInfo.sourceKey, sourceInfo.isConfig)
		if err != nil {
			return err
		}
		if protectSourceMap != nil {
			configMap[configMapKey] = protectSourceMap
		}
	}
	return nil
}

// getConfigMapForKey returns a config map for the key if it exists in ResourceData.
func getConfigMapForKey(
	d *schema.ResourceData, key string, isConfig bool) (map[string]interface{}, error) {
	var mapToReturn map[string]interface{}
	if val, ok := d.GetOk(key); ok {
		keyMap := make(map[string]interface{})
		if keyVersion, ok := val.(string); ok {
			majorVersion, minorVersion, err := parseVersion(keyVersion)
			if err != nil {
				return nil, err
			}
			keyMap["enabled"] = true
			keyMap["version"] = majorVersion
			if minorVersion != "" {
				keyMap["minorVersion"] = minorVersion
			}

		}
		mapToReturn = keyMap
		// If isConfig is true it wraps the keyMap with another map with "config" as the key.
		if isConfig {
			configMap := make(map[string]interface{})
			configMap["config"] = keyMap
			mapToReturn = configMap
		}
	}
	return mapToReturn, nil
}

// parseVersion parses the version and minorVersion given the version string.
func parseVersion(version string) (string, string, error) {
	splits := strings.Split(version, ".")
	switch len(splits) {
	case 1:
		return version, "", nil
	case 2:
		return splits[0], splits[1], nil
	default:
		return "", "", errors.New(fmt.Sprintf("Invalid version: %v", version))
	}
}

// processErrorMessage takes the failure reason and adds the potential cause for the
// failure.
func processErrorMessage(message string, region string, token string) string {
	if strings.Contains(message, movedPermanently) {
		return fmt.Sprintf("Incorrect region specified : %s", region)
	}
	if strings.Contains(message, snsPublishError) {
		return fmt.Sprintf(
			"SNS Publish Error. Incorrect region or clumio_token specified : %s", token)
	}

	return message
}
