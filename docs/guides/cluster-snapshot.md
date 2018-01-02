---
title: Cluster Snapshots
description: Cluster Snapshots
menu:
  product_kubed_0.3.1:
    identifier: tutorials-cluster-snapshots
    name: Cluster Snapshots
    parent: tutorials
    weight: 15
product_name: kubed
menu_name: product_kubed_0.3.1
section_menu_id: tutorials
---

> New to Kubed? Please start [here](/docs/guides/README.md).

# Cluster Snapshots
Kubed supports taking periodic snapshot of a Kubernetes cluster objects. The snapshot data can be stored in various cloud providers, eg, [Amazon S3](#aws-s3), [Google Cloud Storage](#google-cloud-storage-gcs), [Microsoft Azure](#microsoft-azure-storage), [OpenStack Swift](#openstack-swift) and any [locally mounted volumes](#local-backend) like NFS, GlusterFS, etc. Kubed uses Kubernetes discovery api to find all available resources in a cluster and stores them in a file matching the `selfLink` URL for an object. Kubed uses [appscode/osm](https://github.com/appscode/osm) to interact with various cloud providers. This tutorial will show you how to use Kubed to take periodic snapshots of a Kubernetes cluster objects.

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
  gcs:
    bucket: bucket-for-snapshot
    prefix: minikube
  storageSecretName: gcs-secret
  sanitize: true
  schedule: '@every 6h'
```

| Key                             | Description                                                                      |
|---------------------------------|----------------------------------------------------------------------------------|
| `snapshotter.storageSecretName` | `Required`. Name of storage secret                                               |
| `snapshotter.gcs.bucket`        | `Required`. Name of GCS Bucket                                                   |
| `snapshotter.gcs.prefix`        | `Optional`. Path prefix into bucket where snapshot will be stored                |
| `snapshotter.overwrite`         | `Optional`. If set to `true`, snapshot folders are reused, otherwise a new folder is created at snapshot timestamp.     |
| `snapshotter.sanitize`          | `Optional`. If set to `true`, various auto generated ObjectMeta and Spec fields are cleaned up before storing snapshots |
| `snapshotter.schedule`          | `Required`. [Cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) specifying the schedule for snapshot operations. |


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
  config.yaml: c25hcHNob3R0ZXI6CiAgZ2NzOgogICAgYnVja2V0OiBidWNrZXQtZm9yLXNuYXBzaG90CiAgICBwcmVmaXg6IG1pbmlrdWJlCiAgc3RvcmFnZVNlY3JldE5hbWU6IGdjcy1zZWNyZXQKICBzYW5pdGl6ZTogdHJ1ZQogIHNjaGVkdWxlOiAnQGV2ZXJ5IDZoJwo=
kind: Secret
metadata:
  creationTimestamp: 2017-08-01T06:41:07Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "1011"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 65b97f23-7684-11e7-b77d-0800274b060f
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/setup/install.md). Once the operator pod is running, check your bucket from Google Cloud console. You should see the data from initial snapshot operation.

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
  s3:
    endpoint: s3.amazonaws.com
    bucket: kubedb-qa
    prefix: minikube
  storageSecretName: snap-secret
  sanitize: true
  schedule: '@every 6h'
```

| Key                             | Description                                                                     |
|---------------------------------|---------------------------------------------------------------------------------|
| `snapshotter.storageSecretName` | `Required`. Name of storage secret                                              |
| `snapshotter.s3.endpoint`       | `Required`. Endpoint of s3 like service.                                        |
| `snapshotter.s3.bucket`         | `Required`. Name of S3 Bucket                                                   |
| `snapshotter.s3.prefix`         | `Optional`. Sub directory in the bucket where snapshot will be stored           |
| `snapshotter.overwrite`         | `Optional`. If set to `true`, snapshot folders are reused, otherwise a new folder is created at snapshot timestamp.     |
| `snapshotter.sanitize`          | `Optional`. If set to `true`, various auto generated ObjectMeta and Spec fields are cleaned up before storing snapshots |
| `snapshotter.schedule`          | `Required`. [Cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) specifying the schedule for snapshot operations. |

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
  config.yaml: c25hcHNob3R0ZXI6CiAgczM6CiAgICBlbmRwb2ludDogJ3MzLmFtYXpvbmF3cy5jb20nCiAgICBidWNrZXQ6IGJ1Y2tldC1mb3Itc25hcHNob3QKICAgIHByZWZpeDogbWluaWt1YmUKICBzdG9yYWdlU2VjcmV0TmFtZTogc25hcC1zZWNyZXQKICBzYW5pdGl6ZTogdHJ1ZQogIHNjaGVkdWxlOiAnQGV2ZXJ5IDZoJw==
kind: Secret
metadata:
  creationTimestamp: 2017-08-01T06:43:39Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "1179"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: c0624c33-7684-11e7-b77d-0800274b060f
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/setup/install.md). Once the operator pod is running, check your bucket from S3 console. You should see the data from initial snapshot operation.


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
  azure:
    container: bucket-for-snapshot
    prefix: minikube
  storageSecretName: azure-secret
  sanitize: true
  schedule: '@every 6h'
