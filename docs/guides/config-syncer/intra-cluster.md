---
title: Synchronize Configuration across Namespaces
description: Synchronize Configuration across Namespaces
menu:
  product_kubed_{{ .version }}:
    identifier: intra-cluster-syncer
    name: Across Namespaces
    parent: config-syncer
    weight: 10
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: guides
---

> New to Kubed? Please start [here](/docs/concepts/README.md).

# Synchronize Configuration across Namespaces

Say, you are using some Docker private registry. You want to keep its image pull secret synchronized across all namespaces of a Kubernetes cluster. Kubed can do that for you. If a ConfigMap or a Secret has the annotation __`kubed.appscode.com/sync: ""`__, Kubed will create a copy of that ConfigMap/Secret in all existing namespaces. Kubed will also create this ConfigMap/Secret, when you create a new namespace.

If you want to synchronize ConfigMap/Secret to some selected namespaces instead of all namespaces, you can do that by specifying namespace label-selector in the annotation. For example: __`kubed.appscode.com/sync: "app=kubed"`__. Kubed will create a copy of that  ConfigMap/Secret in all namespaces that matches the label-selector. Kubed will also create this Configmap/Secret in newly created namespace if it matches the label-selector.

If the data in the source ConfigMap/Secret is updated, all the copies will be updated. Either delete the source ConfigMap/Secret or remove the annotation from the source ConfigMap/Secret to remove the copies. If the namespace with the source ConfigMap/Secret is deleted, the copies will also be deleted.

If the value of label-selector specified by annotation is updated, Kubed will synchronize the ConfigMap/Secret accordingly, ie. it will create ConfigMap/Secret in the namespaces that are selected by new label-selector (if not already exists) and delete from namespaces that were synced before but not selected by new label-selector.

## Before You Begin

At first, you need to have a Kubernetes cluster and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [kind](https://kind.sigs.k8s.io/docs/user/quick-start/).

## Synchronize ConfigMap

In this tutorial, a ConfigMap will be synced across all Kubernetes namespaces using Kubed. You can do the same for Secrets.

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create namespace demo
namespace "demo" created

$ kubectl get namespaces
NAME          STATUS    AGE
default       Active    6h
kube-public   Active    6h
kube-system   Active    6h
demo          Active    4m
```

Now, create a ConfigMap called `omni` in the `demo` namespace. This will be our source ConfigMap.

```console
$ kubectl apply -f ./docs/examples/config-syncer/demo-0.yaml
configmap "omni" created

$ kubectl get configmaps --all-namespaces | grep omni
demo          omni                                 2         7m
```

```yaml
$ kubectl get configmaps omni -n demo -o yaml
apiVersion: v1
data:
  you: only
  leave: once
kind: ConfigMap
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"leave":"once","you":"only"},"kind":"ConfigMap","metadata":{"annotations":{},"name":"omni","namespace":"demo"}}
  creationTimestamp: 2017-07-26T13:20:15Z
  name: omni
  namespace: demo
  resourceVersion: "10598"
  selfLink: /api/v1/namespaces/demo/configmaps/omni
  uid: 2988e9d5-7205-11e7-af79-08002738e55e
```

Now, apply the `kubed.appscode.com/sync: ""` annotation to ConfigMap `omni`. Kubed operator will notice that and copy the ConfigMap in all existing namespaces.

```console
$ kubectl annotate configmap omni kubed.appscode.com/sync="" -n demo
configmap "omni" annotated

$ kubectl get configmaps --all-namespaces | grep omni
default       omni                                 2         1m
demo          omni                                 2         8m
kube-public   omni                                 2         1m
kube-system   omni                                 2         1m
```

```yaml
$ kubectl get configmaps omni -n demo -o yaml
apiVersion: v1
data:
  you: only
  leave: once
kind: ConfigMap
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"leave":"once","you":"only"},"kind":"ConfigMap","metadata":{"annotations":{},"name":"omni","namespace":"demo"}}
    kubed.appscode.com/sync: "true"
  creationTimestamp: 2017-07-26T13:20:15Z
  name: omni
  namespace: demo
  resourceVersion: "11053"
  selfLink: /api/v1/namespaces/demo/configmaps/omni
  uid: 2988e9d5-7205-11e7-af79-08002738e55e
