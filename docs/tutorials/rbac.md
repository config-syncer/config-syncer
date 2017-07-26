# Using RBAC with KubeDB
This tutorial will show you how to use KubeDB in a [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/) enabled cluster.

## Before You Begin
At first, you need to have a RBAC enabled Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). To create a RBAC enabled cluster using MiniKube, follow the instructions below:

1. If you are currently running a Minukube cluster without RBAC, delete the cluster. This will delete any objects running in the cluster.
```console
$ minikube delete
```

2. Now, create a RBAC cluster with RBAC enabled.
```console
$ minikube start --extra-config=apiserver.Authorization.Mode=RBAC
```

3. Once the cluster is up and running, you need to set ServiceAccount for the `kube-dns` addon to successfully run it.
```console
# Wait for kube-dns deployment to be created.
$  kubectl get deployment -n kube-system --watch

# create kube-dns ServiceAccount
$ kubectl create serviceaccount kube-dns -n kube-system

# Patch kube-dns Deployment to set service account for pods.
$ kubectl patch deployment kube-dns -n kube-system -p '{"spec":{"template":{"spec":{"serviceAccountName":"kube-dns"}}}}'

# Wait for kube-dns pods to start running
$ kubectl get pods -n kube-system --watch
```

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/install.md).
```
$ kubedb init --rbac
```

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. This tutorial will also use a PGAdmin to connect and test PostgreSQL database, once it is running. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f ./docs/examples/rbac/demo-0.yaml
namespace "demo" created
deployment "pgadmin" created
service "pgadmin" created

$ kubectl get pods -n demo --watch
NAME                      READY     STATUS              RESTARTS   AGE
pgadmin-538449054-s046r   0/1       ContainerCreating   0          13s
pgadmin-538449054-s046r   1/1       Running   0          1m
^C⏎                                                                                                                                                             

$ kubectl get service -n demo
NAME      CLUSTER-IP   EXTERNAL-IP   PORT(S)        AGE
pgadmin   10.0.0.92    <pending>     80:31188/TCP   1m

