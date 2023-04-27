# 0.7.0

* Feature - Allow callers to set log output. ([#309](https://github.com/awslabs/amazon-ecr-credential-helper/pull/309) and [#312](https://github.com/awslabs/amazon-ecr-credential-helper/pull/312))
* Upgrade dependencies for bug fixes and security patches. 

# 0.6.0

* Feature - Added support for [AWS SSO](https://aws.amazon.com/single-sign-on/) ([#229](https://github.com/awslabs/amazon-ecr-credential-helper/issues/229))
* Feature - Added support to assume roles via EC2 instance metadata. ([#282](https://github.com/awslabs/amazon-ecr-credential-helper/issues/282))
* Feature - Added support for [Apple Silicon](https://go.dev/doc/go1.16#darwin) ([#291](https://github.com/awslabs/amazon-ecr-credential-helper/pull/291))
* Enhancement - The AWS shared config file (`~/.aws/config`) is now always enabled. (`AWS_SDK_LOAD_CONFIG` environment variable is no longer supported) ([#282](https://github.com/awslabs/amazon-ecr-credential-helper/issues/282))

# 0.5.0

* Feature - Added support for [ECR Public](https://gallery.ecr.aws) ([#253](https://github.com/awslabs/amazon-ecr-credential-helper/pull/253))
* Feature - Added support for [EC2 IMDSv2](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-instance-metadata-service.html) ([#240](https://github.com/awslabs/amazon-ecr-credential-helper/pull/240))
* Enhancement - The AWS shared config file (`~/.aws/config`) is now enabled by default.  This can be disabled by setting the environment variable `AWS_SDK_LOAD_CONFIG` to `false` ([#201](https://github.com/awslabs/amazon-ecr-credential-helper/pull/201))
* Bug - Fixed an issue where output from a `credential_process` was sometimes too big for the default buffer ([#240](https://github.com/awslabs/amazon-ecr-credential-helper/pull/240))

# 0.4.0

* Feature - Added support for chaining assumed roles in the shared config file (`~/.aws/config`) defined by `source_profile` and `credential_source` ([#177](https://github.com/awslabs/amazon-ecr-credential-helper/pull/177))
* Feature - Added support for Web Identities and [IAM Roles for Service Accounts (IRSA)](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) with Kubernetes ([#183](https://github.com/awslabs/amazon-ecr-credential-helper/pull/183))
* Bug - Fixed the `make docker` target when the credential helper git repository is used as a git submodule ([#184](https://github.com/awslabs/amazon-ecr-credential-helper/issues/184))

# 0.3.1

* Bug - Log directory is now automatically created when the helper runs

# 0.3.0

* Feature - Added support for PrivateLink endpoints

# 0.2.0

* Feature - Added support for FIPS endpoints

# 0.1.0

Initial release