```

| Key                             | Description                                                                     |
|---------------------------------|---------------------------------------------------------------------------------|
| `snapshotter.storageSecretName` | `Required`. Name of storage secret                                              |
| `snapshotter.azure.container`   | `Required`. Name of Azure container                                             |
| `snapshotter.azure.prefix`      | `Optional`. Path prefix into bucket where snapshot will be stored               |
| `snapshotter.overwrite`         | `Optional`. If set to `true`, snapshot folders are reused, otherwise a new folder is created at snapshot timestamp.     |
| `snapshotter.sanitize`          | `Optional`. If set to `true`, various auto generated ObjectMeta and Spec fields are cleaned up before storing snapshots |
| `snapshotter.schedule`          | `Required`. [Cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) specifying the schedule for snapshot operations. |


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
  config.yaml: c25hcHNob3R0ZXI6CiAgYXp1cmU6CiAgICBjb250YWluZXI6IGJ1Y2tldC1mb3Itc25hcHNob3QKICAgIHByZWZpeDogbWluaWt1YmUKICBzdG9yYWdlU2VjcmV0TmFtZTogYXp1cmUtc2VjcmV0CiAgc2FuaXRpemU6IHRydWUKICBzY2hlZHVsZTogJ0BldmVyeSA2aCc=
kind: Secret
metadata:
  creationTimestamp: 2017-08-01T06:45:36Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "1314"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 06180b7b-7685-11e7-b77d-0800274b060f
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/setup/install.md). Once the operator pod is running, check your container from Azure portal. You should see the data from initial snapshot operation.


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

Now, let's take a look at the cluster config. Here,

```yaml
$ cat ./docs/examples/cluster-snapshot/swift/config.yaml

snapshotter:
  swift:
    container: bucket-for-snapshot
    prefix: minikube
  storageSecretName: snap-secret
  sanitize: true
  schedule: '@every 6h'
