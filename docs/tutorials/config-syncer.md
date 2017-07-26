> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Synchronize Configuration Across Namespaces
Sometimes you have some configuration that you want to synchronize across all Kubernetes namespaces. Kubed can do that for you. If a ConfigMap or Secret has the label `kubed.appscode.com/sync: true`, Kubed will create a similar ConfigMap / Secret in all existing namespaces. Kubed will also create this ConfigMap/Secret, when you create a new namespace. If the data in the source ConfigMap/Secret is updated, all the copies will be updated.

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

## Synchronize ConfigMap & Secret


















![GCS Snapshot](/docs/images/cluster-snapshot/gcs-snapshot.png)


Now, install Kubed operator in your cluster following the steps [here](/docs/install.md).



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


Sync ConfigMaps and Secrets


kubectl create configmap global-config -n demo --from-literal=special.how=very --from-literal=special.type=charm
kubectl label configmap global-config kubed.appscode.com/sync=true -n demo


~ $ kubectl create configmap global-config --from-literal=special.how=very --from-literal=special.type=charm
configmap "global-config" created
~ $ kubectl delete configmap global-config
configmap "global-config" deleted
~ $ 
~ $ 
~ $ kubectl create configmap global-config -n src --from-literal=special.how=very --from-literal=special.type=charm
configmap "global-config" created
~ $ kubectl label configmap global-config kubed.appscode.com/sync=true -n src
configmap "global-config" labeled
~ $ 
