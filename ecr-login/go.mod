module github.com/awslabs/amazon-ecr-credential-helper/ecr-login

go 1.23

require (
	github.com/aws/aws-sdk-go-v2 v1.38.1
	github.com/aws/aws-sdk-go-v2/config v1.31.3
	github.com/aws/aws-sdk-go-v2/credentials v1.18.7
	github.com/aws/aws-sdk-go-v2/service/ecr v1.49.2
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.37.0
	github.com/aws/smithy-go v1.23.2
	github.com/docker/docker-credential-helpers v0.9.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.28.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.34.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.38.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
