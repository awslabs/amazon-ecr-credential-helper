module github.com/awslabs/amazon-ecr-credential-helper/ecr-login

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.33.0
	github.com/aws/aws-sdk-go-v2/config v1.29.0
	github.com/aws/aws-sdk-go-v2/credentials v1.17.53
	github.com/aws/aws-sdk-go-v2/service/ecr v1.38.5
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.29.3
	github.com/aws/smithy-go v1.22.1
	github.com/docker/docker-credential-helpers v0.8.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.8 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
