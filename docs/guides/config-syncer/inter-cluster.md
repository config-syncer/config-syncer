---
title: Synchronize Configuration across Kubernetes Clusters
description: Synchronize Configuration across Kubernetes Clusters
menu:
  product_kubed_{{ .version }}:
    identifier: inter-cluster-syncer
    name: Across Clusters
    parent: config-syncer
    weight: 15
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: guides
---

> New to Config Syncer? Please start [here](/docs/concepts/README.md).

# Synchronize Configuration across Clusters

You can synchronize a ConfigMap or a Secret into different clusters using Config Syncer. For this you need to provide a `kubeconfig` file consisting cluster contexts and specify context names in comma separated format using __`kubed.appscode.com/sync-contexts`__ annotation. Config Syncer will create a copy of that ConfigMap/Secret in all clusters specified by the annotation. _For each cluster, it will sync into source namespace by default, but if namespace specified in the context (in the `kubeconfig` file), it will sync into that namespace._ Note that, Config Syncer will not create any namespace, it has to be created beforehand.

If the data in the source ConfigMap/Secret is updated, all the copies will be updated. Either delete the source ConfigMap/Secret or remove the annotation from the source ConfigMap/Secret to remove the copies.

If the list of contexts specified by the annotation is updated, Config Syncer will synchronize the ConfigMap/Secret accordingly, ie. it will create ConfigMap/Secret  in the clusters listed in new annotation (if not already exists) and delete ConfigMap/Secret from the clusters that were synced before but not listed in new annotation.

Note that, Config Syncer will error out if multiple contexts listed in annotation point same cluster. Also Config Syncer assumes that none of cluster contexts in `kubeconfig` file points the source cluster.

## Before You Begin

At first, you need to have a Kubernetes cluster and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [kind](https://kind.sigs.k8s.io/docs/user/quick-start/).

## Deploy Config Syncer

To enable config syncer for different clusters, you need a `kubeconfig` file consisting cluster contexts where you want to sync your ConfigMap/Secret.

```yaml
$ cat ./docs/examples/cluster-syncer/demo-kubeconfig.yaml

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

Now, deploy Config Syncer operator in your cluster following the steps [here](/docs/setup/install.md). Below you can see the command to install Config Syncer using Helm 3.

```console
$ helm install kubed appscode/kubed \
  --version {{< param "info.version" >}} \
  --namespace kube-system \
  --set config.clusterName=kind \
  --set config.kubeconfigContent="$(cat ./docs/examples/cluster-syncer/demo-kubeconfig.yaml)"
```

Once the operator pod is running, go to the next section.

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
$ kubectl annotate configmap omni kubed.appscode.com/sync-contexts="context-1,context-2" -n demo
configmap "omni" annotated
```

It will create configmap "omni" in `cluster-1` and `cluster-2`. For `cluster-1` it will sync into source namespace `demo`  since no namespace specified in `context-1` and for `cluster-2` it will sync into `demo-cluster-2` namespace since namespace specified in `context-2`. Here we assume that those namespaces already exits in the respective clusters.

Other concepts like updating source configmap, removing annotation, origin annotation, origin labels, etc. are similar to the tutorial described [here](/docs/guides/config-syncer/intra-cluster.md).

## Next Steps

- Need to keep some configuration synchronized across namespaces? Try [Config Syncer config syncer](/docs/guides/config-syncer/intra-cluster.md).
- Want to hack on Config Syncer? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
