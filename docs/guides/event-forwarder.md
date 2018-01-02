---
title: Event Forwarder
description: Event Forwarder
menu:
  product_kubed_0.3.1:
    identifier: tutorials-event-forwarder
    name: Event Forwarder
    parent: tutorials
    weight: 30
product_name: kubed
menu_name: product_kubed_0.3.1
section_menu_id: tutorials
---

> New to Kubed? Please start [here](/docs/guides/README.md).

# Forward Cluster Events
Kubed can send notifications via Email, SMS or Chat for various cluster events. This tutorial will show you how to use Kubed to setup an event forwarder.


## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).


## Deploy Kubed
To enable config syncer, you need a cluster config like below.

```yaml
$ cat ./docs/examples/event-forwarder/config.yaml

clusterName: unicorn
eventForwarder:
  nodeAdded:
    handle: true
  csrEvents:
    handle: true
  storageAdded:
    handle: true
  ingressAdded:
    handle: true
  warningEvents:
    handle: true
    namespaces:
    - kube-system
  receivers:
  - notifier: Mailgun
    to:
    - ops@example.com
notifierSecretName: notifier-config
```

| Key                                       | Description                                                                                           |
|-------------------------------------------|-------------------------------------------------------------------------------------------------------|
| `eventForwarder.nodeAdded.handle`         | `Optional`. If set to true, notifications are sent when a Node is added.                              |
| `eventForwarder.csrEvents.handle`         | `Optional`. If set to true, notifications are sent when a [CertificateSigningRequest](https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/) is added, approved or denied. |
| `eventForwarder.ingressAdded.handle`      | `Optional`. If set to true, notifications are sent when an Ingress is added.                          |
| `eventForwarder.ingressAdded.namespaces`  | `Optional`. If set, notifications are sent only when Ingress are added in these namespaces. Otherwise, notifications are sent when Ingress are added in any namespace |
| `eventForwarder.storageAdded.handle`      | `Optional`. If set to true, notifications are sent when a StorageClass/PV/PVC is added.               |
| `eventForwarder.storageAdded.namespaces`  | `Optional`. If set, notifications are sent only when PVC are added in these namespaces. Otherwise, notifications are sent when PVC added in any namespace. Since StorageClass and PV are non-namespaced resource, this field has not effect on these. |
| `eventForwarder.warningEvents.handle`     | `Optional`. If set to true, notifications are sent when a `Warning` Event is added.                   |
| `eventForwarder.warningEvents.namespaces` | `Optional`. If set, notifications are sent only when warning events are added in these namespaces. Otherwise, notifications are sent when warning events are added in any namespace |
| `eventForwarder.receiver`                 | `Required`. To learn how to use various notifiers, please visit [here](/docs/guides/notifiers.md). |
| `clusterName`                             | `Optional`. A meaningful identifer for cluster. This cluster name will be prefixed to any notification sent via Email/SMS/Chat so that you can identify the source easily. |

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
  config.yaml: ZXZlbnRGb3J3YXJkZXI6CiAgbm9kZUFkZGVkOgogICAgaGFuZGxlOiB0cnVlCiAgc3RvcmFnZUFkZGVkOgogICAgaGFuZGxlOiB0cnVlCiAgaW5ncmVzc0FkZGVkOgogICAgaGFuZGxlOiB0cnVlCiAgd2FybmluZ0V2ZW50czoKICAgIGhhbmRsZTogdHJ1ZQogICAgbmFtZXNwYWNlczoKICAgIC0ga3ViZS1zeXN0ZW0KICByZWNlaXZlcjoKICAgIG5vdGlmaWVyOiBtYWlsZ3VuCiAgICB0bzoKICAgIC0gb3BzQGV4YW1wbGUuY29tCm5vdGlmaWVyU2VjcmV0TmFtZToga3ViZWQtbm90aWZpZXIK
kind: Secret
metadata:
  creationTimestamp: 2017-07-27T05:35:54Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "70583"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 753220c3-728d-11e7-87f5-08002738e55e
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/setup/install.md). Once the operator pod is running, go to the next section.


## Test Forwarder
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

### Forward Storage Added Event  
In this section, a PVC will be used to show how event forwarder feature can be used. Create a PVC called `myclaim` in the `demo` namespace.

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

### Forward Warning Events
In this section, a Busybox pod will be used to show how warning events are forwarded. Create a Pod called `busybox` in the `demo` namespace.

```yaml
$ cat ./docs/examples/event-forwarder/demo-1.yaml

apiVersion: v1
kind: Pod
metadata:
  name: busybox
  namespace: demo
spec:
  restartPolicy: Never
  containers:
  - name: busybox
    image: busybox
    imagePullPolicy: IfNotPresent
    command:
      - bad
      - "3600"
```
```console
$ kubectl apply -f ./docs/examples/event-forwarder/demo-1.yaml
pod "busybox" created

$ kubectl get pods -n demo --show-all
NAME      READY     STATUS                                                                                                                                                                                                      RESTARTS   AGE
busybox   0/1       rpc error: code = 2 desc = failed to start container "bcc25386c0c9421b04ce9c574405917fc4940a0b324a2b062f02978c46463f07": Error response from daemon: Container command 'bad' not found or does not exist.   0          10m
```

Here, the busybox pod fails to start because it uses a missing command called `bad`. This results in 2 `Warning` events. Now, check your GMail account. You should receive 2 emails like below.

![Pod Failed](/docs/images/event-forwarder/pod-fail-1.png)
![Pod FailedSync](/docs/images/event-forwarder/pod-fail-2.png)


## Filter by Namespaces
You can configure Kubed to forward events for a subset of namespaces. You can also disable sending events for a particular type. Here is an example `config.yaml`:
```yaml
clusterName: unicorn
eventForwarder:
  nodeAdded: {}
  ingressAdded:
    handle: false
  warningEvents:
    handle: true
    namespaces:
    - kube-system
  receivers:
  - notifier: Mailgun
    to:
    - ops@example.com
notifierSecretName: notifier-config
```

In the above example:
 - `eventForwarder.nodeAdded` is set to an empty object `{}`. This means `eventForwarder.nodeAdded.handle` is false. So, notifications are _not_ sent when Nodes are added.
 - `eventForwarder.ingressAdded.handle` is set to `false`. Notifications are _not_ sent when Ingress objects are added.
 - `eventForwarder.storageAdded` is missing. So, _no_ notifications are sent when StorageClass/PV/PVC etc are added.
 - `eventForwarder.warningEvents.handle` is set to `true`. Notifications are sent when Events are added in `kube-system` namespace.


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

To uninstall Kubed operator, please follow the steps [here](/docs/setup/uninstall.md).


## Next Steps
 - Learn how to use Kubed to take periodic snapshots of a Kubernetes cluster [here](/docs/guides/cluster-snapshot.md).
 - To setup a recycle bin for deleted and/or updated Kubernetes objects, please visit [here](/docs/guides/recycle-bin.md).
 - Need to keep some configuration synchronized across namespaces? Try [Kubed config syncer](/docs/guides/config-syncer.md).
 - Out of disk space because of too much logs in Elasticsearch or metrics in InfluxDB? Configure [janitors](/docs/guides/janitors.md) to delete old data.
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/ROADMAP.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