```

| Key                             | Description                                                                     |
|---------------------------------|---------------------------------------------------------------------------------|
| `snapshotter.storageSecretName` | `Required`. Name of storage secret                                              |
| `snapshotter.swift.container`   | `Required`. Name of OpenStack Swift container                                   |
| `snapshotter.swift.prefix`      | `Optional`. Path prefix into bucket where snapshot will be stored               |
| `snapshotter.overwrite`         | `Optional`. If set to `true`, snapshot folders are reused, otherwise a new folder is created at snapshot timestamp.     |
| `snapshotter.sanitize`          | `Optional`. If set to `true`, various auto generated ObjectMeta and Spec fields are cleaned up before storing snapshots |
| `snapshotter.schedule`          | `Required`. [Cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) specifying the schedule for snapshot operations. |


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
  kubed-config.yaml: YXBpVmVyc2lvbjogdjEKZGF0YToKICBrdWJlZC1jb25maWcueWFtbDogWVhCcFZtVnljMmx2YmpvZ2RqRUthMmx1WkRvZ1EyOXVabWxuVFdGd0NtMWxkR0ZrWVhSaE9nb2dJRzVoYldVNklHdDFZbVZrTFdOdmJtWnBad29nSUc1aGJXVnpjR0ZqWlRvZ2EzVmlaUzF6ZVhOMFpXMEtJQ0JzWVdKbGJITTZDaUFnSUNCaGNIQTZJR3QxWW1Wa0NtUmhkR0U2Q2lBZ1kyOXVabWxuTG5saGJXdzZJSHdLSUNBZ0lITnVZWEJ6YUc5MGRHVnlPZ29nSUNBZ0lDQlRkRzl5WVdkbE9nb2dJQ0FnSUNBZ0lITjNhV1owT2dvZ0lDQWdJQ0FnSUNBZ1kyOXVkR0ZwYm1WeU9pQmlkV05yWlhRdFptOXlMWE51WVhCemFHOTBDaUFnSUNBZ0lDQWdJQ0J3Y21WbWFYZzZJRzFwYm1scmRXSmxDaUFnSUNBZ0lDQWdjM1J2Y21GblpWTmxZM0psZEU1aGJXVTZJSE51WVhBdGMyVmpjbVYwQ2lBZ0lDQWdJSE5oYm1sMGFYcGxPaUIwY25WbENpQWdJQ0FnSUhOamFHVmtkV3hsT2lBblFHVjJaWEo1SURab0p3bz0Ka2luZDogU2VjcmV0Cm1ldGFkYXRhOgogIGNyZWF0aW9uVGltZXN0YW1wOiAyMDE3LTA3LTI2VDA2OjI1OjU0WgogIGxhYmVsczoKICAgIGFwcDoga3ViZWQKICBuYW1lOiBrdWJlZC1jb25maWcKICBuYW1lc3BhY2U6IGt1YmUtc3lzdGVtCiAgcmVzb3VyY2VWZXJzaW9uOiAiOTMwMyIKICBzZWxmTGluazogL2FwaS92MS9uYW1lc3BhY2VzL2t1YmUtc3lzdGVtL3NlY3JldHMva3ViZWQtY29uZmlnCiAgdWlkOiA0Nzc3ZjI4Yi03MWNiLTExZTctYTVlYy0wODAwMjczZGY1ZjIKdHlwZTogT3BhcXVlCg==
kind: Secret
metadata:
  creationTimestamp: 2017-08-01T06:46:44Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "1386"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 2e9e3b46-7685-11e7-b77d-0800274b060f
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/setup/install.md). Once the operator pod is running, check your container. You should see the data from initial snapshot operation.


## Local Backend
`Local` backend refers to a local path inside Kubed container. When running Kubed, mount any Kubernetes supported [persistent volume](https://kubernetes.io/docs/concepts/storage/volumes/) and configure kubed to store snapshot data in that volume. Some examples are: `emptyDir` for testing, NFS, Ceph, GlusterFS, etc. Let's take a look at the cluster config. Here,

```yaml
$ cat ./docs/examples/cluster-snapshot/local/config.yaml

snapshotter:
  local:
    path: /var/data
  sanitize: true
  schedule: '@every 6h'