$ minikube ip
192.168.99.100
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{pgadmin-svc-nodeport}_. According to the above example, this URL will be [http://192.168.99.100:31188](http://192.168.99.100:31188). To log into the PGAdmin, use username `admin` and password `admin`.

## Create a PostgreSQL database
KubeDB implements a `Postgres` TPR to define the specification of a PostgreSQL database. Below is the `Postgres` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: p1
  namespace: demo
spec:
  version: 9.5
  doNotPause: true
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi

$ kubedb create -f ./docs/examples/rbac/demo-1.yaml
validating "./docs/examples/rbac/demo-1.yaml"
postgres "p1" created
```

Here,
 - `spec.version` is the version of PostgreSQL database. In this tutorial, a PostgreSQL 9.5 database is going to be created.

 - `spec.doNotPause` tells KubeDB operator that if this tpr is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.

 - `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

KubeDB operator watches for `Postgres` objects using Kubernetes api. When a `Postgres` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching tpr name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present.

```console
$ kubedb describe pg -n demo p1
Name:		p1
Namespace:	demo
StartTimestamp:	Mon, 17 Jul 2017 15:31:34 -0700
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

Service:	
  Name:		p1
  Type:		ClusterIP
  IP:		10.0.0.161
  Port:		db	5432/TCP

Database Secret:
  Name:	p1-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	35 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                Type       Reason               Message
  ---------   --------   -----     ----                --------   ------               -------
  2m          2m         1         Postgres operator   Normal     SuccessfulValidate   Successfully validate Postgres
  2m          2m         1         Postgres operator   Normal     SuccessfulCreate     Successfully created Postgres
  2m          2m         1         Postgres operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  3m          3m         1         Postgres operator   Normal     SuccessfulValidate   Successfully validate Postgres
  3m          3m         1         Postgres operator   Normal     Creating             Creating Kubernetes objects


$ kubectl get statefulset -n demo
NAME      DESIRED   CURRENT   AGE
p1        1         1         1m

$ kubectl get pvc -n demo
NAME        STATUS    VOLUME                                     CAPACITY   ACCESSMODES   STORAGECLASS   AGE
data-p1-0   Bound     pvc-e90b87d4-6b5a-11e7-b9ca-080027f73ab7   50Mi RWO           standard       1m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS    CLAIM            STORAGECLASS   REASON    AGE
pvc-e90b87d4-6b5a-11e7-b9ca-080027f73ab7   50Mi RWO           Delete          Bound     demo/data-p1-0   standard                 1m

$ kubectl get service -n demo
NAME      CLUSTER-IP   EXTERNAL-IP   PORT(S)        AGE
kubedb    None         <none>                       3m
p1        10.0.0.143   <none>        5432/TCP       3m
pgadmin   10.0.0.120   <pending>     80:30576/TCP   6m
```

Since RBAC is enabled, a ClusterRole, ServiceAccount and ClusterRoleBinding with the matching tpr name is also created and used as the service account name for the corresponding StatefulSet. This Role is used by Prometheus exporter sidecar container to connect to the database.

```yaml
$ kubectl get role -n demo p1 -o yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  creationTimestamp: 2017-07-18T20:35:43Z
  name: p1
  namespace: demo
  resourceVersion: "1308"
  selfLink: /apis/rbac.authorization.k8s.io/v1beta1/namespaces/demo/roles/p1
  uid: ab72299c-6bf8-11e7-ab55-080027815c31
rules:
- apiGroups:
  - kubedb.com
  resourceNames:
  - p1
  resources:
  - postgreses
  verbs:
  - get
- apiGroups:
  - ""
  resourceNames:
  - p1-admin-auth
  resources:
  - secrets
  verbs:
  - get

$ kubectl get serviceaccount -n demo p1 -o yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  creationTimestamp: 2017-07-18T20:35:43Z
  name: p1
  namespace: demo
  resourceVersion: "1317"
  selfLink: /api/v1/namespaces/demo/serviceaccounts/p1
  uid: ab73147f-6bf8-11e7-ab55-080027815c31
secrets:
- name: p1-token-qxdtf

$ kubectl get rolebindings -n demo p1 -o yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  creationTimestamp: 2017-07-18T20:35:43Z
  name: p1
  namespace: demo
  resourceVersion: "1310"
  selfLink: /apis/rbac.authorization.k8s.io/v1beta1/namespaces/demo/rolebindings/p1
  uid: ab7d1eb4-6bf8-11e7-ab55-080027815c31
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: p1
subjects:
- kind: ServiceAccount
  name: p1
  namespace: demo
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified tpr:


```yaml
$ kubedb get pg -n demo p1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  creationTimestamp: 2017-07-17T22:31:34Z
  name: p1
  namespace: demo
  resourceVersion: "2677"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/postgreses/p1
  uid: b02ccec1-6b3f-11e7-bdc0-080027aa4456
spec:
  databaseSecret:
    secretName: p1-admin-auth
  doNotPause: true
  init:
    scriptSource:
      gitRepo:
        repository: https://github.com/k8sdb/postgres-init-scripts.git
      scriptPath: postgres-init-scripts/run.sh
  resources: {}
  storage:
    accessModes:
    - ReadWriteOnce
    storageClassName: standard
    resources:
      requests:
        storage: 50Mi
  version: "9.5"
status:
  creationTime: 2017-07-17T22:31:34Z
  phase: Running
```


Please note that KubeDB operator has created a new Secret called `p1-admin-auth` (format: {tpr-name}-admin-auth) for storing the password for `postgres` superuser. This secret contains a `.admin` key with a ini formatted key-value pairs. If you want to use an existing secret please specify that when creating the tpr using `spec.databaseSecret.secretName`.

Now, you can connect to this database from the PGAdmin dashboard using the database pod IP and `postgres` user password. 

```console
$ kubectl get pods p1-0 -n demo -o yaml | grep IP
  hostIP: 192.168.99.100
  podIP: 172.17.0.6

$ kubectl get secrets -n demo p1-admin-auth -o jsonpath='{.data.\.admin}' | base64 -d
POSTGRES_PASSWORD=R9keKKRTqSJUPtNC
```


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall KubeDB operator, please follow the steps [here](/docs/uninstall.md).


## Next Steps
- Learn how to use KubeDB to run a PostgreSQL database [here](/docs/tutorials/postgres.md).
- Learn how to use KubeDB to run an Elasticsearch database [here](/docs/tutorials/elasticsearch.md).
- Wondering what features are coming next? Please visit [here](/ROADMAP.md). 
- Want to hack on KubeDB? Check our [contribution guidelines](/CONTRIBUTING.md).
