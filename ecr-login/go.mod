module github.com/awslabs/amazon-ecr-credential-helper/ecr-login

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.32.4
	github.com/aws/aws-sdk-go-v2/config v1.28.2
	github.com/aws/aws-sdk-go-v2/credentials v1.17.43
	github.com/aws/aws-sdk-go-v2/service/ecr v1.36.4
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.27.4
	github.com/aws/smithy-go v1.22.0
	github.com/docker/docker-credential-helpers v0.8.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.32.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
