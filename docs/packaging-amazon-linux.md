# Packaging for Amazon Linux

Packaging sources for Amazon Linux are contained in the
(`amazonlinux` branch)[https://github.com/awslabs/amazon-ecr-credential-helper/tree/amazonlinux].

This branch is updated every time a release is prepared for Amazon Linux.

To prepare a release, update the `amazon-ecr-credential-helper.spec` file in
the root of this repository with the new version number (or release number)
and with appropriate changes in the `changelog` section.

To build an RPM locally for testing, you can use the `rpm` target in the
Makefile (`make rpm`) or use the `docker-rpm` target, which will use an
`amazonlinux:2` Docker container as a build environment.  Both targets will
generate artifacts in this repository which can be cleaned by running
`make clean`.
