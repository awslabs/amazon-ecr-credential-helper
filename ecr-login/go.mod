module github.com/awslabs/amazon-ecr-credential-helper/ecr-login

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.20.1
	github.com/aws/aws-sdk-go-v2/config v1.18.29
	github.com/aws/aws-sdk-go-v2/credentials v1.13.28
	github.com/aws/aws-sdk-go-v2/service/ecr v1.19.2
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.16.5
	github.com/aws/smithy-go v1.14.1
	github.com/docker/docker-credential-helpers v0.8.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.38 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.36 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.29 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.20.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
