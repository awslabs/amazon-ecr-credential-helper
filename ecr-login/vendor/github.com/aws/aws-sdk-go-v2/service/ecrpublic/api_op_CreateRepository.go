// Code generated by smithy-go-codegen DO NOT EDIT.

package ecrpublic

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates a repository in a public registry. For more information, see [Amazon ECR repositories] in the
// Amazon Elastic Container Registry User Guide.
//
// [Amazon ECR repositories]: https://docs.aws.amazon.com/AmazonECR/latest/userguide/Repositories.html
func (c *Client) CreateRepository(ctx context.Context, params *CreateRepositoryInput, optFns ...func(*Options)) (*CreateRepositoryOutput, error) {
	if params == nil {
		params = &CreateRepositoryInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateRepository", params, optFns, c.addOperationCreateRepositoryMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateRepositoryOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateRepositoryInput struct {

	// The name to use for the repository. This appears publicly in the Amazon ECR
	// Public Gallery. The repository name can be specified on its own (for example
	// nginx-web-app ) or prepended with a namespace to group the repository into a
	// category (for example project-a/nginx-web-app ).
	//
	// This member is required.
	RepositoryName *string

	// The details about the repository that are publicly visible in the Amazon ECR
	// Public Gallery.
	CatalogData *types.RepositoryCatalogDataInput

	// The metadata that you apply to each repository to help categorize and organize
	// your repositories. Each tag consists of a key and an optional value. You define
	// both of them. Tag keys can have a maximum character length of 128 characters,
	// and tag values can have a maximum length of 256 characters.
	Tags []types.Tag

	noSmithyDocumentSerde
}

type CreateRepositoryOutput struct {

	// The catalog data for a repository. This data is publicly visible in the Amazon
	// ECR Public Gallery.
	CatalogData *types.RepositoryCatalogData

	// The repository that was created.
	Repository *types.Repository

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateRepositoryMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpCreateRepository{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpCreateRepository{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CreateRepository"); err != nil {
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
	if err = addOpCreateRepositoryValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateRepository(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opCreateRepository(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CreateRepository",
	}
}
