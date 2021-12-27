// Copyright 2021. Clumio, Inc.

// Acceptance test for resource_clumio_callback.
package clumio

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

func init() {
	resource.AddTestSweepers("clumio_callback_resource", &resource.Sweeper{
		Name: "clumio_callback_resource",
		F:    testSweepAWSResources,
	})
}

func testSweepAWSResources(region string) error {
	return destroyResources(nil)
}

func TestAccResourceClumioCallback(t *testing.T) {
	if isAcceptanceTest() {
		accountId, topicArn, canonicalUser, err := setUpResources(t)
		require.Nil(t, err)
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheckAws(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: getTestAccResourceClumioCallback(
						topicArn, accountId, canonicalUser),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(
							"clumio_callback_resource.test", "sns_topic",
							regexp.MustCompile(topicArn)),
					),
				},
			},
		})
		err = destroyResources(aws.String(topicArn))
		require.Nil(t, err)
	} else {
		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: getTestAccResourceClumioCallback(
						"topicArn", "accountId", "canonicalUser"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(
							"clumio_callback_resource.test", "sns_topic",
							regexp.MustCompile("topicArn")),
					),
				},
			},
		})
	}

}

// setUpResources creates the AWS resources required for the acceptance test
func setUpResources(t *testing.T) (string, string, string, error) {
	testAccPreCheckAws(t)
	log.Println("Creating acceptance test resources for clumio_callback_resource.")
	ctx := context.TODO()
	cfg, err := getAWSConfig(ctx)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	stsClient := sts.NewFromConfig(cfg)
	getCallerIdentityOut, err := stsClient.GetCallerIdentity(
		ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	s3Client := s3.NewFromConfig(cfg)
	listBucketsOut, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	canonicalUser := listBucketsOut.Owner.ID
	bucketName := fmt.Sprintf("clumio-tf-acc-test-bucket-%s", *getCallerIdentityOut.Account)
	err = createAccTestS3Bucket(ctx, s3Client, bucketName)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	roleArn, err := createAccTestRole(ctx, cfg, bucketName)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// Create SNS topic
	snsClient := sns.NewFromConfig(cfg)
	createTopicOut, err := snsClient.CreateTopic(ctx, &sns.CreateTopicInput{
		Name: aws.String("clumio-tf-acc-test-topic"),
	})
	if err != nil {
		return "", "", "", err
	}

	functionArn, err := createAccTestLambda(
		ctx, cfg, roleArn, bucketName, createTopicOut.TopicArn)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	_, err = snsClient.Subscribe(ctx, &sns.SubscribeInput{
		Protocol:              aws.String("lambda"),
		TopicArn:              createTopicOut.TopicArn,
		Attributes:            nil,
		Endpoint:              functionArn,
		ReturnSubscriptionArn: false,
	})
	return *getCallerIdentityOut.Account, *createTopicOut.TopicArn, *canonicalUser, err
}

// destroyResources cleans up the AWS resources created for the acceptance test.
func destroyResources(topicArn *string) error {
	log.Println("Destroying AWS resources created during the acceptance test for clumio_callback_resource.")
	ctx := context.TODO()
	cfg, err := getAWSConfig(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	stsClient := sts.NewFromConfig(cfg)
	getCallerIdentityOut, err := stsClient.GetCallerIdentity(
		ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Println(err)
		return err
	}

	err = deleteAccTestRole(ctx, cfg)
	if err != nil {
		log.Println(err)
		return err
	}

	err = deleteAccTestTopic(ctx, cfg, topicArn, *getCallerIdentityOut.Account)
	if err != nil {
		log.Println(err)
		return err
	}

	err = deleteAccTestLambda(ctx, cfg)
	if err != nil {
		log.Println(err)
		return err
	}

	err = deleteAccTestS3Bucket(ctx, cfg, *getCallerIdentityOut.Account)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// function to return the tf script snippet to be used in the acceptance/unit test
func getTestAccResourceClumioCallback(
	topicArn string, accountId string, canonicalUser string) string {
	return fmt.Sprintf(
		testAccResourceClumioCallback, topicArn, accountId, accountId, canonicalUser)
}

// getAWSConfig constructs and returns the AWS config
func getAWSConfig(ctx context.Context) (aws.Config, error) {
	loadOpts := []func(options *config.LoadOptions) error{
		config.WithRetryer(func() aws.Retryer {
			return retry.AddWithMaxAttempts(retry.NewStandard(), kMaxRetries)
		}),
	}
	profile := os.Getenv("AWS_PROFILE")
	if profile != "" {
		loadOpts = append(loadOpts, config.WithSharedConfigProfile(profile))
	}

	sharedCredsFile := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	if sharedCredsFile != "" {
		loadOpts = append(
			loadOpts, config.WithSharedCredentialsFiles([]string{sharedCredsFile}))
	}
	cfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		log.Println(err)
		return aws.Config{}, err
	}
	return cfg, nil
}

const testAccResourceClumioCallback = `
provider clumio{}
resource "clumio_callback_resource" "test" {
  sns_topic = "%s"
  token = "89130d52-67cb-4752-a896-730cf067aeb1"
  role_external_id = "Aasdfjhg8943kbnlasdklhkljghlkash7892r"
  account_id = "%s"
  region = "us-west-2"
  role_id = "TestRole-us-west-2-89130d52-67cb-4752-a896-730cf067aeb1"
  role_arn = "arn:aws:iam::482567874266:role/clumio/TestRole-us-west-2-89130d52-67cb-4752-a896-730cf067aeb1"
  clumio_event_pub_id = "arn:aws:sns:us-west-2:482567874266:ClumioInventoryTopic_89130d52-67cb-4752-a896-730cf067aeb1"
  type = "service"
  bucket_name = "clumio-tf-acc-test-bucket-%s"
  canonical_user = "%s"
  config_version = "1"
  discover_version = "3"
  protect_config_version = "18"
  protect_ebs_version = "19"
  protect_rds_version = "18"
  protect_ec2_mssql_version = "1"
  protect_warm_tier_version = "2"
  protect_warm_tier_dynamodb_version = "2"
  protect_dynamodb_version = "1"
}
`

// StatementEntry dictates what this policy allows or doesn't allow.
type StatementEntry struct {
	Effect   string
	Action   []string
	Resource string
}

// PolicyDocument is our definition of our policies to be uploaded to AWS Identity and Access Management (IAM).
type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

// CreateLambdaS3PolicyDoc creates a policy document that contains permissions to read from
// and write to a S3 bucket.
func CreateLambdaS3PolicyDoc(bucketName string) ([]byte, error) {
	policy := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"s3:PutObject",
					"s3:GetObjectAcl",
					"s3:GetObject",
					"s3:PutObjectAcl",
				},
				Resource: fmt.Sprintf("arn:aws:s3:::%s/*", bucketName),
			},
		},
	}
	b, err := json.Marshal(&policy)
	return b, err
}

