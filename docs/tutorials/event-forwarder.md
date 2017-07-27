> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Forward Cluster Events
Kubed can send notifications via Email, SMS or Chat for various cluster events. This tutorial will show you how to use Kubed to setup an event forwarder.


## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).


## Deploy Kubed
To enable config syncer, you need a cluster config like below.

```yaml
$ cat ./docs/examples/event-forwarder/config.yaml

eventForwarder:
  nodeAdded: {}
  ingressAdded: {}
  storageAdded: {}
  warningEvents: {}
  receiver:
    notifier: mailgun
    to:
    - ops@example.com
notifierSecretName: kubed-notifier
```

| Key                            | Description                                                                       |
|--------------------------------|-----------------------------------------------------------------------------------|
| `eventForwarder.nodeAdded`     | `Optional`. If set, notifications are sent when a Node is added.                  |
| `eventForwarder.ingressAdded`  | `Optional`. If set, notifications are sent when an Ingress is added.              |
| `eventForwarder.storageAdded`  | `Optional`. If set, notifications are sent when a StorageClass/PV/PVC is added.   |
| `eventForwarder.warningEvents` | `Optional`. If set, notifications are sent when a `Warning` Event is added.       |
| `eventForwarder.receiver`      | `Required`. To learn how to use various notifiers, please visit [here](/docs/tutorials/notifiers.md). |

Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/event-forwarder/config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  config.yaml: ZXZlbnRGb3J3YXJkZXI6CiAgbm9kZUFkZGVkOiB7fQogIGluZ3Jlc3NBZGRlZDoge30KICBzdG9yYWdlQWRkZWQ6IHt9CiAgd2FybmluZ0V2ZW50czoge30KICByZWNlaXZlcjoKICAgIG5vdGlmaWVyOiBtYWlsZ3VuCiAgICB0bzoKICAgIC0gb3BzQGV4YW1wbGUuY29tCm5vdGlmaWVyU2VjcmV0TmFtZToga3ViZWQtbm90aWZpZXIK
kind: Secret
metadata:
  creationTimestamp: 2017-07-27T00:53:01Z
  name: kubed-config
  namespace: kube-system
  resourceVersion: "56568"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: f0cb8f70-7265-11e7-af79-08002738e55e
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/install.md). Once the operator pod is running, go to the next section.


## Test Forwarder
In this tutorial, a PVC will be used to show how event forwarder feature can be used.

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

Create a PVC called `myclaim` in the `demo` namespace.

```console
$ kubectl apply -f ./docs/examples/event-forwarder/demo-0.yaml
persistentvolumeclaim "myclaim" configured
```
```yaml
$ kubectl get pvc myclaim -n demo -o yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  annotations:
    control-plane.alpha.kubernetes.io/leader: '{"holderIdentity":"a56b7269-71ef-11e7-af79-08002738e55e","leaseDurationSeconds":15,"acquireTime":"2017-07-27T01:24:08Z","renewTime":"2017-07-27T01:24:10Z","leaderTransitions":0}'
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"PersistentVolumeClaim","metadata":{"annotations":{},"name":"myclaim","namespace":"demo"},"spec":{"accessModes":["ReadWriteOnce"],"resources":{"requests":{"storage":"50Mi"}},"storageClassName":"standard"}}
    pv.kubernetes.io/bind-completed: "yes"
    pv.kubernetes.io/bound-by-controller: "yes"
    volume.beta.kubernetes.io/storage-provisioner: k8s.io/minikube-hostpath
  creationTimestamp: 2017-07-27T01:24:08Z
  name: myclaim
  namespace: demo
  resourceVersion: "58641"
  selfLink: /api/v1/namespaces/demo/persistentvolumeclaims/myclaim
  uid: 49b9851c-726a-11e7-af79-08002738e55e
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 50Mi
  storageClassName: standard
  volumeName: pvc-49b9851c-726a-11e7-af79-08002738e55e
status:
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 50Mi
  phase: Bound
```

Now, assuming you configured a GMail account as the receiver for events, you should see an email like below:

![PVC Added Notification](/docs/images/event-forwarder/pvc-added-notification.png)


## Filter by Namespaces
You can configure Kubed to forward events for a subset of namespaces. You can also disable sending events for a particular type. Here is an example `config.yaml`:
```yaml
eventForwarder:
  nodeAdded: {}
  ingressAdded:
    namespaces:
    - default
  warningEvents:
    namespaces:
    - kube-system
  receiver:
    notifier: mailgun
    to:
    - ops@example.com
notifierSecretName: kubed-notifier
```

In the above example:
 - `eventForwarder.nodeAdded` is set to an empty object `{}`. So, notifications are sent when Nodes are added.
 - `eventForwarder.ingressAdded` is set. Notifications are sent when Ingress objects are added in `default` namespace.
 - `eventForwarder.storageAdded` is missing. So, _no_ notifications are sent when StorageClass/PV/PVC etc are added.
 - `eventForwarder.warningEvents` is set. Notifications are sent when Events are added in `kube-system` namespace.


## Disable Recycle Bin
If you would like to disable this feature, remove the `eventForwarder` portion of your Kubed cluster config. Then update the `kubed-config` Secret and restart Kubed operator pod(s).


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run the following commands:
```console
$ kubectl delete pvc myclaim -n demo
persistentvolumeclaim "myclaim" deleted

$ kubectl delete ns demo
namespace "demo" deleted
```

To uninstall Kubed operator, please follow the steps [here](/docs/uninstall.md).

## Next Steps
