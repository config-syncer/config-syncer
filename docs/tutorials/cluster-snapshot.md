> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Cluster Snapshots
Kubed supports taking periodic snapshot of a Kubernetes cluster objects. The snapshot data can be stored in various cloud providers, eg, Amazon S3, Google Cloud Storage, Microsoft Azure, OpenStack Swift and any locally mounted volumes like NFS, GlusterFS, etc. Kubed uses Kubernetes discovery api to find all available resources in a cluster and stores them in a file matching the `selfLink` URL for an object. This tutorial will show you how to use Kubed to take periodic snapshots of a Kubernetes cluster objects.

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

## Google Cloud Storage (GCS)
In this section, we are going to use Google Cloud Storage to store snapshot data. To configure this backend, a Kubertnetes Secret with the following keys is needed:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic gcs-secret -n kube-system \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "gcs-secret" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret gcs-secret app=kubed -n kube-system
secret "gcs-secret" labeled
```

```yaml
$ kubectl get secret gcs-secret -n kube-system -o yaml
apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInB...tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T04:37:44Z
  labels:
    app: kubed
  name: gcs-secret
  namespace: kube-system
  resourceVersion: "1671"
  selfLink: /api/v1/namespaces/kube-system/secrets/gcs-secret
  uid: 2aacabc8-71bc-11e7-a5ec-0800273df5f2
type: Opaque
```

Now, let's take a look at the cluster config. Here,

```yaml
$ cat ./docs/examples/cluster-snapshot/gcs/config.yaml

snapshotter:
  Storage:
    gcs:
      bucket: bucket-for-snapshot
      prefix: minikube
    storageSecretName: gcs-secret
  sanitize: true
  schedule: '@every 6h'
```

| Key                                     | Description                                                                     |
|-----------------------------------------|---------------------------------------------------------------------------------|
| `snapshotter.storage.storageSecretName` | `Required`. Name of storage secret                                              |
| `snapshotter.storage.gcs.bucket`        | `Required`. Name of GCS Bucket                                                  |
| `snapshotter.storage.gcs.prefix`        | `Optional`. Path prefix into bucket where snapshot will be stored               |
| `snapshotter.sanitize`                  | `Optional`. If set to `true`, various auto generated ObjectMeta and Spec fields are cleaned up before storing snapshots |
| `snapshotter.schedule`                  | `Required`. [Cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) specifying the schedule for snapshot operations. |


Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/cluster-snapshot/gcs/config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  config.yaml: c25hcHNob3R0ZXI6CiAgU3RvcmFnZToKICAgIGdjczoKICAgICAgYnVja2V0OiBidWNrZXQtZm9yLXNuYXBzaG90CiAgICAgIHByZWZpeDogbWluaWt1YmUKICAgIHN0b3JhZ2VTZWNyZXROYW1lOiBnY3Mtc2VjcmV0CiAgc2FuaXRpemU6IHRydWUKICBzY2hlZHVsZTogJ0BldmVyeSA2aCcK
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T05:51:11Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "6831"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 6d5babe7-71c6-11e7-a5ec-0800273df5f2
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/install.md). Once the operator pod is running, check your bucket from Google Cloud console. You should see the data from initial snapshot operation.

![GCS Snapshot](/docs/images/cluster-snapshot/gcs-snapshot.png)


## AWS S3
Kubed supports Amazon S3 or [Minio](https://minio.io/) servers as snapshot storage backend. To configure this backend, create a Secret with the following secret keys:

| Key                     | Description                                                |
|-------------------------|------------------------------------------------------------|
| `AWS_ACCESS_KEY_ID`     | `Required`. AWS / Minio access key ID                      |
| `AWS_SECRET_ACCESS_KEY` | `Required`. AWS / Minio secret access key                  |

```console
$ echo -n '<your-aws-access-key-id-here>' > AWS_ACCESS_KEY_ID
$ echo -n '<your-aws-secret-access-key-here>' > AWS_SECRET_ACCESS_KEY
$ kubectl create secret generic s3-secret -n kube-system \
    --from-file=./AWS_ACCESS_KEY_ID \
    --from-file=./AWS_SECRET_ACCESS_KEY
secret "s3-secret" created

$ kubectl label secret s3-secret app=kubed -n kube-system
secret "s3-secret" labeled
```

```yaml
$ kubectl get secret s3-secret -n kube-system -o yaml

