#!/bin/bash

# Normalize to working directory being build root (up one level from ./snap)
SOURCES=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )

# docker build needs context, but we don't really want to provide any
TMPCONTEXT=$(mktemp -d)
docker build \
       --network=host \
       -t amazon/amazon-ecr-credential-helper:snapcraft \
       -f scripts/snapcraft.dockerfile \
       "${TMPCONTEXT}"

rmdir "${TMPCONTEXT}"

docker run \
       --rm -it\
       --net=host \
       -v ${SOURCES}:/build \
       amazon/amazon-ecr-credential-helper:snapcraft \
       snapcraft \
       $@
