FROM ubuntu:18.04

RUN DEBIAN_FRONTEND=noninteractive apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    golang-go \
    make \
    git \
    snapcraft

WORKDIR /build
CMD ["snapcraft"]
