module github.com/awslabs/amazon-ecr-credential-helper/ecr-login

require (
	github.com/aws/aws-sdk-go-v2 v1.16.7
	github.com/aws/aws-sdk-go-v2/config v1.5.0
	github.com/aws/aws-sdk-go-v2/credentials v1.12.8
	github.com/aws/aws-sdk-go-v2/service/ecr v1.17.8
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.13.8
	github.com/aws/smithy-go v1.12.0
	github.com/docker/docker-credential-helpers v0.6.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.8.0
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
)

go 1.13
