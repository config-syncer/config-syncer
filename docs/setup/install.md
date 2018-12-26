---
title: Kubed Install
description: Kubed Install
menu:
  product_kubed_0.9.0:
    identifier: kubed-install
    name: Install
    parent: setup
    weight: 10
product_name: kubed
menu_name: product_kubed_0.9.0
section_menu_id: setup
---

> New to Kubed? Please start [here](/docs/concepts/README.md).

# Installation Guide

## Create Cluster Config
Before you can install Kubed, you need a cluster config for Kubed. Cluster config is defined in YAML format. You find an example config in [./hack/deploy/config.yaml](https://raw.githubusercontent.com/appscode/kubed/0.9.0/hack/deploy/config.yaml).

```yaml
$ cat https://raw.githubusercontent.com/appscode/kubed/0.9.0/hack/deploy/config.yaml

clusterName: unicorn
enableConfigSyncer: true
eventForwarder:
  receivers:
  - notifier: Mailgun
    to:
    - ops@example.com
  rules:
  - namespaces:
    - kube-system
    operations:
    - CREATE
    resources:
    - group: ""
      resources:
      - events
  - operations:
    - CREATE
    resources:
    - group: ""
      resources:
      - nodes
      - persistentvolumes
      - persistentvolumeclaims
    - group: storage.k8s.io
      resources:
      - storageclasses
    - group: extensions
      resources:
      - ingresses
    - group: voyager.appscode.com
      resources:
      - ingresses
    - group: certificates.k8s.io
      resources:
      - certificatesigningrequests
    - group: networking.k8s.io
      resources:
      - networkpolicies
notifierSecretName: notifier-config
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  ttl: 168h0m0s
```

To understand the various configuration options, check Kubed [tutorials](/docs/guides/README.md). Once you are satisfied with the configuration, create a Secret with the Kubed cluster config under `config.yaml` key.

```console
$ kubectl create secret generic kubed-config -n kube-system \
    --from-literal=config.yaml=$(curl -fsSL https://raw.githubusercontent.com/appscode/kubed/0.9.0/hack/deploy/config.yaml)
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

```

You may have to create another [Secret for notifiers](/docs/guides/cluster-events/notifiers.md), usually called `notifier-config`. If you are [storing cluster snapshots](/docs/guides/disaster-recovery/cluster-snapshot.md) in cloud storage, you have to create another Secret to provide cloud credentials.

### Generate Config using script
If you are familiar with GO, you can use the [./hack/config/main.go](https://github.com/appscode/kubed/blob/0.9.0/hack/config/main.go) script to generate a cluster config. Open this file in your favorite editor, update the config returned from `#CreateClusterConfig()` method. Then run the script to generate updated config in [./hack/deploy/config.yaml](https://raw.githubusercontent.com/appscode/kubed/0.9.0/hack/deploy/config.yaml).

```console
go run ./hack/config/main.go
```

### Verifying Cluster Config
Kubed includes a check command to verify a cluster config. Download the pre-built binary from [appscode/kubed Github releases](https://github.com/appscode/kubed/releases) and put the binary to some directory in your `PATH`.

```console
$ kubed check --clusterconfig=./hack/deploy/config.yaml
Cluster config was parsed successfully.
```

Kubed can be installed via a script or as a Helm chart.

<ul class="nav nav-tabs" id="installerTab" role="tablist">
  <li class="nav-item">
    <a class="nav-link active" id="script-tab" data-toggle="tab" href="#script" role="tab" aria-controls="script" aria-selected="true">Script</a>
  </li>
  <li class="nav-item">
    <a class="nav-link" id="helm-tab" data-toggle="tab" href="#helm" role="tab" aria-controls="helm" aria-selected="false">Helm</a>
  </li>
</ul>
<div class="tab-content" id="installerTabContent">
  <div class="tab-pane fade show active" id="script" role="tabpanel" aria-labelledby="script-tab">

## Using Script

Kubed can be installed via installer script included in the [/hack/deploy](https://github.com/appscode/kubed/tree/0.9.0/hack/deploy) folder.

```console
# set cluster-name to something meaningful to you, say, prod, prod-us-east, qa, etc.
# so that you can distinguish notifications sent by kubed
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/0.9.0/hack/deploy/kubed.sh \
    | bash -s -- --cluster-name=<your-cluster-name>
```

#### Customizing Installer

You can see the full list of flags available to installer using `-h` flag.

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/0.9.0/hack/deploy/kubed.sh | bash -s -- -h
kubed.sh - install Kubernetes cluster daemon

kubed.sh [options]

options:
-h, --help                             show brief help
-n, --namespace=NAMESPACE              specify namespace (default: kube-system)
    --rbac                             create RBAC roles and bindings (default: true)
    --docker-registry                  docker registry used to pull kubed images (default: appscode)
    --image-pull-secret                name of secret used to pull kubed operator images
    --run-on-master                    run kubed operator on master
    --cluster-name                     name of cluster (default: unicorn)
    --enable-apiserver                 enable/disable kubed apiserver
    --use-kubeapiserver-fqdn-for-aks   if true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)
    --enable-analytics                 send usage events to Google Analytics (default: true)
    --uninstall                        uninstall kubed
```

If you would like to run Kubed operator pod in `master` instances, pass the `--run-on-master` flag:

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/0.9.0/hack/deploy/kubed.sh \
    | bash -s -- --run-on-master [--rbac]
```

Kubed operator will be installed in a `kube-system` namespace by default. If you would like to run Kubed operator pod in `kubed` namespace, pass the `--namespace=kubed` flag:

```console
$ kubectl create namespace kubed
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/0.9.0/hack/deploy/kubed.sh \
    | bash -s -- --namespace=kubed [--run-on-master] [--rbac]
```

If you are using a private Docker registry, you need to pull the following docker image:

 - [appscode/kubed](https://hub.docker.com/r/appscode/kubed)

To pass the address of your private registry and optionally a image pull secret use flags `--docker-registry` and `--image-pull-secret` respectively.

```console
$ kubectl create namespace kubed
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/0.9.0/hack/deploy/kubed.sh \
    | bash -s -- --docker-registry=MY_REGISTRY [--image-pull-secret=SECRET_NAME] [--rbac]
```

</div>
<div class="tab-pane fade" id="helm" role="tabpanel" aria-labelledby="helm-tab">

## Using Helm
Kubed can be installed via [Helm](https://helm.sh/) using the [chart](https://github.com/appscode/kubed/tree/0.9.0/chart/kubed) from [AppsCode Charts Repository](https://github.com/appscode/charts). To install the chart with the release name `my-release`:

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm search appscode/kubed
NAME            CHART VERSION APP VERSION   DESCRIPTION
appscode/kubed  0.9.0    0.9.0    Kubed by AppsCode - Kubernetes daemon

# set cluster-name to something meaningful to you, say, prod, prod-us-east, qa, etc.
# so that you can distinguish notifications sent by kubed

# Kubernetes 1.8.x
$ helm install appscode/kubed --name kubed --version 0.9.0 \
  --namespace kube-system \
  --set config.clusterName=<your-cluster-name> \
  --set apiserver.enabled=false

# Kubernetes 1.9.0 or later
$ helm install appscode/kubed --name kubed --version 0.9.0 \
  --namespace kube-system \
  --set config.clusterName=<your-cluster-name>
```

To see the detailed configuration options, visit [here](https://github.com/appscode/kubed/tree/0.9.0/chart/kubed).

</div>

### Installing in GKE Cluster

If you are installing Kubed on a GKE cluster, you will need cluster admin permissions to install Kubed operator. Run the following command to grant admin permision to the cluster.

```console
# get current google identity
$ gcloud info | grep Account
Account: [user@example.org]

$ kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=user@example.org
```


## Verify installation
To check if Kubed operator pods have started, run the following command:

```console
$ kubectl get pods --all-namespaces -l app=kubed --watch
```

Once the operator pods are running, you can cancel the above command by typing `Ctrl+C`.


## Configuring RBAC
Kubed creates a custom resource: `SearchResult`. Kubed installer will create a user facing cluster role:

| ClusterRole           | Aggregates To     | Desription                            |
|-----------------------|-------------------|---------------------------------------|
| appscode:voyager:view | admin, edit, view | Allows read-only access to Kubed resources, intended to be granted within a namespace using a RoleBinding. |

These user facing roles supports [ClusterRole Aggregation](https://kubernetes.io/docs/admin/authorization/rbac/#aggregated-clusterroles) feature in Kubernetes 1.9 or later clusters.


## Update Cluster Config
If you would like to update cluster config, update the `kubed-config` Secret. Kubed will notice the change in config file and automatically apply the updated configuration.
