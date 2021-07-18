---
title: Kubed Install
description: Kubed Install
menu:
  product_kubed_{{ .version }}:
    identifier: kubed-install
    name: Install
    parent: setup
    weight: 10
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: setup
---

> New to Kubed? Please start [here](/docs/concepts/README.md).

# Installation Guide

Kubed operator can be installed via a script or as a Helm chart.

<ul class="nav nav-tabs" id="installerTab" role="tablist">
  <li class="nav-item">
    <a class="nav-link active" id="helm3-tab" data-toggle="tab" href="#helm3" role="tab" aria-controls="helm3" aria-selected="true">Helm 3 (Recommended)</a>
  </li>
  <li class="nav-item">
    <a class="nav-link" id="helm2-tab" data-toggle="tab" href="#helm2" role="tab" aria-controls="helm2" aria-selected="false">Helm 2</a>
  </li>
  <li class="nav-item">
    <a class="nav-link" id="script-tab" data-toggle="tab" href="#script" role="tab" aria-controls="script" aria-selected="false">YAML</a>
  </li>
</ul>
<div class="tab-content" id="installerTabContent">
  <div class="tab-pane fade show active" id="helm3" role="tabpanel" aria-labelledby="helm3-tab">

## Using Helm 3

Kubed can be installed via [Helm](https://helm.sh/) using the [chart](https://github.com/kubeops/kubed/tree/{{< param "info.version" >}}/charts/kubed) from [AppsCode Charts Repository](https://github.com/appscode/charts). To install the chart with the release name `my-release`:

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm search repo appscode/kubed --version {{< param "info.version" >}}
NAME            CHART VERSION APP VERSION DESCRIPTION
appscode/kubed  {{< param "info.version" >}}    {{< param "info.version" >}}  Kubed by AppsCode - Kubernetes daemon

$ helm install kubed appscode/kubed \
  --version {{< param "info.version" >}} \
  --namespace kube-system
```

To see the detailed configuration options, visit [here](https://github.com/kubeops/kubed/tree/{{< param "info.version" >}}/charts/kubed).

</div>
<div class="tab-pane fade" id="helm2" role="tabpanel" aria-labelledby="helm2-tab">

## Using Helm 2

Kubed can be installed via [Helm](https://helm.sh/) using the [chart](https://github.com/kubeops/kubed/tree/{{< param "info.version" >}}/charts/kubed) from [AppsCode Charts Repository](https://github.com/appscode/charts). To install the chart with the release name `my-release`:

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm search appscode/kubed --version {{< param "info.version" >}}
NAME            CHART VERSION APP VERSION DESCRIPTION
appscode/kubed  {{< param "info.version" >}}    {{< param "info.version" >}}  Kubed by AppsCode - Kubernetes daemon

$ helm install appscode/kubed --name kubed \
  --version {{< param "info.version" >}} \
  --namespace kube-system
```

To see the detailed configuration options, visit [here](https://github.com/kubeops/kubed/tree/{{< param "info.version" >}}/charts/kubed).

</div>
<div class="tab-pane fade" id="script" role="tabpanel" aria-labelledby="script-tab">

## Using YAML

If you prefer to not use Helm, you can generate YAMLs from Kubed chart and deploy using `kubectl`. Here we are going to show the prodecure using Helm 3.

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm search repo appscode/kubed --version {{< param "info.version" >}}
NAME            CHART VERSION APP VERSION DESCRIPTION
appscode/kubed  {{< param "info.version" >}}    {{< param "info.version" >}}  Kubed by AppsCode - Kubernetes daemon

$ helm template kubed appscode/kubed \
  --version {{< param "info.version" >}} \
  --namespace kube-system \
  --no-hooks | kubectl apply -f -
```

To see the detailed configuration options, visit [here](https://github.com/kubeops/kubed/tree/{{< param "info.version" >}}/charts/kubed).

</div>
</div>

### Installing in GKE Cluster

If you are installing Kubed on a GKE cluster, you will need cluster admin permissions to install Kubed operator. Run the following command to grant admin permision to the cluster.

```console
$ kubectl create clusterrolebinding "cluster-admin-$(whoami)" \
  --clusterrole=cluster-admin \
  --user="$(gcloud config get-value core/account)"
```

In addition, if your GKE cluster is a [private cluster](https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters), you will need to either add an additional firewall rule that allows master nodes access port `8443/tcp` on worker nodes, or change the existing rule that allows access to ports `443/tcp` and `10250/tcp` to also allow access to port `8443/tcp`. The procedure to add or modify firewall rules is described in the official GKE documentation for private clusters mentioned before.

## Verify installation

Kubed includes a check command to verify a cluster config. Download the pre-built binary from [appscode/kubed Github releases](https://github.com/kubeops/kubed/releases) and put the binary to some directory in your `PATH`.

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

Kubed can be installed via installer script included in the [/hack/deploy](https://github.com/kubeops/kubed/tree/{{< param "info.version" >}}/hack/deploy) folder.

```console
# set cluster-name to something meaningful to you, say, prod, prod-us-east, qa, etc.
# so that you can distinguish notifications sent by kubed
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/{{< param "info.version" >}}/hack/deploy/kubed.sh \
    | bash -s -- --cluster-name=<your-cluster-name>
```

#### Customizing Installer

You can see the full list of flags available to installer using `-h` flag.

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/{{< param "info.version" >}}/hack/deploy/kubed.sh | bash -s -- -h
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
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/{{< param "info.version" >}}/hack/deploy/kubed.sh \
    | bash -s -- --run-on-master [--rbac]
```

Kubed operator will be installed in a `kube-system` namespace by default. If you would like to run Kubed operator pod in `kubed` namespace, pass the `--namespace=kubed` flag:

```console
$ kubectl create namespace kubed
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/{{< param "info.version" >}}/hack/deploy/kubed.sh \
    | bash -s -- --namespace=kubed [--run-on-master] [--rbac]
```

If you are using a private Docker registry, you need to pull the following docker image:

 - [appscode/kubed](https://hub.docker.com/r/appscode/kubed)

To pass the address of your private registry and optionally a image pull secret use flags `--docker-registry` and `--image-pull-secret` respectively.

```console
$ kubectl create namespace kubed
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/{{< param "info.version" >}}/hack/deploy/kubed.sh \
    | bash -s -- --docker-registry=MY_REGISTRY [--image-pull-secret=SECRET_NAME] [--rbac]
```

</div>
<div class="tab-pane fade" id="helm" role="tabpanel" aria-labelledby="helm-tab">

## Using Helm
Kubed can be installed via [Helm](https://helm.sh/) using the [chart](https://github.com/kubeops/kubed/tree/{{< param "info.version" >}}/chart/kubed) from [AppsCode Charts Repository](https://github.com/appscode/charts). To install the chart with the release name `my-release`:

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm search appscode/kubed
NAME            CHART VERSION APP VERSION   DESCRIPTION
appscode/kubed  {{< param "info.version" >}}    {{< param "info.version" >}}    Kubed by AppsCode - Kubernetes daemon

# set cluster-name to something meaningful to you, say, prod, prod-us-east, qa, etc.
# so that you can distinguish notifications sent by kubed

# Kubernetes 1.8.x
$ helm install appscode/kubed --name kubed --version {{< param "info.version" >}} \
  --namespace kube-system \
  --set config.clusterName=<your-cluster-name> \
  --set apiserver.enabled=false

# Kubernetes 1.9.0 or later
$ helm install appscode/kubed --name kubed --version {{< param "info.version" >}} \
  --namespace kube-system \
  --set config.clusterName=<your-cluster-name>
```

To see the detailed configuration options, visit [here](https://github.com/kubeops/kubed/tree/{{< param "info.version" >}}/chart/kubed).

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
