> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Synchronize Configuration Across Namespaces
Sometimes you have some configuration that you want to synchronize across all Kubernetes namespaces. Kubed can do that for you. If a ConfigMap or Secret has the annotation `kubed.appscode.com/sync: true`, Kubed will create a similar ConfigMap / Secret in all existing namespaces. Kubed will also create this ConfigMap/Secret, when you create a new namespace. If the data in the source ConfigMap/Secret is updated, all the copies will be updated. Either delete the source ConfigMap/Secret or remove the annotation from the source ConfigMap/Secret to purge the copies. If the namespace with the source ConfigMap/Secret is deleted, the copies are left intact.

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

## Deploy Kubed
To enable config syncer, you need a cluster config like below.

```yaml
$ cat ./docs/examples/config-syncer/config.yaml

enableConfigSyncer: true
```

| Key                   | Description                                                                               |
|-----------------------|-------------------------------------------------------------------------------------------|
| `enableConfigSyncer`  | `Required`. If set to `true`, ConfigMap/Secret synchronization operation will be enabled. |


Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/config-syncer/config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  config.yaml: ZW5hYmxlQ29uZmlnU3luY2VyOiB0cnVlCg==
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

Now, deploy Kubed operator in your cluster following the steps [here](/docs/install.md). Once the operator pod is running, go to the next section.

## Synchronize ConfigMap
In this tutorial, a ConfigMap will be synced across Kubernetes namespaces using Kubed. You can do the same using Secrets.

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create namespace demo
namespace "demo" created

~ $ kubectl get namespaces
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

Now, apply the `kubed.appscode.com/sync: true` annotaiotn to ConfigMap `omni`.

```console
$ kubectl annotate configmap omni kubed.appscode.com/sync=true -n demo
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
  leave: once
  you: only
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










![GCS Snapshot](/docs/images/cluster-snapshot/gcs-snapshot.png)


Now, install Kubed operator in your cluster following the steps [here](/docs/install.md).





Sync ConfigMaps and Secrets


kubectl create configmap omni -n demo --from-literal=special.how=very --from-literal=special.type=charm
kubectl annotate configmap omni kubed.appscode.com/sync=true -n demo

$ kubectl get configmaps --all-namespaces | grep omni
default       omni                                 2         1m
demo          omni                                 2         1m
kube-public   omni                                 2         1m
kube-system   omni                                 2         1m




~ $ kubectl create configmap omni --from-literal=special.how=very --from-literal=special.type=charm
configmap "omni" created
~ $ kubectl delete configmap omni
configmap "omni" deleted
~ $ 
~ $ 
~ $ kubectl create configmap omni -n src --from-literal=special.how=very --from-literal=special.type=charm
configmap "omni" created
~ $ kubectl label configmap omni kubed.appscode.com/sync=true -n src
configmap "omni" labeled
~ $ 



$ kubectl annotate configmap omni kubed.appscode.com/sync- -n demo
configmap "omni" annotated

$ kubectl get configmaps --all-namespaces | grep omni
demo          omni                                 2         18m