apiVersion: v1
data:
  AWS_ACCESS_KEY_ID: PHlvdXItYXdzLWFjY2Vzcy1rZXktaWQtaGVyZT4=
  AWS_SECRET_ACCESS_KEY: PHlvdXItYXdzLXNlY3JldC1hY2Nlc3Mta2V5LWhlcmU+
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T05:26:19Z
  labels:
    app: kubed
  name: s3-secret
  namespace: kube-system
  resourceVersion: "5180"
  selfLink: /api/v1/namespaces/kube-system/secrets/s3-secret
  uid: f4353b86-71c2-11e7-a5ec-0800273df5f2
type: Opaque
```

Now, let's take a look at the cluster config. Here,

```yaml
$ cat ./docs/examples/cluster-snapshot/s3/config.yaml

snapshotter:
  Storage:
    s3:
      endpoint: 's3.amazonaws.com'
      bucket: kubedb-qa
      prefix: minikube
    storageSecretName: snap-secret
  sanitize: true
  schedule: '@every 6h'
```

| Key                                     | Description                                                                     |
|-----------------------------------------|---------------------------------------------------------------------------------|
| `snapshotter.storage.storageSecretName` | `Required`. Name of storage secret                                              |
| `snapshotter.storage.s3.bucket`         | `Required`. Name of S3 Bucket                                                   |
| `snapshotter.storage.s3.prefix`         | `Optional`. Path prefix into bucket where snapshot will be stored               |
| `snapshotter.sanitize`                  | `Optional`. If set to `true`, various auto generated ObjectMeta and Spec fields are cleaned up before storing snapshots |
| `snapshotter.schedule`                  | `Required`. [Cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) specifying the schedule for snapshot operations. |

Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/cluster-snapshot/s3/config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  config.yaml: c25hcHNob3R0ZXI6CiAgU3RvcmFnZToKICAgIHMzOgogICAgICBlbmRwb2ludDogJ3MzLmFtYXpvbmF3cy5jb20nCiAgICAgIGJ1Y2tldDoga3ViZWRiLXFhCiAgICAgIHByZWZpeDogbWluaWt1YmUKICAgIHN0b3JhZ2VTZWNyZXROYW1lOiBzbmFwLXNlY3JldAogIHNhbml0aXplOiB0cnVlCiAgc2NoZWR1bGU6ICdAZXZlcnkgNmgn
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T05:32:09Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "5572"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: c4dc8215-71c3-11e7-a5ec-0800273df5f2
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/install.md). Once the operator pod is running, check your bucket from S3 console. You should see the data from initial snapshot operation.


## Microsoft Azure Storage
Kubed supports Microsoft Azure Storage as snapshot storage backend. To configure this backend, create a Secret with the following secret keys:

| Key                     | Description                                                |
|-------------------------|------------------------------------------------------------|
| `AZURE_ACCOUNT_NAME`    | `Required`. Azure Storage account name                     |
| `AZURE_ACCOUNT_KEY`     | `Required`. Azure Storage account key                      |

```console
$ echo -n '<your-azure-storage-account-name>' > AZURE_ACCOUNT_NAME
$ echo -n '<your-azure-storage-account-key>' > AZURE_ACCOUNT_KEY
$ kubectl create secret generic azure-secret -n kube-system \
    --from-file=./AZURE_ACCOUNT_NAME \
    --from-file=./AZURE_ACCOUNT_KEY
secret "azure-secret" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret azure-secret app=kubed -n kube-system
secret "azure-secret" labeled
```

```yaml
$ kubectl get secret azure-secret -n kube-system -o yaml

apiVersion: v1
data:
  AZURE_ACCOUNT_KEY: PHlvdXItYXp1cmUtc3RvcmFnZS1hY2NvdW50LWtleT4=
  AZURE_ACCOUNT_NAME: PHlvdXItYXp1cmUtc3RvcmFnZS1hY2NvdW50LW5hbWU+
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T05:58:21Z
  labels:
    app: kubed
  name: azure-secret
  namespace: kube-system
  resourceVersion: "7427"
  selfLink: /api/v1/namespaces/kube-system/secrets/azure-secret
  uid: 6e197570-71c7-11e7-a5ec-0800273df5f2
type: Opaque
```

Now, let's take a look at the cluster config. Here,

```yaml
$ cat ./docs/examples/cluster-snapshot/azure/config.yaml

snapshotter:
  Storage:
    azure:
      container: bucket-for-snapshot
      prefix: minikube
    storageSecretName: azure-secret
  sanitize: true
  schedule: '@every 6h'
