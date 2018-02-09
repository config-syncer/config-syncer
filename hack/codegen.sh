#!/bin/bash

set -x

GOPATH=$(go env GOPATH)
PACKAGE_NAME=github.com/appscode/kubed
REPO_ROOT="$GOPATH/src/$PACKAGE_NAME"
DOCKER_REPO_ROOT="/go/src/$PACKAGE_NAME"

pushd $REPO_ROOT

# Generate defaults
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 defaulter-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/pkg/apis/kubed" \
    --input-dirs "$PACKAGE_NAME/pkg/apis/kubed/v1alpha1" \
    --extra-peer-dirs "$PACKAGE_NAME/pkg/apis/kubed" \
    --extra-peer-dirs "$PACKAGE_NAME/pkg/apis/kubed/v1alpha1" \
    --output-file-base "zz_generated.defaults"

# Generate deep copies
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 deepcopy-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/pkg/apis/kubed" \
    --input-dirs "$PACKAGE_NAME/pkg/apis/kubed/v1alpha1" \
    --output-file-base zz_generated.deepcopy

# Generate conversions
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 conversion-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/pkg/apis/kubed" \
    --input-dirs "$PACKAGE_NAME/pkg/apis/kubed/v1alpha1" \
    --output-file-base zz_generated.conversion

# Generate openapi
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 openapi-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/pkg/apis/kubed/v1alpha1" \
    --output-package "$PACKAGE_NAME/pkg/apis/kubed/v1alpha1"

# Generate the internal clientset (client/clientset_generated/internalversion)
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 client-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-base "$PACKAGE_NAME/pkg/apis/" \
   --input "kubed/" \
   --clientset-path "$PACKAGE_NAME/pkg/client/clientset" \
   --clientset-name "internalversion"

# Generate the versioned clientset (client/clientset_generated/clientset)
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 client-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-base "$PACKAGE_NAME/pkg/apis/" \
   --input "kubed/v1alpha1/" \
   --clientset-path "$PACKAGE_NAME/pkg/client/" \
   --clientset-name "clientset"

# generate lister
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 lister-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-dirs="$PACKAGE_NAME/pkg/apis/kubed" \
   --input-dirs="$PACKAGE_NAME/pkg/apis/kubed/v1alpha1" \
   --output-package "$PACKAGE_NAME/pkg/client/listers"

popd
