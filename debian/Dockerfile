FROM debian:sid

RUN DEBIAN_FRONTEND=noninteractive apt-get -o Acquire::ForceIPv4=true update && \
    DEBIAN_FRONTEND=noninteractive apt-get -o Acquire::ForceIPv4=true install -yq \
      debhelper \
      dh-golang \
      golang-go \
      git \
      git-buildpackage \
      golang-github-aws-aws-sdk-go-dev \
      golang-github-docker-docker-credential-helpers-dev \
      golang-github-mitchellh-go-homedir-dev \
      golang-github-sirupsen-logrus-dev \
      golang-github-golang-mock-dev \
      golang-github-pkg-errors-dev

WORKDIR /sources/amazon-ecr-credential-helper

CMD ["sh", "-c", "gbp buildpackage --git-tarball-dir=/sources --git-ignore-new --git-verbose -uc -us"]
