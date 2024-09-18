// Code generated by smithy-go-codegen DO NOT EDIT.

package ecr

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates or updates the permissions policy for your registry.
//
// A registry policy is used to specify permissions for another Amazon Web
// Services account and is used when configuring cross-account replication. For
// more information, see [Registry permissions]in the Amazon Elastic Container Registry User Guide.
//
// [Registry permissions]: https://docs.aws.amazon.com/AmazonECR/latest/userguide/registry-permissions.html
func (c *Client) PutRegistryPolicy(ctx context.Context, params *PutRegistryPolicyInput, optFns ...func(*Options)) (*PutRegistryPolicyOutput, error) {
	if params == nil {
		params = &PutRegistryPolicyInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "PutRegistryPolicy", params, optFns, c.addOperationPutRegistryPolicyMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*PutRegistryPolicyOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type PutRegistryPolicyInput struct {

	// The JSON policy text to apply to your registry. The policy text follows the
	// same format as IAM policy text. For more information, see [Registry permissions]in the Amazon Elastic
	// Container Registry User Guide.
	//
	// [Registry permissions]: https://docs.aws.amazon.com/AmazonECR/latest/userguide/registry-permissions.html
	//
	// This member is required.
	PolicyText *string

	noSmithyDocumentSerde
}

type PutRegistryPolicyOutput struct {

	// The JSON policy text for your registry.
	PolicyText *string

	// The registry ID associated with the request.
	RegistryId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationPutRegistryPolicyMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpPutRegistryPolicy{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpPutRegistryPolicy{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "PutRegistryPolicy"); err != nil {
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
	if err = addOpPutRegistryPolicyValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opPutRegistryPolicy(options.Region), middleware.Before); err != nil {
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
	return nil
}

func newServiceMetadataMiddleware_opPutRegistryPolicy(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "PutRegistryPolicy",
	}
}
