# Amazon ECR Docker Credential Helper

[![Build Status](https://travis-ci.org/awslabs/amazon-ecr-credential-helper.svg?branch=master)](https://travis-ci.org/awslabs/amazon-ecr-credential-helper)

The Amazon ECR Docker Credential Helper is a
[credential helper](https://github.com/docker/docker-credential-helpers)
for the Docker daemon that makes it easier to use
[Amazon EC2 Container Registry](https://aws.amazon.com/ecr/).

## Prerequisites

You must have at least Docker 1.11 installed on your system.

You also must have AWS credentials available in one of the standard locations:

* The `~/.aws/credentials` file
* The `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables
* An [IAM role for Amazon EC2](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html)
* If you are working with an assumed role please set the environment variable: `AWS_SDK_LOAD_CONFIG=true` also.

The Amazon ECR Docker Credential Helper uses the same credentials as the AWS
CLI and the AWS SDKs. For more information about configuring AWS credentials,
see
[Configuration and Credential Files](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-config-files)
in the *AWS Command Line Interface User Guide*.

The credentials must have a policy applied that
[allows access to Amazon ECR](http://docs.aws.amazon.com/AmazonECR/latest/userguide/ecr_managed_policies.html).

## Building and Setup
To build the Amazon ECR Docker Credential Helper, you must have Go 1.5 or greater, and you must have `git`
and `make` installed on your system.

Clone this repository into your existing `GOPATH` under
`src/github.com/awslabs/amazon-ecr-credential-helper`, then run `make`.  The
resulting binary can be found in `bin/local/docker-credential-ecr-login`.

Or if you already have Docker environment, just clone this repository anywhere
and run `make docker`. This command builds the binary by Go inside docker container and
output it to local directory. With `TARGET_GOOS` environment variable, you can also
cross compile the binary.

Place the `docker-credential-ecr-login` binary on your `PATH` and set the contents
of your `~/.docker/config.json` file to be:

```json
{
	"credsStore": "ecr-login"
}
```

## Usage

`docker push 123457689012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

`docker push 123457689012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

There is no need to use `docker login` or `docker logout`.

## Troubleshooting

Logs from the Amazon ECR Docker Credential Helper are stored in `~/.ecr/log`.

For more information about Amazon ECR, see the the
[Amazon EC2 Container Registry User Guide](http://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html).

## License

The Amazon ECR Docker Credential Helper is licensed under the Apache 2.0
License.
