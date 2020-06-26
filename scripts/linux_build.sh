#!/bin/sh

# This must be run in a golang alpine container

apk add --no-cache make gcc musl-dev linux-headers
cd /terraform-provider-quorum
make distlocal ldflags='-linkmode external -extldflags "-static"'