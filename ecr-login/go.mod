module github.com/awslabs/amazon-ecr-credential-helper/ecr-login

require (
	github.com/aws/aws-sdk-go-v2 v1.16.16
	github.com/aws/aws-sdk-go-v2/config v1.17.7
	github.com/aws/aws-sdk-go-v2/credentials v1.12.20
	github.com/aws/aws-sdk-go-v2/service/ecr v1.17.17
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.13.16
	github.com/aws/smithy-go v1.13.3
	github.com/docker/docker-credential-helpers v0.6.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.0
)

go 1.13
