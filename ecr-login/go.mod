module github.com/awslabs/amazon-ecr-credential-helper/ecr-login

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.48.0
	github.com/aws/aws-sdk-go-v2/config v1.33.21
	github.com/aws/aws-sdk-go-v2/credentials v1.21.3
	github.com/aws/aws-sdk-go-v2/service/ecr v1.69.0
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.47.12
	github.com/aws/aws-sdk-go-v2/service/signin v1.14.0
	github.com/aws/smithy-go v1.28.9
	github.com/docker/docker-credential-helpers v0.9.13
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.10.2
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.20.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.39.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.44.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.53.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.20.0 // indirect
)