// createAccTestS3Bucket creates a S3 bucket which is used for acceptance test.
func createAccTestS3Bucket(
	ctx context.Context, s3Client *s3.Client, bucketName string) error {
	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraintUsWest2,
		},
	})
	if err != nil {
		var aerr smithy.APIError
		if errors.As(err, &aerr) {
			_, ok := aerr.(*s3types.BucketAlreadyOwnedByYou)
			if ok {
				return nil
			}
		}
		log.Println(err)
		return err
	}
	return nil
}

// createAccTestRole creates an IAM role along with a role policy which allows a lambda
// to write to a S3 bucket.
func createAccTestRole(ctx context.Context, cfg aws.Config, bucketName string) (
	*string, error) {
	// create policy for lambda to write to s3
	iamClient := iam.NewFromConfig(cfg)
	assumeRolePolicy := `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": ["sts:AssumeRole"], "Principal": {"Service": "lambda.amazonaws.com"}}]}`
	policyDoc, err := CreateLambdaS3PolicyDoc(bucketName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	createRoleOut, err := iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String("clumio-tf-acc-test-role"),
		Description:              aws.String("Clumio TF acceptance test role"),
		AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
	})
	var roleArn *string
	if err != nil {
		var aerr smithy.APIError
		if errors.As(err, &aerr) {
			_, ok := aerr.(*iamtypes.EntityAlreadyExistsException)
			if !ok {
				log.Println(err)
				return nil, err
			}
		}
		getRoleOut, err := iamClient.GetRole(
			ctx, &iam.GetRoleInput{RoleName: aws.String("clumio-tf-acc-test-role")})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		roleArn = getRoleOut.Role.Arn
	} else {
		roleArn = createRoleOut.Role.Arn
	}
	_, err = iamClient.PutRolePolicy(ctx, &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(string(policyDoc)),
		PolicyName:     aws.String("clumio-tf-acc-test-role-policy"),
		RoleName:       aws.String("clumio-tf-acc-test-role"),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = iamClient.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		RoleName: aws.String("clumio-tf-acc-test-role"),
		PolicyArn: aws.String(
			"arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return roleArn, nil
}

// createAccTestLambda creates a lambda along with the necessary permissions to read an
// event from SNS topic and write a status message to S3 bucket.
func createAccTestLambda(ctx context.Context, cfg aws.Config, roleArn *string,
	bucketName string, topicArn *string) (*string, error) {
	//Create the zip file to be uploaded
	archive, err := os.Create("test/clumio-tf-acc-test-lambda.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	lambda_file, err := os.Open("test/acc_test_lambda.py")
	if err != nil {
		panic(err)
	}
	defer lambda_file.Close()
	w1, err := zipWriter.Create("acc_test_lambda.py")
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(w1, lambda_file); err != nil {
		panic(err)
	}
	if err = zipWriter.Close(); err != nil {
		return nil, err
	}
	time.Sleep(10 * time.Second)
	lambdaClient := lambda.NewFromConfig(cfg)
	fileData, err := os.ReadFile("test/clumio-tf-acc-test-lambda.zip")
	if err != nil {
		return nil, err
	}
	lambdaOut, err := lambdaClient.CreateFunction(ctx, &lambda.CreateFunctionInput{
		Code: &types.FunctionCode{
			ZipFile: fileData,
		},
		FunctionName: aws.String("clumio-tf-acc-test-lambda"),
		Role:         roleArn,
		Description:  aws.String("Clumio TF provider acceptance test lambda"),
		MemorySize:   nil,
		PackageType:  "",
		Publish:      true,
		Runtime:      "python3.9",
		Timeout:      aws.Int32(60),
		Handler:      aws.String("acc_test_lambda.clumio_event_handler"),
		Environment: &types.Environment{
			Variables: map[string]string{
				"BUCKET_NAME": bucketName,
			},
		},
	})
	var functionArn *string
	var aerr smithy.APIError
	if err != nil {
		if errors.As(err, &aerr) {
			_, ok := aerr.(*types.ResourceConflictException)
			if !ok {
				log.Println(err)
				return nil, err
			}
		}
		getLambdaOut, err := lambdaClient.GetFunction(ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String("clumio-tf-acc-test-lambda"),
		})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		functionArn = getLambdaOut.Configuration.FunctionArn
	} else {
		functionArn = lambdaOut.FunctionArn
	}
	_, err = lambdaClient.AddPermission(ctx, &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: functionArn,
		Principal:    aws.String("sns.amazonaws.com"),
		StatementId:  aws.String("AllowExecutionFromSNS"),
		SourceArn:    topicArn,
	})
	if err != nil {
		if errors.As(err, &aerr) {
			_, ok := aerr.(*types.ResourceConflictException)
			if !ok {
				log.Println(err)
				return nil, err
			}
		}
	}
	err = os.Remove("test/clumio-tf-acc-test-lambda.zip")
	return functionArn, err
}

