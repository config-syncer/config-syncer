#!/bin/bash

set -x

GOPATH=$(go env GOPATH)
PACKAGE_NAME=github.com/appscode/kubed
REPO_ROOT="$GOPATH/src/$PACKAGE_NAME"
DOCKER_REPO_ROOT="/go/src/$PACKAGE_NAME"

pushd $REPO_ROOT

rm -rf client listers informers

# Generate defaults
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 defaulter-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/apis/kubed" \
    --input-dirs "$PACKAGE_NAME/apis/kubed/v1alpha1" \
    --extra-peer-dirs "$PACKAGE_NAME/apis/kubed" \
    --extra-peer-dirs "$PACKAGE_NAME/apis/kubed/v1alpha1" \
    --output-file-base "zz_generated.defaults"

# Generate deep copies
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 deepcopy-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/apis/kubed" \
    --input-dirs "$PACKAGE_NAME/apis/kubed/v1alpha1" \
    --output-file-base zz_generated.deepcopy

# Generate conversions
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 conversion-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/apis/kubed" \
    --input-dirs "$PACKAGE_NAME/apis/kubed/v1alpha1" \
    --output-file-base zz_generated.conversion

# Generate openapi
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 openapi-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/apis/kubed/v1alpha1" \
    --output-package "$PACKAGE_NAME/apis/kubed/v1alpha1"

# Generate the internal clientset (client/clientset_generated/internalclientset)
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 client-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-base "$PACKAGE_NAME/apis/" \
   --input "kubed/" \
   --clientset-path "$PACKAGE_NAME/client/" \
   --clientset-name internalclientset

# Generate the versioned clientset (client/clientset_generated/clientset)
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 client-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-base "$PACKAGE_NAME/apis/" \
   --input "kubed/v1alpha1" \
   --clientset-path "$PACKAGE_NAME/" \
   --clientset-name "client"

# generate lister
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 lister-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-dirs="$PACKAGE_NAME/apis/kubed" \
   --input-dirs="$PACKAGE_NAME/apis/kubed/v1alpha1" \
   --output-package "$PACKAGE_NAME/listers"

# generate informer
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.9 informer-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-dirs "$PACKAGE_NAME/apis/kubed/v1alpha1" \
   --versioned-clientset-package "$PACKAGE_NAME/client" \
   --listers-package "$PACKAGE_NAME/listers" \
   --output-package "$PACKAGE_NAME/informers"

popd