```

| Key                                     | Description                                                                     |
|-----------------------------------------|---------------------------------------------------------------------------------|
| `snapshotter.storage.storageSecretName` | `Required`. Name of storage secret                                              |
| `snapshotter.storage.azure.container`   | `Required`. Name of Azure container                                             |
| `snapshotter.storage.azure.prefix`      | `Optional`. Path prefix into bucket where snapshot will be stored               |
| `snapshotter.sanitize`                  | `Optional`. If set to `true`, various auto generated ObjectMeta and Spec fields are cleaned up before storing snapshots |
| `snapshotter.schedule`                  | `Required`. [Cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) specifying the schedule for snapshot operations. |


Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/cluster-snapshot/azure/config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  config.yaml: c25hcHNob3R0ZXI6CiAgU3RvcmFnZToKICAgIGF6dXJlOgogICAgICBjb250YWluZXI6IGJ1Y2tldC1mb3Itc25hcHNob3QKICAgICAgcHJlZml4OiBtaW5pa3ViZQogICAgc3RvcmFnZVNlY3JldE5hbWU6IGF6dXJlLXNlY3JldAogIHNhbml0aXplOiB0cnVlCiAgc2NoZWR1bGU6ICdAZXZlcnkgNmgn
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T06:01:42Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "7555"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: e5b8f78f-71c7-11e7-a5ec-0800273df5f2
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/install.md). Once the operator pod is running, check your container from Azure portal. You should see the data from initial snapshot operation.


## OpenStack Swift
Kubed supports OpenStack Swift as snapshot storage backend. To configure this backend, create a Secret with the following secret keys:

| Key                      | Description                                                |
|--------------------------|------------------------------------------------------------|
| `ST_AUTH`                | For keystone v1 authentication                             |
| `ST_USER`                | For keystone v1 authentication                             |
| `ST_KEY`                 | For keystone v1 authentication                             |
| `OS_AUTH_URL`            | For keystone v2 authentication                             |
| `OS_REGION_NAME`         | For keystone v2 authentication                             |
| `OS_USERNAME`            | For keystone v2 authentication                             |
| `OS_PASSWORD`            | For keystone v2 authentication                             |
| `OS_TENANT_ID`           | For keystone v2 authentication                             |
| `OS_TENANT_NAME`         | For keystone v2 authentication                             |
| `OS_AUTH_URL`            | For keystone v3 authentication                             |
| `OS_REGION_NAME`         | For keystone v3 authentication                             |
| `OS_USERNAME`            | For keystone v3 authentication                             |
| `OS_PASSWORD`            | For keystone v3 authentication                             |
| `OS_USER_DOMAIN_NAME`    | For keystone v3 authentication                             |
| `OS_PROJECT_NAME`        | For keystone v3 authentication                             |
| `OS_PROJECT_DOMAIN_NAME` | For keystone v3 authentication                             |
| `OS_STORAGE_URL`         | For authentication based on tokens                         |
| `OS_AUTH_TOKEN`          | For authentication based on tokens                         |


```console
$ echo -n '<your-auth-url>' > OS_AUTH_URL
$ echo -n '<your-tenant-id>' > OS_TENANT_ID
$ echo -n '<your-tenant-name>' > OS_TENANT_NAME
$ echo -n '<your-username>' > OS_USERNAME
$ echo -n '<your-password>' > OS_PASSWORD
$ echo -n '<your-region>' > OS_REGION_NAME
$ kubectl create secret generic swift-secret -n kube-system \
    --from-file=./OS_AUTH_URL \
    --from-file=./OS_TENANT_ID \
    --from-file=./OS_TENANT_NAME \
    --from-file=./OS_USERNAME \
    --from-file=./OS_PASSWORD \
    --from-file=./OS_REGION_NAME
secret "swift-secret" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret swift-secret app=kubed -n kube-system
secret "swift-secret" labeled
```

```yaml
$ kubectl get secret swift-secret -n kube-system -o yaml

apiVersion: v1
data:
  OS_AUTH_URL: PHlvdXItYXV0aC11cmw+
  OS_PASSWORD: PHlvdXItcGFzc3dvcmQ+
  OS_REGION_NAME: PHlvdXItcmVnaW9uPg==
  OS_TENANT_ID: PHlvdXItdGVuYW50LWlkPg==
  OS_TENANT_NAME: PHlvdXItdGVuYW50LW5hbWU+
  OS_USERNAME: PHlvdXItdXNlcm5hbWU+
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T06:23:22Z
  labels:
    app: kubed
  name: swift-secret
  namespace: kube-system
  resourceVersion: "9134"
  selfLink: /api/v1/namespaces/kube-system/secrets/swift-secret
  uid: ec735b2d-71ca-11e7-a5ec-0800273df5f2
