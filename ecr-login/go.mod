module github.com/awslabs/amazon-ecr-credential-helper/ecr-login

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.8
	github.com/aws/aws-sdk-go-v2/config v1.32.31
	github.com/aws/aws-sdk-go-v2/credentials v1.19.30
	github.com/aws/aws-sdk-go-v2/service/ecr v1.60.0
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.41.0
	github.com/aws/aws-sdk-go-v2/service/signin v1.13.0
	github.com/aws/smithy-go v1.28.8
	github.com/docker/docker-credential-helpers v0.9.12
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.9.4
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.33.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.45.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.20.0 // indirect
)