```

Now, create a new namespace called `other`. Kubed will copy ConfigMap `omni` into that namespace.

```console
$ kubectl create ns other
namespace "other" created

$ kubectl get configmaps --all-namespaces | grep omni
default       omni                                 2         5m
demo          omni                                 2         12m
kube-public   omni                                 2         5m
kube-system   omni                                 2         5m
other         omni                                 2         1m
```

Alas! there is a typo is the ConfigMap data. Let's fix that.

```console
$ kubectl apply -f ./docs/examples/config-syncer/demo-1.yaml
configmap "omni" configured

$ kubectl get configmaps --all-namespaces | grep omni
default       omni                                 2         9m
demo          omni                                 2         16m
kube-public   omni                                 2         9m
kube-system   omni                                 2         9m
other         omni                                 2         5m
```

```yaml
$ kubectl get configmaps omni -n other -o yaml
apiVersion: v1
data:
  you: only
  live: once
kind: ConfigMap
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"live":"once","you":"only"},"kind":"ConfigMap","metadata":{"annotations":{},"name":"omni","namespace":"demo"}}
  creationTimestamp: 2017-07-26T13:31:13Z
  name: omni
  namespace: other
  resourceVersion: "11594"
  selfLink: /api/v1/namespaces/other/configmaps/omni
  uid: b193f40f-7206-11e7-af79-08002738e55e
```

Kubed operator notices that the source ConfigMap `omni` has been updated and propagates the change to all the copies in other namespaces.

## Namespace Selector

Lets' change annotation value of source ConfigMap `omni`.

```console
$ kubectl annotate configmap omni kubed.appscode.com/sync="app=kubed" -n demo --overwrite
configmap "omni" annotated

$ kubectl get configmaps --all-namespaces | grep omni
demo          omni                                 2         8m
```

Kubed operator removes the ConfigMap from all namespaces (except source) since no namespace matches the label-selector `app=kubed`.
Now, lets' apply `app=kubed` annotation to `other` namespace. Kubed operator will then sync the ConfigMap to `other` namespace.

```console
$ kubectl label namespace other app=kubed
namespace "other" labeled

$ kubectl get configmaps --all-namespaces | grep omni
demo          omni                                 2         8m
other         omni                                 2         5m
```

## Restricting Source Namespace

By default, Kubed will watch all namespaces for configmaps and secrets with `kubed.appscode.com/sync` annotation. But you can restrict the source namespace for configmaps and secrets by passing `config.configSourceNamespace` value during installation.

```console
$ helm install kubed appscode/kubed \
  --namespace=kube-system \
  --set imagePullPolicy=Always \
  --set config.configSourceNamespace=demo
```

## Remove Annotation

Now, lets' remove the annotation from source ConfigMap `omni`. Please note that `-` after annotation key `kubed.appscode.com/sync-`. This tells kubectl to remove this annotation from ConfigMap `omni`.

```console
$ kubectl annotate configmap omni kubed.appscode.com/sync- -n demo
configmap "omni" annotated

$ kubectl get configmaps --all-namespaces | grep omni
demo          omni                                 2         18m
```

## Origin Annotation

Since 0.9.0, Kubed operator will apply `kubed.appscode.com/origin` annotation on ConfigMap or Secret copies.

![origin annotation](/docs/images/config-syncer/config-origin.png)

## Origin Labels

Kubed  operator will apply following labels on ConfigMap or Secret copies:

- `kubed.appscode.com/origin.name`
- `kubed.appscode.com/origin.namespace`
- `kubed.appscode.com/origin.cluster`

This annotations are used by Kubed operator to list the copies for a specific source ConfigMap/Secret.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run the following commands:

```console
$ kubectl delete ns other
namespace "other" deleted

$ kubectl delete ns demo
namespace "demo" deleted
```

To uninstall Kubed operator, please follow the steps [here](/docs/setup/uninstall.md).

## Next Steps

- Learn how to sync config-maps or secrets across multiple cluster [here](/docs/guides/config-syncer/inter-cluster.md).
- Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
