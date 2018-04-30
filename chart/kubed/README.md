# Kubed
[Kubed by AppsCode](https://github.com/appscode/kubed) - A Kubernetes cluster manager daemon.

## TL;DR;

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm install appscode/kubed
```

## Introduction

This chart bootstraps a [Kubed controller](https://github.com/appscode/kubed) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.8+

## Installing the Chart
To install the chart with the release name `my-release`:

```console
$ helm install appscode/kubed --name my-release
```

The command deploys Kubed operator on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release`:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Kubed chart and their default values.


| Parameter                        | Description                                                       | Default            |
| ---------------------------------| ------------------------------------------------------------------|--------------------|
| `replicaCount`                   | Number of kubed operator replicas to create (only 1 is supported) | `1`                |
| `kubed.registry`                 | Docker registry used to pull Kubed image                          | `appscode`         |
| `kubed.repository`               | Kubed container image                                             | `kubed`            |
| `kubed.tag`                      | Kubed container image tag                                         | `0.6.0-rc.0`       |
| `imagePullSecrets`               | Specify image pull secrets                                        | `nil` (does not add image pull secrets to deployed pods) |
| `imagePullPolicy`                | Image pull policy                                                 | `IfNotPresent`     |
| `criticalAddon`                  | If true, installs kubed operator as critical addon                | `false`            |
| `logLevel`                       | Log level for kubed                                               | `3`                |
| `nodeSelector`                   | Node labels for pod assignment                                    | `{}`               |
| `rbac.create`                    | If `true`, create and use RBAC resources                          | `true`             |
| `serviceAccount.create`          | If `true`, create a new service account                           | `true`             |
| `serviceAccount.name`            | Service account to be used. If not set and `serviceAccount.create` is `true`, a name is generated using the fullname template | `` |
| `apiserver.groupPriorityMinimum` | The minimum priority the group should have.                       | 10000              |
| `apiserver.versionPriority`      | The ordering of this API inside of the group.                     | 15                 |
| `apiserver.ca`                   | CA certificate used by main Kubernetes api server                 | ``                 |
| `enableAnalytics`                | Send usage events to Google Analytics                             | `true`             |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install --name my-release --set image.tag=v0.2.1 appscode/kubed
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```console
$ helm install --name my-release --values values.yaml appscode/kubed
```

## RBAC
By default the chart will not install the recommended RBAC roles and rolebindings.

You need to have the flag `--authorization-mode=RBAC` on the api server. See the following document for how to enable [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/).

To determine if your cluster supports RBAC, run the the following command:

```console
$ kubectl api-versions | grep rbac
```

If the output contains "beta", you may install the chart with RBAC enabled (see below).

### Enable RBAC role/rolebinding creation

To enable the creation of RBAC resources (On clusters with RBAC). Do the following:

```console
$ helm install --name my-release appscode/kubed --set rbac.create=true
```
