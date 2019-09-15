---
title: Event Forwarder
description: Event Forwarder
menu:
  product_kubed_{{ .version }}:
    identifier: event-forwarder
    name: Event Forwarder
    parent: cluster-events
    weight: 10
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: guides
---

> New to Kubed? Please start [here](/docs/concepts/README.md).

# Forward Cluster Events
Kubed can send notifications via Email, SMS or Chat for various cluster events. This document will show you how to use Kubed to setup an event forwarder.


## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).


## Deploy Kubed
To enable config syncer, you need a cluster config like below.

```yaml
$ cat ./docs/examples/event-forwarder/config.yaml

clusterName: unicorn
eventForwarder:
  receivers:
  - notifier: Mailgun
    to:
    - ops@example.com
  rules:
  # notify for warning events in kube-system namespace
  - namespaces:
    - kube-system
    operations:
    - CREATE
    resources:
    - group: "" # core API group
      resources:
      - events
  # notify for both CREATE and DELETE operations in any namespace
  - resources:
    - group: ""  # core API group
      resources:
      - nodes
      - persistentvolumes
      - persistentvolumeclaims
    - group: storage.k8s.io
      resources:
      - storageclasses
    - group: extensions
      resources:
      - ingresses
    - group: voyager.appscode.com
      resources:
      - ingresses
    - group: certificates.k8s.io
      resources:
      - certificatesigningrequests
notifierSecretName: notifier-config
```

The configuration format is inpired by [audit policy file format](https://kubernetes.io/docs/tasks/debug-application-cluster/audit/). The policy is defined [here](https://github.com/appscode/kubed/blob/4d4f7b9d03a84910e04c52a6801a9b0f71fae8e7/apis/kubed/v1alpha1/types.go#L75).
The matcher logic is implemented [here](https://github.com/appscode/kubed/blob/4d4f7b9d03a84910e04c52a6801a9b0f71fae8e7/pkg/eventer/resourcehandler.go#L52).

**NB:** The event forwarder configuration format has been redesigned in 0.8.0 and should be updates accordingly if you are upgrading from a previous version.

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


## Supported Kubernetes Objects
Following Kubernetes objects are supported by event forwarder:

- __v1:__
  - ConfigMap
  - Event
  - LimitRange
  - Namespace
  - Node
  - PersistentVolume
  - PersistentVolumeClaim
  - ReplicationController
  - Secret
  - Service
  - ServiceAccount
- __apps/v1beta1:__
  - Deployment
  - StatefulSet
- __batch/v1:__
  - Job
- __batch/v1beta1:__
  - CronJob
- __extensions/v1beta1:__
  - Deployment
  - Ingress
  - ReplicaSet
- __networking.k8s.io/v1:__
  - NetworkPolicy
- __kubedb/v1alpha1:__
  - DormantDatabase
  - Elasticsearch
  - Memcached
  - MongoDB
  - MySQL
  - Postgres
  - Redis
  - Snapshot
- __monitoring.coreos.com/v1:__
  - Prometheus
  - ServiceMonitor
  - Alertmanager
- __rbac/v1:__
  - ClusterRole
  - ClusterRoleBinding
  - Role
  - RoleBinding
- __monitoring.appscode.com/v1alpha1:__
  - ClusterAlert
  - NodeAlert
  - PodAlert
- __stash.appscode.com/v1alpha1:__
  - Restic
  - Recovery
- __storage/v1:__
  - StorageClass
- __voyager.appscode.com/v1beta1:__
  - Certificate
  - Ingress

To add support for additional object types, please [file an issue](https://github.com/appscode/kubed/issues/new?title=Support+Object+Kind+[xyz]+in+SearchEngine).


## Disable Event Forwarder
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
 - See the list of supported notifiers [here](/docs/guides/cluster-events/notifiers.md).
 - Learn how to use Kubed to protect your Kubernetes cluster from disasters [here](/docs/guides/disaster-recovery/).
 - Need to keep configmaps/secrets synchronized across namespaces or clusters? Try [Kubed config syncer](/docs/guides/config-syncer/).
 - Out of disk space because of too much logs in Elasticsearch or metrics in InfluxDB? Configure [janitors](/docs/guides/janitors.md) to delete old data.
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
