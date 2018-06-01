---
title: Synchronize Configuration across Kubernetes Clusters
description: Synchronize Configuration across Kubernetes Clusters
menu:
  product_kubed_0.7.0:
    identifier: inter-cluster-syncer
    name: Across Clusters
    parent: config-syncer
    weight: 15
product_name: kubed
menu_name: product_kubed_0.7.0
section_menu_id: guides
---

> New to Kubed? Please start [here](/docs/concepts/README.md).

# Synchronize Configuration across Clusters

You can synchronize a ConfigMap or a Secret into different clusters using Kubed. For this you need to provide a `kube-config` file consisting cluster contexts and specify context names in comma separated format using __`kubed.appscode.com/sync-contexts`__ annotation. Kubed will create a copy of that ConfigMap/Secret in all clusters specified by the annotation. _For each cluster, it will sync into source namespace by default, but if namespace specified in the context (in the `kube-config` file), it will sync into that namespace._ Note that, Kubed will not create any namespace, it has to be created beforehand.

If the data in the source ConfigMap/Secret is updated, all the copies will be updated. Either delete the source ConfigMap/Secret or remove the annotation from the source ConfigMap/Secret to remove the copies.

If the list of contexts specified by the annotation is updated, Kubed will synchronize the ConfigMap/Secret accordingly, ie. it will create ConfigMap/Secret  in the clusters listed in new annotation (if not already exists) and delete ConfigMap/Secret from the clusters that were synced before but not listed in new annotation.

Note that, Kubed will error out if multiple contexts listed in annotation point same cluster. Also Kubed assumes that none of cluster contexts in `kube-config` file points the source cluster.

## Before You Begin

At first, you need to have a Kubernetes cluster and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). You also need a `kube-config` file consisting cluster contexts where you want to sync your ConfigMap/Secret.

## Deploy Kubed

To enable config syncer for different clusters, you need a cluster config like below.

```yaml
$ cat ./docs/examples/cluster-syncer/config.yaml

clusterName: minikube
enableConfigSyncer: true
kubeConfigFile: /srv/kubed/kubeConfigFile
```

| Key                  | Description                                                                                      |
|----------------------|--------------------------------------------------------------------------------------------------|
| `clusterName`        | `Optional`. Specifies the source cluster name used in label `kubed.appscode.com/origin.cluster`. |
| `enableConfigSyncer` | `Required`. If set to `true`, ConfigMap/Secret synchronization operation will be enabled.        |
| `kubeConfigFile`     | `Required`. Specifies the path of `kube-config` file.                                            |

Lets' consider following demo `kube-config` file:

```yaml
$ cat ./docs/examples/config-syncer/demo-kubeconfig.yaml

apiVersion: v1
kind: Config
clusters:
- name: cluster-1
  cluster:
    certificate-authority-data: ...
    server: https://1.2.3.4
- name: cluster-2
  cluster:
    certificate-authority-data: ...
    server: https://2.3.4.5
users:
- name: user-1
  user:
    client-certificate: ...
    client-key: ...
- name: user-2
  user:
    client-certificate: ...
    client-key: ...
contexts:
- name: context-1
  context:
    cluster: cluster-1
    user: user-1
- name: context-2
  context:
    cluster: cluster-2
    user: user-2
    namespace: demo-cluster-2
```

Now, create a Secret with the Kubed cluster config under `config.yaml` key. Also include required `kube-config` file under `kubeConfigFile` key. You can use separate secret for `kube-config`.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/cluster-syncer/config.yaml \
    --from-file=kubeConfigFile=./docs/examples/cluster-syncer/demo-kubeconfig.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  config.yaml: ZW5hYmxlQ29uZmlnU3luY2VyOiB0cnVlCg==
  kubeConfigFile: base64 endoded contents of kube-config file
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T10:25:33Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "25114"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: c207c236-71ec-11e7-a5ec-0800273df5f2
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/setup/install.md).  It will mount the secret inside operator pod in path `/srv/kubed`. So kubed cluster config file will be available in path `/srv/kubed/config.yaml` and `kube-config` file will be available in path `/srv/kubed/kubeConfigFile`.  Once the operator pod is running, go to the next section.

## Synchronize ConfigMap

At first, create a ConfigMap called `omni` in the `demo` namespace. This will be our source ConfigMap.

```console
$ kubectl create namespace demo
namespace "demo" created

$ kubectl apply -f ./docs/examples/config-syncer/demo.yaml
configmap "omni" created
```

Now, apply the `kubed.appscode.com/sync-contexts: "context-1,context-2"` annotation to ConfigMap `omni`.

```console
$ kubectl annotate configmap omni kubed.appscode.com/sync="context-1,context-2" -n demo
configmap "omni" annotated
```

It will create configmap "omni" in `cluster-1` and `cluster-2`. For `cluster-1` it will sync into source namespace `demo`  since no namespace specified in `context-1` and for `cluster-2` it will sync into `demo-cluster-2` namespace since namespace specified in `context-2`. Here we assume that those namespaces already exits in the respective clusters.

Other concepts like updating source configmap, removing annotation, origin annotation, origin labels, etc. are similar to the tutorial described [here](/docs/guides/config-syncer/intra-cluster.md).

## Next Steps
 - Need to keep some configuration synchronized across namespaces? Try [Kubed config syncer](/docs/guides/config-syncer/intra-cluster.md).
 - Learn how to use Kubed to protect your Kubernetes cluster from disasters [here](/docs/guides/disaster-recovery/).
 - Want to keep an eye on your cluster with automated notifications? Setup Kubed [event forwarder](/docs/guides/cluster-events/).
 - Out of disk space because of too much logs in Elasticsearch or metrics in InfluxDB? Configure [janitors](/docs/guides/janitors.md) to delete old data.
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
