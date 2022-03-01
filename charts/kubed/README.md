# Config Syncer

[Config Syncer by AppsCode](https://github.com/kubeops/config-syncer) - A Kubernetes cluster manager daemon

## TL;DR;

```bash
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm search repo appscode/kubed --version=v0.13.2
$ helm upgrade -i kubed appscode/kubed -n kube-system --create-namespace --version=v0.13.2
```

## Introduction

This chart deploys a Config Syncer operator on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.11+

## Installing the Chart

To install/upgrade the chart with the release name `kubed`:

```bash
$ helm upgrade -i kubed appscode/kubed -n kube-system --create-namespace --version=v0.13.2
```

The command deploys a Config Syncer operator on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall the `kubed`:

```bash
$ helm uninstall kubed -n kube-system
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the `kubed` chart and their default values.

|              Parameter               |                                                                                                            Description                                                                                                             |            Default             |
|--------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------|
| nameOverride                         | Overrides name template                                                                                                                                                                                                            | <code>""</code>                |
| fullnameOverride                     | Overrides fullname template                                                                                                                                                                                                        | <code>""</code>                |
| replicaCount                         | Number of Config Syncer operator replicas to create (only 1 is supported)                                                                                                                                                          | <code>1</code>                 |
| operator.registry                    | Docker registry used to pull Config Syncer operator image                                                                                                                                                                          | <code>appscode</code>          |
| operator.repository                  | Config Syncer operator container image                                                                                                                                                                                             | <code>kubed</code>             |
| operator.tag                         | Config Syncer operator container image tag                                                                                                                                                                                         | <code>v0.13.2</code>           |
| operator.resources                   | Compute Resources required by the operator container                                                                                                                                                                               | <code>{}</code>                |
| operator.securityContext             | Security options the operator container should run with                                                                                                                                                                            | <code>{}</code>                |
| imagePullSecrets                     | Specify an array of imagePullSecrets. Secrets must be manually created in the namespace. <br> Example: <br> `helm template charts/kubed \` <br> `--set imagePullSecrets[0].name=sec0 \` <br> `--set imagePullSecrets[1].name=sec1` | <code>[]</code>                |
| imagePullPolicy                      | Container image pull policy                                                                                                                                                                                                        | <code>IfNotPresent</code>      |
| criticalAddon                        | If true, installs Config Syncer operator as critical addon                                                                                                                                                                         | <code>false</code>             |
| logLevel                             | Log level for operator                                                                                                                                                                                                             | <code>3</code>                 |
| annotations                          | Annotations applied to operator deployment                                                                                                                                                                                         | <code>{}</code>                |
| podAnnotations                       | Annotations passed to operator pod(s).                                                                                                                                                                                             | <code>{}</code>                |
| nodeSelector                         | Node labels for pod assignment                                                                                                                                                                                                     | <code>{}</code>                |
| tolerations                          | Tolerations for pod assignment                                                                                                                                                                                                     | <code>[]</code>                |
| affinity                             | Affinity rules for pod assignment                                                                                                                                                                                                  | <code>{}</code>                |
| podSecurityContext                   | Security options the operator pod should run with.                                                                                                                                                                                 | <code>{"fsGroup":65535}</code> |
| serviceAccount.create                | Specifies whether a service account should be created                                                                                                                                                                              | <code>true</code>              |
| serviceAccount.annotations           | Annotations to add to the service account                                                                                                                                                                                          | <code>{}</code>                |
| serviceAccount.name                  | The name of the service account to use. If not set and create is true, a name is generated using the fullname template                                                                                                             | <code>""</code>                |
| apiserver.securePort                 | Port used by Config Syncer server                                                                                                                                                                                                  | <code>"8443"</code>            |
| apiserver.useKubeapiserverFqdnForAks | If true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)                                                                                                             | <code>true</code>              |
| apiserver.healthcheck.enabled        | healthcheck configures the readiness and liveliness probes for the operator pod.                                                                                                                                                   | <code>false</code>             |
| apiserver.servingCerts.generate      | If true, generates on install/upgrade the certs that allow the kube-apiserver (and potentially ServiceMonitor) to authenticate operators pods. Otherwise specify certs in `apiserver.servingCerts.{caCrt, serverCrt, serverKey}`.  | <code>true</code>              |
| apiserver.servingCerts.caCrt         | CA certficate used by serving certificate of Config Syncer server.                                                                                                                                                                 | <code>""</code>                |
| apiserver.servingCerts.serverCrt     | Serving certficate used by Config Syncer server.                                                                                                                                                                                   | <code>""</code>                |
| apiserver.servingCerts.serverKey     | Private key for the serving certificate used by Config Syncer server.                                                                                                                                                              | <code>""</code>                |
| enableAnalytics                      | If true, sends usage analytics                                                                                                                                                                                                     | <code>true</code>              |
| config.clusterName                   | Set cluster-name to something meaningful to you, say, prod, prod-us-east, qa, etc. so that you can distinguish notifications sent by kubed                                                                                         | <code>unicorn</code>           |
| config.configSourceNamespace         | If set, configmaps and secrets from only this namespace will be synced                                                                                                                                                             | <code>""</code>                |
| config.kubeconfigContent             | kubeconfig file content for configmap and secret syncer                                                                                                                                                                            | <code>""</code>                |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm upgrade -i`. For example:

```bash
$ helm upgrade -i kubed appscode/kubed -n kube-system --create-namespace --version=v0.13.2 --set replicaCount=1
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```bash
$ helm upgrade -i kubed appscode/kubed -n kube-system --create-namespace --version=v0.13.2 --values values.yaml
```