```

| Key                             | Description                                                                     |
|---------------------------------|---------------------------------------------------------------------------------|
| `snapshotter.local.path`        | `Optional`. Path where snapshot will be stored                                  |
| `snapshotter.overwrite`         | `Optional`. If set to `true`, snapshot folders are reused, otherwise a new folder is created at snapshot timestamp.     |
| `snapshotter.sanitize`          | `Optional`. If set to `true`, various auto generated ObjectMeta and Spec fields are cleaned up before storing snapshots |
| `snapshotter.schedule`          | `Required`. [Cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) specifying the schedule for snapshot operations. |


Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/cluster-snapshot/local/kubed-config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  kubed-config.yaml: YXBpVmVyc2lvbjogdjEKZGF0YToKICBrdWJlZC1jb25maWcueWFtbDogWVhCcFZtVnljMmx2YmpvZ2RqRUthMmx1WkRvZ1EyOXVabWxuVFdGd0NtMWxkR0ZrWVhSaE9nb2dJRzVoYldVNklHdDFZbVZrTFdOdmJtWnBad29nSUc1aGJXVnpjR0ZqWlRvZ2EzVmlaUzF6ZVhOMFpXMEtJQ0JzWVdKbGJITTZDaUFnSUNCaGNIQTZJR3QxWW1Wa0NtUmhkR0U2Q2lBZ1kyOXVabWxuTG5saGJXdzZJSHdLSUNBZ0lITnVZWEJ6YUc5MGRHVnlPZ29nSUNBZ0lDQlRkRzl5WVdkbE9nb2dJQ0FnSUNBZ0lHeHZZMkZzT2dvZ0lDQWdJQ0FnSUNBZ2NHRjBhRG9nTDNaaGNpOWtZWFJoQ2lBZ0lDQWdJQ0FnYzNSdmNtRm5aVk5sWTNKbGRFNWhiV1U2SUdkamN5MXpaV055WlhRS0lDQWdJQ0FnYzJGdWFYUnBlbVU2SUhSeWRXVUtJQ0FnSUNBZ2MyTm9aV1IxYkdVNklDZEFaWFpsY25rZ05tZ25DZz09CmtpbmQ6IFNlY3JldAptZXRhZGF0YToKICBjcmVhdGlvblRpbWVzdGFtcDogMjAxNy0wNy0yNlQwOToyNzoxMloKICBsYWJlbHM6CiAgICBhcHA6IGt1YmVkCiAgbmFtZToga3ViZWQtY29uZmlnCiAgbmFtZXNwYWNlOiBrdWJlLXN5c3RlbQogIHJlc291cmNlVmVyc2lvbjogIjIxMjQ5IgogIHNlbGZMaW5rOiAvYXBpL3YxL25hbWVzcGFjZXMva3ViZS1zeXN0ZW0vc2VjcmV0cy9rdWJlZC1jb25maWcKICB1aWQ6IDliNDU1OTRiLTcxZTQtMTFlNy1hNWVjLTA4MDAyNzNkZjVmMgp0eXBlOiBPcGFxdWUK
kind: Secret
metadata:
  creationTimestamp: 2017-08-01T06:47:43Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "1453"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 51a10e2d-7685-11e7-b77d-0800274b060f
type: Opaque
```

Now, deploy Kubed operator in your cluster. Since, `snapshotter.local.path` is set to `/var/data`, mount your local volume at that path. You can find example installation scripts below:
 - [Without RBAC](/docs/examples/cluster-snapshot/local/without-rbac.yaml)
 - [With RBAC](/docs/examples/cluster-snapshot/local/with-rbac.yaml)

Once the operator pod is running, check your container. You should see the data from initial snapshot operation.

## Instant Snapshot
To take an instant snapshot of a cluster, you can use the `snapshot` command from kubed. Download the pre-built binary from [appscode/kubed Github releases](https://github.com/appscode/kubed/releases) and put the binary to some directory in your `PATH`.

```console
$ kubed snapshot --context=minikube --backup-dir=/tmp/minikube

$ ls -l /tmp/minikube
total 24
drwxrwxr-x 3 tamal tamal  4096 Jul 26 02:42 api/
-rwxr-xr-x 1 tamal tamal 15477 Jul 26 02:42 api_resources.yaml*
drwxrwxr-x 5 tamal tamal  4096 Jul 26 02:42 apis/
```

## Disable Snapshotter
If you would like to disable this feature, remove the `snapshotter` portion of your Kubed cluster config. Then update the `kubed-config` Secret and restart Kubed operator pod(s).


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, uninstall Kubed operator following the steps [here](/docs/setup/uninstall.md).


## Next Steps
 - To setup a recycle bin for deleted and/or updated Kubernetes objects, please visit [here](/docs/guides/recycle-bin.md).
 - Need to keep some configuration synchronized across namespaces? Try [Kubed config syncer](/docs/guides/config-syncer.md).
 - Want to keep an eye on your cluster with automated notifications? Setup Kubed [event forwarder](/docs/guides/event-forwarder.md).
 - Out of disk space because of too much logs in Elasticsearch or metrics in InfluxDB? Configure [janitors](/docs/guides/janitors.md) to delete old data.
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/ROADMAP.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
