---
title: Overview | Developer Guide
description: Developer Guide Overview
menu:
  product_kubed_{{ .version }}:
    identifier: developer-guide-readme
    name: Overview
    parent: developer-guide
    weight: 15
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: setup
---

## Development Guide
This document is intended to be the canonical source of truth for things like supported toolchain versions for building Config Syncer.
If you find a requirement that this doc does not capture, please submit an issue on github.

This document is intended to be relative to the branch in which it is found. It is guaranteed that requirements will change over time
for the development branch, but release branches of Config Syncer should not change.

### Build Config Syncer
Some of the Config Syncer development helper scripts rely on a fairly up-to-date GNU tools environment, so most recent Linux distros should
work just fine out-of-the-box.

#### Setup GO
Config Syncer is written in Google's GO programming language. Currently, Config Syncer is developed and tested on **go 1.8.3**. If you haven't set up a GO
development environment, please follow [these instructions](https://golang.org/doc/code.html) to install GO.

#### Download Source

```console
$ go get kubeops.dev/config-syncer
$ cd $(go env GOPATH)/src/kubeops.dev/config-syncer
```

#### Install Dev tools
To install various dev tools for Config Syncer, run the following command:
```console
$ ./hack/builddeps.sh
```

#### Build Binary
```console
$ ./hack/make.py
$ kubed version
```

#### Run Binary
```console
$ kubed run \
  --secure-port=8443 \
  --kubeconfig="$HOME/.kube/config" \
  --authorization-kubeconfig="$HOME/.kube/config" \
  --authentication-kubeconfig="$HOME/.kube/config" \
  --authentication-skip-lookup \
  --clusterconfig=./hack/deploy/config.yaml
```

#### Dependency management
Config Syncer uses [Glide](https://github.com/Masterminds/glide) to manage dependencies. Dependencies are already checked in the `vendor` folder.
If you want to update/add dependencies, run:
```console
$ glide slow
```

#### Build Docker images
To build and push your custom Docker image, follow the steps below. To release a new version of Config Syncer, please follow the [release guide](/docs/setup/developer-guide/release.md).

```console
# Build Docker image
$ ./hack/docker/setup.sh; ./hack/docker/setup.sh push

# Add docker tag for your repository
$ docker tag appscode/kubed:<tag> <image>:<tag>

# Push Image
$ docker push <image>:<tag>
```

#### Generate CLI Reference Docs
```console
$ ./hack/gendocs/make.sh
```

### Testing Config Syncer
#### Unit tests
```console
$ ./hack/make.py test unit
```

#### Run e2e tests
Config Syncer uses [Ginkgo](http://onsi.github.io/ginkgo/) to run e2e tests.
```console
$ ./hack/make.py test e2e
```

To run e2e tests against remote backends, you need to set cloud provider credentials in `./hack/config/.env`. You can see an example file in `./hack/config/.env.example`.