type: Opaque
```

Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/cluster-snapshot/swift/kubed-config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  kubed-config.yaml: YXBpVmVyc2lvbjogdjEKa2luZDogQ29uZmlnTWFwCm1ldGFkYXRhOgogIG5hbWU6IGt1YmVkLWNvbmZpZwogIG5hbWVzcGFjZToga3ViZS1zeXN0ZW0KICBsYWJlbHM6CiAgICBhcHA6IGt1YmVkCmRhdGE6CiAgY29uZmlnLnlhbWw6IHwKICAgIHNuYXBzaG90dGVyOgogICAgICBTdG9yYWdlOgogICAgICAgIHN3aWZ0OgogICAgICAgICAgY29udGFpbmVyOiBidWNrZXQtZm9yLXNuYXBzaG90CiAgICAgICAgICBwcmVmaXg6IG1pbmlrdWJlCiAgICAgICAgc3RvcmFnZVNlY3JldE5hbWU6IHNuYXAtc2VjcmV0CiAgICAgIHNhbml0aXplOiB0cnVlCiAgICAgIHNjaGVkdWxlOiAnQGV2ZXJ5IDZoJwo=
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T06:25:54Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "9303"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 4777f28b-71cb-11e7-a5ec-0800273df5f2
type: Opaque
```

Now, let's take a look at the cluster config. Here,

| Key                                     | Description                                                                     |
|-----------------------------------------|---------------------------------------------------------------------------------|
| `snapshotter.storage.storageSecretName` | `Required`. Name of storage secret                                              |
| `snapshotter.storage.swift.container`   | `Required`. Name of OpenStack Swift container                                   |
| `snapshotter.storage.swift.prefix`      | `Optional`. Path prefix into bucket where snapshot will be stored               |
| `snapshotter.sanitize`                  | `Optional`. If set to `true`, various auto generated ObjectMeta and Spec fields are cleaned up from snapshots |
| `snapshotter.schedule`                  | `Required`. [Cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) specifying the schedule for snapshot operations. |

Now, deploy Kubed operator in your cluster following the steps [here](/docs/install.md). Once the operator pod is running, check your container. You should see the data from initial snapshot operation.


## Local Backend
`Local` backend refers to a local path inside snapshot job container. Any Kubernetes supported [persistent volume](https://kubernetes.io/docs/concepts/storage/volumes/) can be used here. Some examples are: `emptyDir` for testing, NFS, Ceph, GlusterFS, etc.
To configure this backend, no secret is needed. Following parameters are available for `Local` backend.

| Parameter                 | Description                                                                             |
|---------------------------|-----------------------------------------------------------------------------------------|
| `spec.databaseName`       | `Required`. Name of database                                                            |
| `spec.local.path`         | `Required`. Path where this volume will be mounted in the job container. Example: /repo |
| `spec.local.volumeSource` | `Required`. Any Kubernetes [volume](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes) |
| `spec.resources`          | `Optional`. Compute resources required by Jobs used to take snapshot or initialize databases from snapshot.  To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/). |

```console
$ kubectl create -f ./docs/examples/snapshot/local/local-snapshot.yaml
snapshot "local-snapshot" created
```

```yaml
$ kubectl get snapshot local-snapshot -o yaml
apiVersion: kubed.com/v1alpha1
kind: Snapshot
metadata:
  creationTimestamp: 2017-06-28T12:14:48Z
  name: local-snapshot
  namespace: default
  resourceVersion: "2000"
  selfLink: /apis/kubed.com/v1alpha1/namespaces/kube-system/snapshots/local-snapshot
  uid: 617e3487-5bfb-11e7-bb52-08002711f4aa
  labels:
    kubed.com/kind: Postgres
spec:
  databaseName: postgres-db
  local:
    path: /repo
    volumeSource:
      emptyDir: {}
  resources:
    requests:
      memory: "64Mi"
      cpu: "250m"
    limits:
      memory: "128Mi"
      cpu: "500m"
```

## instant backup




## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, uninstall Kubed operator. Please follow the steps [here](/docs/uninstall.md) to uninstall Kubed operator.




## Next Steps
- Learn about the details of Postgres tpr [here](/docs/concepts/postgres.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Thinking about monitoring your database? Kubed works [out-of-the-box with Prometheus](/docs/tutorials/monitoring.md).
- Learn how to use Kubed in a [RBAC](/docs/tutorials/rbac.md) enabled cluster.
- Wondering what features are coming next? Please visit [here](/ROADMAP.md). 
- Want to hack on Kubed? Check our [contribution guidelines](/CONTRIBUTING.md).


