# Kubed

[Kubed by AppsCode](https://github.com/appscode/kubed) - A Kubernetes cluster manager daemon

## TL;DR;

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm install kubed appscode/kubed -n kube-system
```

## Introduction

This chart deploys a Kubed operator on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.11+

## Installing the Chart

To install the chart with the release name `kubed`:

```console
$ helm install kubed appscode/kubed -n kube-system
```

The command deploys a Kubed operator on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `kubed`:

```console
$ helm delete kubed -n kube-system
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the `kubed` chart and their default values.

|              Parameter               |                                                                                                            Description                                                                                                             |       Default       |
|--------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------|
| nameOverride                         | Overrides name template                                                                                                                                                                                                            | `""`                |
| fullnameOverride                     | Overrides fullname template                                                                                                                                                                                                        | `""`                |
| replicaCount                         | Number of Kubed operator replicas to create (only 1 is supported)                                                                                                                                                                  | `1`                 |
| operator.registry                    | Docker registry used to pull Kubed operator image                                                                                                                                                                                  | `appscode`          |
| operator.repository                  | Kubed operator container image                                                                                                                                                                                                     | `kubed`             |
| operator.tag                         | Kubed operator container image tag                                                                                                                                                                                                 | `v0.12.0`           |
| operator.resources                   | Compute Resources required by the operator container                                                                                                                                                                               | `{}`                |
| operator.securityContext             | Security options the operator container should run with                                                                                                                                                                            | `{}`                |
| imagePullSecrets                     | Specify an array of imagePullSecrets. Secrets must be manually created in the namespace. <br> Example: <br> `helm template charts/kubed \` <br> `--set imagePullSecrets[0].name=sec0 \` <br> `--set imagePullSecrets[1].name=sec1` | `[]`                |
| imagePullPolicy                      | Container image pull policy                                                                                                                                                                                                        | `IfNotPresent`      |
| criticalAddon                        | If true, installs Kubed operator as critical addon                                                                                                                                                                                 | `false`             |
| logLevel                             | Log level for operator                                                                                                                                                                                                             | `3`                 |
| annotations                          | Annotations applied to operator deployment                                                                                                                                                                                         | `{}`                |
| podAnnotations                       | Annotations passed to operator pod(s).                                                                                                                                                                                             | `{}`                |
| nodeSelector                         | Node labels for pod assignment                                                                                                                                                                                                     | `{}`                |
| tolerations                          | Tolerations for pod assignment                                                                                                                                                                                                     | `[]`                |
| affinity                             | Affinity rules for pod assignment                                                                                                                                                                                                  | `{}`                |
| podSecurityContext                   | Security options the operator pod should run with.                                                                                                                                                                                 | `{"fsGroup":65535}` |
| serviceAccount.create                | Specifies whether a service account should be created                                                                                                                                                                              | `true`              |
| serviceAccount.annotations           | Annotations to add to the service account                                                                                                                                                                                          | `{}`                |
| serviceAccount.name                  | The name of the service account to use. If not set and create is true, a name is generated using the fullname template                                                                                                             | `""`                |
| apiserver.securePort                 | Port used by Kubed server                                                                                                                                                                                                          | `"8443"`            |
| apiserver.useKubeapiserverFqdnForAks | If true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)                                                                                                             | `true`              |
| apiserver.healthcheck.enabled        | healthcheck configures the readiness and liveliness probes for the operator pod.                                                                                                                                                   | `false`             |
| apiserver.servingCerts.generate      | If true, generates on install/upgrade the certs that allow the kube-apiserver (and potentially ServiceMonitor) to authenticate operators pods. Otherwise specify certs in `apiserver.servingCerts.{caCrt, serverCrt, serverKey}`.  | `true`              |
| apiserver.servingCerts.caCrt         | CA certficate used by serving certificate of Kubed server.                                                                                                                                                                         | `""`                |
| apiserver.servingCerts.serverCrt     | Serving certficate used by Kubed server.                                                                                                                                                                                           | `""`                |
| apiserver.servingCerts.serverKey     | Private key for the serving certificate used by Kubed server.                                                                                                                                                                      | `""`                |
| enableAnalytics                      | If true, sends usage analytics                                                                                                                                                                                                     | `true`              |
| config.clusterName                   | Set cluster-name to something meaningful to you, say, prod, prod-us-east, qa, etc. so that you can distinguish notifications sent by kubed                                                                                         | `unicorn`           |
| config.configSourceNamespace         | If set, configmaps and secrets from only this namespace will be synced                                                                                                                                                             | `""`                |
| config.kubeconfigContent             | kubeconfig file content for configmap and secret syncer                                                                                                                                                                            | `""`                |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install kubed appscode/kubed -n kube-system --set replicaCount=1
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```console
$ helm install kubed appscode/kubed -n kube-system --values values.yaml
```