// deleteAccTestRole deletes the AWS IAM role created for the acceptance test.
func deleteAccTestRole(ctx context.Context, cfg aws.Config) error {
	iamClient := iam.NewFromConfig(cfg)
	listRolePoliciesOut, err := iamClient.ListRolePolicies(ctx, &iam.ListRolePoliciesInput{
		RoleName: aws.String("clumio-tf-acc-test-role"),
	})
	if err != nil {
		var aerr smithy.APIError
		if errors.As(err, &aerr) {
			_, ok := aerr.(*iamtypes.NoSuchEntityException)
			if ok {
				return nil
			}
		}
		log.Println(err)
		return err
	}
	for _, policyName := range listRolePoliciesOut.PolicyNames {
		_, err = iamClient.DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
			PolicyName: aws.String(policyName),
			RoleName:   aws.String("clumio-tf-acc-test-role"),
		})
		if err != nil {
			log.Println(err)
			return err
		}
	}

	_, err = iamClient.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
		RoleName: aws.String("clumio-tf-acc-test-role"),
		PolicyArn: aws.String(
			"arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = iamClient.DeleteRole(ctx, &iam.DeleteRoleInput{
		RoleName: aws.String("clumio-tf-acc-test-role"),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// deleteAccTestTopic deletes the AWS SNS topic created for the acceptance test.
func deleteAccTestTopic(
	ctx context.Context, cfg aws.Config, topicArn *string, accountId string) error {
	snsClient := sns.NewFromConfig(cfg)
	if topicArn == nil || *topicArn == "" {
		topicArn = aws.String(fmt.Sprintf(
			"arn:aws:sns:%s:%s:clumio-tf-acc-test-topic", cfg.Region, accountId))
	}
	_, err := snsClient.DeleteTopic(ctx, &sns.DeleteTopicInput{
		TopicArn: topicArn,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// deleteAccTestTopic deletes the AWS lambda function created for the acceptance test.
func deleteAccTestLambda(ctx context.Context, cfg aws.Config) error {
	lambdaClient := lambda.NewFromConfig(cfg)
	_, err := lambdaClient.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
		FunctionName: aws.String("clumio-tf-acc-test-lambda"),
	})
	if err != nil {
		var aerr smithy.APIError
		if errors.As(err, &aerr) {
			_, ok := aerr.(*types.ResourceNotFoundException)
			if ok {
				return nil
			}
		}
		log.Println(err)
		return err
	}
	return nil
}

// deleteAccTestTopic deletes the AWS S3 bucket created for the acceptance test.
func deleteAccTestS3Bucket(ctx context.Context, cfg aws.Config, accountId string) error {
	s3Client := s3.NewFromConfig(cfg)
	bucketName := fmt.Sprintf("clumio-tf-acc-test-bucket-%s", accountId)
	listObjsOut, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var aerr smithy.APIError
		if errors.As(err, &aerr) {
			_, ok := aerr.(*s3types.NoSuchBucket)
			if ok {
				return nil
			}
		}
		log.Println(err)
		return err
	}

	objIdentifiers := make([]s3types.ObjectIdentifier, 0)
	for _, obj := range listObjsOut.Contents {
		objIdentifiers = append(objIdentifiers, s3types.ObjectIdentifier{
			Key: obj.Key,
		})
	}
	if len(objIdentifiers) > 0 {
		_, err = s3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &s3types.Delete{
				Objects: objIdentifiers,
				Quiet:   false,
			},
		})
	}
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
