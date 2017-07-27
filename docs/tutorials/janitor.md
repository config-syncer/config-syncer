> New to Kubed? Please start [here](/docs/tutorials/README.md).


# Using Janitors
Kubed includes janitors for cleaning up Elasticseach and InfluxDB databases. 
Kubed taking periodic snapshot of a Kubernetes cluster objects. The snapshot data can be stored in various cloud providers, eg, [Amazon S3](#aws-s3), [Google Cloud Storage](#google-cloud-storage-gcs), [Microsoft Azure](#microsoft-azure-storage), [OpenStack Swift](#openstack-swift) and any [locally mounted volumes](#local-backend) like NFS, GlusterFS, etc. Kubed uses Kubernetes discovery api to find all available resources in a cluster and stores them in a file matching the `selfLink` URL for an object. This tutorial will show you how to use Kubed to take periodic snapshots of a Kubernetes cluster objects.


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

## Disable Snapshotter
If you would like to disable this feature, remove the `snapshotter` portion of your Kubed cluster config. Then update the `kubed-config` Secret and restart Kubed operator pod(s).


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, uninstall Kubed operator following the steps [here](/docs/uninstall.md).

## Next Steps
