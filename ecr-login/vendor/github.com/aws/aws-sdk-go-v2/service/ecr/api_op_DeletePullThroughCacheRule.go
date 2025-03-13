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

// Deletes a pull through cache rule.
func (c *Client) DeletePullThroughCacheRule(ctx context.Context, params *DeletePullThroughCacheRuleInput, optFns ...func(*Options)) (*DeletePullThroughCacheRuleOutput, error) {
	if params == nil {
		params = &DeletePullThroughCacheRuleInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DeletePullThroughCacheRule", params, optFns, c.addOperationDeletePullThroughCacheRuleMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DeletePullThroughCacheRuleOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DeletePullThroughCacheRuleInput struct {

	// The Amazon ECR repository prefix associated with the pull through cache rule to
	// delete.
	//
	// This member is required.
	EcrRepositoryPrefix *string

	// The Amazon Web Services account ID associated with the registry that contains
	// the pull through cache rule. If you do not specify a registry, the default
	// registry is assumed.
	RegistryId *string

	noSmithyDocumentSerde
}

type DeletePullThroughCacheRuleOutput struct {

	// The timestamp associated with the pull through cache rule.
	CreatedAt *time.Time

	// The Amazon Resource Name (ARN) of the Amazon Web Services Secrets Manager
	// secret associated with the pull through cache rule.
	CredentialArn *string

	// The ARN of the IAM role associated with the pull through cache rule.
	CustomRoleArn *string

	// The Amazon ECR repository prefix associated with the request.
	EcrRepositoryPrefix *string

	// The registry ID associated with the request.
	RegistryId *string

	// The upstream registry URL associated with the pull through cache rule.
	UpstreamRegistryUrl *string

	// The upstream repository prefix associated with the pull through cache rule.
	UpstreamRepositoryPrefix *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDeletePullThroughCacheRuleMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpDeletePullThroughCacheRule{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpDeletePullThroughCacheRule{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DeletePullThroughCacheRule"); err != nil {
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
	if err = addOpDeletePullThroughCacheRuleValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDeletePullThroughCacheRule(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opDeletePullThroughCacheRule(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DeletePullThroughCacheRule",
	}
}
