# Amazon ECR Docker Credential Helper

![Amazon ECR logo](docs/ecr.png "Amazon ECR")

[![Build Status](https://travis-ci.org/awslabs/amazon-ecr-credential-helper.svg?branch=master)](https://travis-ci.org/awslabs/amazon-ecr-credential-helper)

The Amazon ECR Docker Credential Helper is a
[credential helper](https://github.com/docker/docker-credential-helpers)
for the Docker daemon that makes it easier to use
[Amazon Elastic Container Registry](https://aws.amazon.com/ecr/).

## Prerequisites

You must have at least Docker 1.11 installed on your system.

You also must have AWS credentials available.  See the [AWS credentials section](#aws-credentials) for details on how to
use different AWS credentials.

## Installing

### Amazon Linux 2
You can install the Amazon ECR Credential Helper from the [`docker` or `ecs`
extras](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/amazon-linux-ami-basics.html#extras-library).

```bash
$ sudo amazon-linux-extras enable docker
$ sudo yum install amazon-ecr-credential-helper
```

Once you have installed the credential helper, see the
[Configuration section](#Configuration) for instructions on how to configure
Docker to work with the helper.

### Mac OS
A community-maintained Homebrew formula is available in the core tap.

[![Homebrew package](https://repology.org/badge/version-for-repo/homebrew/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)

```bash
$ brew install docker-credential-helper-ecr
```

Once you have installed the credential helper, see the
[Configuration section](#Configuration) for instructions on how to configure
Docker to work with the helper.

### Debian Buster (and future versions)
You can install the Amazon ECR Credential Helper from the Debian Buster
archives.  This package will also be included in future releases of Debian.

[![Debian Stable package](https://repology.org/badge/version-for-repo/debian_stable/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)
[![Debian Testing package](https://repology.org/badge/version-for-repo/debian_testing/amazon-ecr-credential-helper.svg)](https://repology.org/metapackage/amazon-ecr-credential-helper/versions)
[![Debian Unstable package](https://repology.org/badge/version-for-repo/debian_unstable/amazon-ecr-credential-helper.svg)](https://repology.org/metapackage/amazon-ecr-credential-helper/versions)

```bash
$ sudo apt update
$ sudo apt install amazon-ecr-credential-helper
```

Once you have installed the credential helper, see the
[Configuration section](#Configuration) for instructions on how to configure
Docker to work with the helper.

### Ubuntu 19.04 Disco Dingo and newer
You can install the Amazon ECR Credential Helper from the Ubuntu 19.04 Disco
Dingo (and newer) archives.

[![Ubuntu 19.04 package](https://repology.org/badge/version-for-repo/ubuntu_19_04/amazon-ecr-credential-helper.svg)](https://repology.org/metapackage/amazon-ecr-credential-helper/versions)
[![Ubuntu 19.10 package](https://repology.org/badge/version-for-repo/ubuntu_19_10/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)
[![Ubuntu 20.04 package](https://repology.org/badge/version-for-repo/ubuntu_20_04/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)

```bash
$ sudo apt update
$ sudo apt install amazon-ecr-credential-helper
```

Once you have installed the credential helper, see the
[Configuration section](#Configuration) for instructions on how to configure
Docker to work with the helper.

### Arch Linux
A community-maintained package is available in the Arch User Repository.

[![AUR package](https://repology.org/badge/version-for-repo/aur/amazon-ecr-credential-helper.svg)](https://repology.org/metapackage/amazon-ecr-credential-helper/versions)

```bash
$ git clone https://aur.archlinux.org/amazon-ecr-credential-helper.git
$ cd amazon-ecr-credential-helper
$ makepkg -si
```

Once you have installed the credential helper, see the
[Configuration section](#Configuration) for instructions on how to configure
Docker to work with the helper.

### From Source
To build and install the Amazon ECR Docker Credential Helper, we suggest Go
1.12+, `git` and `make` installed on your system.

If you just installed Go, make sure you also have added it to your PATH or 
Environment Vars (Windows). For example:

```
$ export GOPATH=$HOME/go
$ export PATH=$PATH:$GOPATH/bin
```

Or in Windows:

```
setx GOPATH %USERPROFILE%\go
<your existing PATH definitions>;%USERPROFILE%\go\bin
```

If you haven't defined the PATH, the command below will fail silently, and
running `docker-credential-ecr-login` will output: `command not found`

You can install this via `go get` with:

```
go get -u github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cli/docker-credential-ecr-login
```


If you already have Docker environment, just clone this repository anywhere
and run `make docker`. This command builds the binary with Go inside the Docker
container and output it to local directory.

With `TARGET_GOOS` environment variable, you can also cross compile the binary.

Once you have installed the credential helper, see the
[Configuration section](#Configuration) for instructions on how to configure
Docker to work with the helper.

## Configuration

### Docker

Place the `docker-credential-ecr-login` binary on your `PATH` and set the
contents of your `~/.docker/config.json` file to be:

```json
{
	"credsStore": "ecr-login"
}
```

This configures the Docker daemon to use the credential helper for all Amazon
ECR registries.

With Docker 1.13.0 or greater, you can configure Docker to use different
credential helpers for different registries. To use this credential helper for
a specific ECR registry, create a `credHelpers` section with the URI of your
ECR registry:

```json
{
	"credHelpers": {
		"aws_account_id.dkr.ecr.region.amazonaws.com": "ecr-login"
	}
}
```

This is useful if you use `docker` to operate on registries that use different
authentication credentials.

### AWS credentials

The Amazon ECR Docker Credential Helper allows you to use AWS credentials stored in different locations.  Standard ones
include:

* The shared credentials file (`~/.aws/credentials`)
* The `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables
* An [IAM role for an Amazon ECS task](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html)
* An [IAM role for Amazon EC2](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html)

To use credentials associated with a different named profile in the shared credentials file (`~/.aws/credentials`), you
may set the `AWS_PROFILE` environment variable. 

The Amazon ECR Docker Credential Helper reads and supports some configuration options specified in the AWS
shared configuration file (`~/.aws/config`).  To disable these options, you must set the `AWS_SDK_LOAD_CONFIG` environment
variable to `false`.  The supported options include:

* Assumed roles specified with `role_arn` and `source_profile`
* External credential processes specified with `credential_process`
* Web Identities like [IAM Roles for Service Accounts in
  Kubernetes](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) (*Note: Kubernetes
  users using containers with a non-root user may encounter permission issues described in [this
  bug](https://github.com/kubernetes-sigs/external-dns/pull/1185) and may need to employ a workaround adjusting the
  Kubernetes `securityContext`.*)

The Amazon ECR Docker Credential Helper uses the same credentials as the AWS
CLI and the AWS SDKs. For more information about configuring AWS credentials,
see
[Configuration and Credential Files](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-config-files)
in the *AWS Command Line Interface User Guide*.

The credentials must have a policy applied that
[allows access to Amazon ECR](http://docs.aws.amazon.com/AmazonECR/latest/userguide/ecr_managed_policies.html).

## Usage

`docker pull 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

`docker push 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

If you have configured additional profiles for use with the AWS CLI, you can use
those profiles by specifying the `AWS_PROFILE` environment variable when invoking `docker`.
For example:

`AWS_PROFILE=myprofile docker pull 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

There is no need to use `docker login` or `docker logout`.

## Troubleshooting

Logs from the Amazon ECR Docker Credential Helper are stored in `~/.ecr/log`.

For more information about Amazon ECR, see the the
[Amazon Elastic Container Registry User Guide](http://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html).

## Security disclosures

If you think you’ve found a potential security issue, please do not post it in the Issues.  Instead, please follow the instructions [here](https://aws.amazon.com/security/vulnerability-reporting/) or [email AWS security directly](mailto:aws-security@amazon.com).

## License

The Amazon ECR Docker Credential Helper is licensed under the Apache 2.0
License.
