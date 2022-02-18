// Copyright 2021. Clumio, Inc.

// clumio_callback_resource definition and CRUD implementation.

package clumio_callback

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
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
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

// ClumioCallback returns the resource for Clumio Callback. This resource is similar to
// the cloud formation custom resource. It will publish an event to the specified SNS
// topic and then wait for the status payload in the given S3 bucket.
func ClumioCallback() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Callback Resource used while on-boarding AWS clients." +
			" The purpose of this resource is to send a SNS event with the necessary" +
			" details of the AWS connection configuration done on the client AWS" +
			" account so that necessary connection post processing can be done in Clumio.",
		DeprecationMessage: "This resource is deprecated. " +
			"Use clumio_post_process_aws_connection resource instead.",
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
	client := meta.(*common.ApiClient)
	bucketName := d.Get("bucket_name").(string)
	accountId := fmt.Sprintf("%v", d.Get("account_id"))

	regionalSns := client.SnsAPI
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

	templateConfigs, err := common.GetTemplateConfiguration(d, false)
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
	s3obj := client.S3API
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
			var status common.StatusObject
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
