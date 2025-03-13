// Code generated by smithy-go-codegen DO NOT EDIT.

package ecr

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"time"
)

// Updates an existing pull through cache rule.
func (c *Client) UpdatePullThroughCacheRule(ctx context.Context, params *UpdatePullThroughCacheRuleInput, optFns ...func(*Options)) (*UpdatePullThroughCacheRuleOutput, error) {
	if params == nil {
		params = &UpdatePullThroughCacheRuleInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdatePullThroughCacheRule", params, optFns, c.addOperationUpdatePullThroughCacheRuleMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdatePullThroughCacheRuleOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdatePullThroughCacheRuleInput struct {

	// The repository name prefix to use when caching images from the source registry.
	//
	// This member is required.
	EcrRepositoryPrefix *string

	// The Amazon Resource Name (ARN) of the Amazon Web Services Secrets Manager
	// secret that identifies the credentials to authenticate to the upstream registry.
	CredentialArn *string

	// Amazon Resource Name (ARN) of the IAM role to be assumed by Amazon ECR to
	// authenticate to the ECR upstream registry. This role must be in the same account
	// as the registry that you are configuring.
	CustomRoleArn *string

	// The Amazon Web Services account ID associated with the registry associated with
	// the pull through cache rule. If you do not specify a registry, the default
	// registry is assumed.
	RegistryId *string

	noSmithyDocumentSerde
}

type UpdatePullThroughCacheRuleOutput struct {

	// The Amazon Resource Name (ARN) of the Amazon Web Services Secrets Manager
	// secret associated with the pull through cache rule.
	CredentialArn *string

	// The ARN of the IAM role associated with the pull through cache rule.
	CustomRoleArn *string

	// The Amazon ECR repository prefix associated with the pull through cache rule.
	EcrRepositoryPrefix *string

	// The registry ID associated with the request.
	RegistryId *string

	// The date and time, in JavaScript date format, when the pull through cache rule
	// was updated.
	UpdatedAt *time.Time

	// The upstream repository prefix associated with the pull through cache rule.
	UpstreamRepositoryPrefix *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdatePullThroughCacheRuleMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpUpdatePullThroughCacheRule{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpUpdatePullThroughCacheRule{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdatePullThroughCacheRule"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addCredentialSource(stack, options); err != nil {
		return err
	}
	if err = addOpUpdatePullThroughCacheRuleValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdatePullThroughCacheRule(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opUpdatePullThroughCacheRule(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdatePullThroughCacheRule",
	}
}
