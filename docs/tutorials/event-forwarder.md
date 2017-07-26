> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Forward Cluster Events
Kubed can 
Sometimes you have some configuration that you want to synchronize across all Kubernetes namespaces. Kubed can do that for you. If a ConfigMap or Secret has the annotation __`kubed.appscode.com/sync: true`__, Kubed will create a similar ConfigMap / Secret in all existing namespaces. Kubed will also create this ConfigMap/Secret, when you create a new namespace. If the data in the source ConfigMap/Secret is updated, all the copies will be updated. Either delete the source ConfigMap/Secret or remove the annotation from the source ConfigMap/Secret to purge the copies. If the namespace with the source ConfigMap/Secret is deleted, the copies are left intact.

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

