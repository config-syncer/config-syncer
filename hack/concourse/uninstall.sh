#!/usr/bin/env bash

set -x

# uninstall operator
./hack/deploy/kubed.sh --uninstall --purge

# remove docker images
source "hack/libbuild/common/lib.sh"
detect_tag ''

# delete docker image on exit
./hack/libbuild/docker.py del_tag $DOCKER_REGISTRY kubed $TAG
