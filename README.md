# Amazon ECR Docker Credential Helper

![Amazon ECR logo](docs/ecr.png "Amazon ECR")

[![Build](https://github.com/awslabs/amazon-ecr-credential-helper/actions/workflows/build.yaml/badge.svg)](https://github.com/awslabs/amazon-ecr-credential-helper/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/awslabs/amazon-ecr-credential-helper)](https://goreportcard.com/report/github.com/awslabs/amazon-ecr-credential-helper)
[![latest packaged version(s)](https://repology.org/badge/latest-versions/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)

The Amazon ECR Docker Credential Helper is a
[credential helper](https://github.com/docker/docker-credential-helpers)
for the Docker daemon that makes it easier to use
[Amazon Elastic Container Registry](https://aws.amazon.com/ecr/).

## Table of Contents
  * [Prerequisites](#prerequisites)
  * [Installing](#installing)
    + [Amazon Linux 2023 (AL2023)](#amazon-linux-2023-al2023)
    + [Amazon Linux 2 (AL2)](#amazon-linux-2-al2)
    + [Mac OS](#mac-os)
    + [Debian Buster (and future versions)](#debian-buster-and-future-versions)
    + [Ubuntu 19.04 Disco Dingo and newer](#ubuntu-1904-disco-dingo-and-newer)
    + [Arch Linux](#arch-linux)
    + [Alpine Linux](#alpine-linux)
    + [Windows](#windows)
    + [From Source](#from-source)
  * [Configuration](#configuration)
    + [Docker](#docker)
    + [AWS credentials](#aws-credentials)
    + [Amazon ECR Docker Credential Helper](#amazon-ecr-docker-credential-helper-1)
  * [Usage](#usage)
  * [Troubleshooting](#troubleshooting)
  * [Security disclosures](#security-disclosures)
  * [License](#license)

## Prerequisites

You must have at least Docker 1.11 installed on your system.

You also must have AWS credentials available.  See the [AWS credentials section](#aws-credentials) for details on how to
use different AWS credentials.

## Installing

### Amazon Linux 2023 (AL2023)
You can install the Amazon ECR Credential Helper from the Amazon Linux 2023 repositories.

```bash
$ sudo dnf install -y amazon-ecr-credential-helper
```

Once you have installed the credential helper, see the
[Configuration section](#configuration) for instructions on how to configure
Docker to work with the helper.

### Amazon Linux 2 (AL2)
You can install the Amazon ECR Credential Helper from the [`docker` or `ecs`
extras](https://docs.aws.amazon.com/linux/al2/ug/al2-extras.html).

```bash
$ sudo amazon-linux-extras enable docker
$ sudo yum install amazon-ecr-credential-helper
```

Once you have installed the credential helper, see the
[Configuration section](#configuration) for instructions on how to configure
Docker to work with the helper.

### Mac OS
A community-maintained Homebrew formula is available in the core tap.

[![Homebrew package](https://repology.org/badge/version-for-repo/homebrew/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)

```bash
$ brew install docker-credential-helper-ecr
```

On macOS, another community-maintained installation method is to use MacPorts.

[![MacPorts package](https://repology.org/badge/version-for-repo/macports/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)

```bash
$ sudo port install docker-credential-helper-ecr
```

Once you have installed the credential helper, see the
[Configuration section](#configuration) for instructions on how to configure
Docker to work with the helper.

### Debian Buster (and future versions)
You can install the Amazon ECR Credential Helper from the Debian Buster
archives.  This package will also be included in future releases of Debian.

[![Debian 10 package](https://repology.org/badge/version-for-repo/debian_10/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)
[![Debian 11 package](https://repology.org/badge/version-for-repo/debian_11/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)
[![Debian 12 package](https://repology.org/badge/version-for-repo/debian_12/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)
[![Debian Unstable package](https://repology.org/badge/version-for-repo/debian_unstable/amazon-ecr-credential-helper.svg)](https://repology.org/metapackage/amazon-ecr-credential-helper/versions)

```bash
$ sudo apt update
$ sudo apt install amazon-ecr-credential-helper
```

Once you have installed the credential helper, see the
[Configuration section](#configuration) for instructions on how to configure
Docker to work with the helper.

### Ubuntu 19.04 Disco Dingo and newer
You can install the Amazon ECR Credential Helper from the Ubuntu 19.04 Disco
Dingo (and newer) archives.

[![Ubuntu 20.04 package](https://repology.org/badge/version-for-repo/ubuntu_20_04/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)
[![Ubuntu 22.04 package](https://repology.org/badge/version-for-repo/ubuntu_22_04/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)

```bash
$ sudo apt update
$ sudo apt install amazon-ecr-credential-helper
```

Once you have installed the credential helper, see the
[Configuration section](#configuration) for instructions on how to configure
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
[Configuration section](#configuration) for instructions on how to configure
Docker to work with the helper.

### Alpine Linux
A community-maintained package is available in the [Alpine Linux aports Repository](https://pkgs.alpinelinux.org/packages?name=docker-credential-ecr-login).

[![Alpine Linux Edge package](https://repology.org/badge/version-for-repo/alpine_edge/amazon-ecr-credential-helper.svg)](https://repology.org/project/amazon-ecr-credential-helper/versions)

```bash
$ apk add docker-credential-ecr-login
```
> [!NOTE] 
> Badge only shows edge, check [repository](https://pkgs.alpinelinux.org/packages?name=docker-credential-ecr-login) for stable releases or add `--repository=http://dl-cdn.alpinelinux.org/alpine/edge/community`

Once you have installed the credential helper, see the
[Configuration section](#configuration) for instructions on how to configure
Docker to work with the helper.

### Windows
Windows executables are available via [GitHub releases](https://github.com/awslabs/amazon-ecr-credential-helper/releases).

> [!NOTE]
> Windows ARM support is considered [experimental](#experimental-features).
>
> See https://github.com/awslabs/amazon-ecr-credential-helper/issues/795

### From Source
To build and install the Amazon ECR Docker Credential Helper, we suggest Go
1.19 or later, `git` and `make` installed on your system.

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

You can install this via the `go` command line tool.

To install run:

```
go install github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cli/docker-credential-ecr-login@latest
```

> [!WARNING]
> Disclaimer: the [Dockerfile](./Dockerfile) in this repository is used to test cross-compilation of the
> Amazon ECR credential helper binaries in GitHub Actions CI and as a developer utility for building locally from source.
> It is a reference implementation and not security hardened for building and running production containers.

If you already have Docker environment, just clone this repository anywhere
and run `make build-in-docker`. This command builds the binary with Go inside the Docker
container and output it to local directory.

With `TARGET_GOOS` environment variable, you can also cross compile the binary.

Once you have installed the credential helper, see the
[Configuration section](#configuration) for instructions on how to configure
Docker to work with the helper.

## Configuration

### Docker

There is no need to use `docker login` or `docker logout`.

Place the `docker-credential-ecr-login` binary on your `PATH`.
On Windows, depending on whether the executable is ran in the User or System context, the corresponding `Path` user or system variable needs to be used.

Following that the configuration for the docker client needs to be updated in `~/.docker/config.json` to use the **ecr-login** helper.
Depending on the operating system and context under which docker client will be executed, this configuration can be found in different places.
  
On Linux systems:
- `/home/<username>/.docker/config.json` for **user** context
- `/root/.docker/config.json` for **root** context
  
On Windows:
- `C:\Users\<username>\.docker\config.json` for **user** context
- `C:\Windows\System32\config\systemprofile\.docker\config.json` for the **SYSTEM** context

Set the contents of the file to the following:

```json
{
	"credsStore": "ecr-login"
}
```
This configures the Docker daemon to use the credential helper for all Amazon
ECR registries.

With Docker 1.13.0 or greater, you can configure Docker to use different
credential helpers for different ECR registries. To use this credential helper for
a specific ECR registry, create a `credHelpers` section with the URI of your
ECR registry:

```json
{
	"credHelpers": {
		"public.ecr.aws": "ecr-login",
		"<aws_account_id>.dkr.ecr.<region>.amazonaws.com": "ecr-login"
	}
}
```

This is useful if you use `docker` to operate on registries that use different
authentication credentials.

If you need to authenticate with multiple registries, including non-ECR registries, you can combine credHelpers with auths. For example:
```json
{
  "credHelpers": {
    "<aws_account_id>.dkr.ecr.<region>.amazonaws.com": "ecr-login"
  },
  "auths": {
      "ghcr.io": {
        "auth": [GITHUB_PERSONAL_ACCESS_TOKEN]
      },
      "https://index.docker.io/v1/": {
        "auth": [docker.io-auth-token]
      },
      "registry.gitlab.com": {
        "auth": [gitlab-auth-token]
      }
	}
}
```

### AWS credentials

The Amazon ECR Docker Credential Helper allows you to use AWS credentials stored in different locations. Standard ones
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
[Configuration and Credential Files](http://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)
in the *AWS Command Line Interface User Guide*.

The credentials must have a policy applied that
[allows access to Amazon ECR](https://docs.aws.amazon.com/AmazonECR/latest/userguide/security-iam-awsmanpol.html).

### Amazon ECR Docker Credential Helper

| Environment Variable         | Sample Value  | Description                                                        |
| ---------------------------- | ------------- | ------------------------------------------------------------------ |
| AWS_ECR_DISABLE_CACHE        | true          | Disables the local file auth cache if set to a non-empty value     |
| AWS_ECR_CACHE_DIR            | ~/.ecr        | Specifies the local file auth cache directory location             |
| AWS_ECR_IGNORE_CREDS_STORAGE | true          | Ignore calls to docker login or logout and pretend they succeeded  |

## Usage

`docker pull 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

`docker push 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

`docker pull public.ecr.aws/amazonlinux/amazonlinux:latest`

If you have configured additional profiles for use with the AWS CLI, you can use
those profiles by specifying the `AWS_PROFILE` environment variable when invoking `docker`.
For example:

`AWS_PROFILE=myprofile docker pull 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repository:my-tag`

There is no need to use `docker login` or `docker logout`.

## Troubleshooting

If you have previously authenticated with an ECR repository by using the `docker login` command manually
then Docker may have stored an auth token which has since expired.
Docker will continue to attempt to use that cached auth token
instead of utilizing the credential helper. You must explicitly remove the previously cached expired
token using `docker logout 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repository`. After that
Docker will start utilizing the ECR credential helper to fetch fresh credentials, and you will no longer
need to use `docker login` or `docker logout`.

Logs from the Amazon ECR Docker Credential Helper are stored in `~/.ecr/log`.

For more information about Amazon ECR, see the the
[Amazon Elastic Container Registry User Guide](http://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html).

## Experimental features

Features marked as experimental are optionally made available to users to test and provide feedback.

If you test any experimental feaures, you can give feedback via the feature's tracking issue regarding:
* Your experience with the feature
* Issues or problems
* Suggested improvements

Experimental features are incomplete in design and implementation. Backwards incompatible
changes may be introduced at any time or support dropped entirely. Therefore experimental 
features are **not recommended** for use in production environments.

## Security disclosures

If you think you’ve found a potential security issue, please do not post it in the Issues.  Instead, please follow the instructions [here](https://aws.amazon.com/security/vulnerability-reporting/) or [email AWS security directly](mailto:aws-security@amazon.com).

## License

The Amazon ECR Docker Credential Helper is licensed under the Apache 2.0
License.
