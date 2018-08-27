#!/usr/bin/env bash

set -eoux pipefail

ORG_NAME=appscode
REPO_NAME=kubed
OPERATOR_NAME=kubed
APP_LABEL=kubed #required for `kubectl describe deploy -n kube-system -l app=$APP_LABEL`

export APPSCODE_ENV=test-concourse
export DOCKER_REGISTRY=appscodeci

# get concourse-common
pushd $REPO_NAME
git status # required, otherwise you'll get error `Working tree has modifications.  Cannot add.`. why?
git subtree pull --prefix hack/libbuild https://github.com/appscodelabs/libbuild.git master --squash -m 'concourse'
popd

source $REPO_NAME/hack/libbuild/concourse/init.sh

pushd $GOPATH/src/github.com/$ORG_NAME/$REPO_NAME

# install dependencies
./hack/builddeps.sh
./hack/docker/setup.sh build
./hack/docker/setup.sh push

./hack/deploy/kubed.sh --docker-registry=$DOCKER_REGISTRY
./hack/make.py test e2e --v=3 --kubeconfig=/root/.kube/config --selfhosted-operator=true
popd
