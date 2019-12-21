# Kubed
[Kubed by AppsCode](https://github.com/appscode/kubed) - A Kubernetes cluster manager daemon.

## TL;DR;

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm install kubed appscode/kubed -n kube-system
```

## Introduction

This chart bootstraps a [Kubed controller](https://github.com/appscode/kubed) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.12+

## Installing the Chart
To install the chart with the release name `my-release`:

```console
$ helm install my-release appscode/kubed -n kube-system
```

The command deploys Kubed operator on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release`:

```console
$ helm delete my-release -n kube-system
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Kubed chart and their default values.


| Parameter                        | Description                                                       | Default            |
| ---------------------------------| ------------------------------------------------------------------|--------------------|
| `replicaCount`                   | Number of kubed operator replicas to create (only 1 is supported) | `1`                |
| `kubed.registry`                 | Docker registry used to pull Kubed image                          | `appscode`         |
| `kubed.repository`               | Kubed container image                                             | `kubed`            |
| `kubed.tag`                      | Kubed container image tag                                         | `v0.11.0`          |
| `imagePullSecrets`               | Specify image pull secrets                                        | `[]`               |
| `imagePullPolicy`                | Image pull policy                                                 | `IfNotPresent`     |
| `criticalAddon`                  | If true, installs kubed operator as critical addon                | `false`            |
| `logLevel`                       | Log level for kubed                                               | `3`                |
| `affinity`                       | Affinity rules for pod assignment                                 | `{}`               |
| `persistence.enabled`            | Use persistent volume to store data                               | `false`            |
| `persistence.size`               | Size of persistent volume claim                                   | `10Gi`             |
| `persistence.existingClaim`      | Use an existing PVC to persist data                               | `nil`              |
| `persistence.storageClassName`   | Type of persistent volume claim                                   | `nil`              |
| `persistence.accessModes`        | Persistence access modes                                          | `[ReadWriteOnce]`  |
| `persistence.subPath`            | Mount a sub dir of the persistent volume                          | `nil`              |
| `nodeSelector`                   | Node labels for pod assignment                                    | `{}`               |
| `tolerations`                    | Tolerations used pod assignment                                   | `{}`               |
| `resources`                      | Compute resources for the kubed container                         | `{}`               |
| `rbac.create`                    | If `true`, create and use RBAC resources                          | `true`             |
| `serviceAccount.create`          | If `true`, create a new service account                           | `true`             |
| `serviceAccount.name`            | Service account to be used. If not set and `serviceAccount.create` is `true`, a name is generated using the fullname template | `` |
| `apiserver.enabled`              | If `true`, enable kubed api server                                | `true`             |
| `apiserver.groupPriorityMinimum` | The minimum priority the group should have.                       | 10000              |
| `apiserver.versionPriority`      | The ordering of this API inside of the group.                     | 15                 |
| `apiserver.useKubeapiserverFqdnForAks` | If true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 | `true`             |
| `apiserver.healthcheck.enabled`  | Enable readiness and liveliness probes                            | `true`             |
| `enableAnalytics`                | Send usage events to Google Analytics                             | `true`             |
| `config.clusterName`             | Set cluster name to something meaningful to you, say, prod, prod-us-east, qa, etc. so that you can distinguish notifications sent by kubed | `unicorn`          |
| `config.enableConfigSyncer`      | If `true`, enables configmap and secret syncer                    | `true`             |
| `config.enableEventForwarder`    | If `true`, enables event forwarder                                | `false`            |
| `config.enableRecycleBin`        | If `true`, enables recycle bin for deleted objects                | `true`             |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install my-release appscode/kubed -n kube-system --set image.tag=v0.2.1
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```console
$ helm install my-release appscode/kubed -n kube-system --values values.yaml
```
