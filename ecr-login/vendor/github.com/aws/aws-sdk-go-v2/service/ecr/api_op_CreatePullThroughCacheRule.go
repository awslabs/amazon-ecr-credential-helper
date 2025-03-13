// Code generated by smithy-go-codegen DO NOT EDIT.

package ecr

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"time"
)

// Creates a pull through cache rule. A pull through cache rule provides a way to
// cache images from an upstream registry source in your Amazon ECR private
// registry. For more information, see [Using pull through cache rules]in the Amazon Elastic Container Registry
// User Guide.
//
// [Using pull through cache rules]: https://docs.aws.amazon.com/AmazonECR/latest/userguide/pull-through-cache.html
func (c *Client) CreatePullThroughCacheRule(ctx context.Context, params *CreatePullThroughCacheRuleInput, optFns ...func(*Options)) (*CreatePullThroughCacheRuleOutput, error) {
	if params == nil {
		params = &CreatePullThroughCacheRuleInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreatePullThroughCacheRule", params, optFns, c.addOperationCreatePullThroughCacheRuleMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreatePullThroughCacheRuleOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreatePullThroughCacheRuleInput struct {

	// The repository name prefix to use when caching images from the source registry.
	//
	// There is always an assumed / applied to the end of the prefix. If you specify
	// ecr-public as the prefix, Amazon ECR treats that as ecr-public/ .
	//
	// This member is required.
	EcrRepositoryPrefix *string

	// The registry URL of the upstream public registry to use as the source for the
	// pull through cache rule. The following is the syntax to use for each supported
	// upstream registry.
	//
	//   - Amazon ECR ( ecr ) – dkr.ecr..amazonaws.com
	//
	//   - Amazon ECR Public ( ecr-public ) – public.ecr.aws
	//
	//   - Docker Hub ( docker-hub ) – registry-1.docker.io
	//
	//   - GitHub Container Registry ( github-container-registry ) – ghcr.io
	//
	//   - GitLab Container Registry ( gitlab-container-registry ) –
	//   registry.gitlab.com
	//
	//   - Kubernetes ( k8s ) – registry.k8s.io
	//
	//   - Microsoft Azure Container Registry ( azure-container-registry ) –
	//   .azurecr.io
	//
	//   - Quay ( quay ) – quay.io
	//
	// This member is required.
	UpstreamRegistryUrl *string

	// The Amazon Resource Name (ARN) of the Amazon Web Services Secrets Manager
	// secret that identifies the credentials to authenticate to the upstream registry.
	CredentialArn *string

	// Amazon Resource Name (ARN) of the IAM role to be assumed by Amazon ECR to
	// authenticate to the ECR upstream registry. This role must be in the same account
	// as the registry that you are configuring.
	CustomRoleArn *string

	// The Amazon Web Services account ID associated with the registry to create the
	// pull through cache rule for. If you do not specify a registry, the default
	// registry is assumed.
	RegistryId *string

	// The name of the upstream registry.
	UpstreamRegistry types.UpstreamRegistry

	// The repository name prefix of the upstream registry to match with the upstream
	// repository name. When this field isn't specified, Amazon ECR will use the ROOT .
	UpstreamRepositoryPrefix *string

	noSmithyDocumentSerde
}

type CreatePullThroughCacheRuleOutput struct {

	// The date and time, in JavaScript date format, when the pull through cache rule
	// was created.
	CreatedAt *time.Time

	// The Amazon Resource Name (ARN) of the Amazon Web Services Secrets Manager
	// secret associated with the pull through cache rule.
	CredentialArn *string

	// The ARN of the IAM role associated with the pull through cache rule.
	CustomRoleArn *string

	// The Amazon ECR repository prefix associated with the pull through cache rule.
	EcrRepositoryPrefix *string

	// The registry ID associated with the request.
	RegistryId *string

	// The name of the upstream registry associated with the pull through cache rule.
	UpstreamRegistry types.UpstreamRegistry

	// The upstream registry URL associated with the pull through cache rule.
	UpstreamRegistryUrl *string

	// The upstream repository prefix associated with the pull through cache rule.
	UpstreamRepositoryPrefix *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreatePullThroughCacheRuleMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpCreatePullThroughCacheRule{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpCreatePullThroughCacheRule{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CreatePullThroughCacheRule"); err != nil {
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
	if err = addOpCreatePullThroughCacheRuleValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreatePullThroughCacheRule(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opCreatePullThroughCacheRule(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CreatePullThroughCacheRule",
	}
}
