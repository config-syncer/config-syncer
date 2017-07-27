# Using Kubed

> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Kubed API Server
Kubed includes an api server. It has 2 categories of endpoints:
 - Search objects
 - Reverse Lookup





> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Cluster Snapshots
Kubed supports taking periodic snapshot of a Kubernetes cluster objects. The snapshot data can be stored in various cloud providers, eg, [Amazon S3](#aws-s3), [Google Cloud Storage](#google-cloud-storage-gcs), [Microsoft Azure](#microsoft-azure-storage), [OpenStack Swift](#openstack-swift) and any [locally mounted volumes](#local-backend) like NFS, GlusterFS, etc. Kubed uses Kubernetes discovery api to find all available resources in a cluster and stores them in a file matching the `selfLink` URL for an object. This tutorial will show you how to use Kubed to take periodic snapshots of a Kubernetes cluster objects.



> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Synchronize Configuration across Namespaces
Sometimes you have some configuration that you want to synchronize across all Kubernetes namespaces. Kubed can do that for you. If a ConfigMap or Secret has the annotation __`kubed.appscode.com/sync: true`__, Kubed will create a similar ConfigMap / Secret in all existing namespaces. Kubed will also create this ConfigMap/Secret, when you create a new namespace. If the data in the source ConfigMap/Secret is updated, all the copies will be updated. Either delete the source ConfigMap/Secret or remove the annotation from the source ConfigMap/Secret to purge the copies. If the namespace with the source ConfigMap/Secret is deleted, the copies are left intact.


> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Forward Cluster Events
Kubed can send notifications via Email, SMS or Chat for various cluster events. This tutorial will show you how to use Kubed to setup an event forwarder.

> New to Kubed? Please start [here](/docs/tutorials/README.md).


# Using Janitors
Kubed includes janitors for Elasticsearch and InfluxDB. These janitors can delete data older than a configures TTL. Kubernetes supports storing cluster logs in Elasticsearch and cluster metrics in InfluxDB. You use these janitors to clean up old data from Elasticsearch and InfluxDB before those fill up your node disks.
> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Monitoring Kubed

KubeDB has native support for monitoring via Prometheus. KubeDB operator exposes Prometheus native monitoring data via `/metrics` endpoint on `:56790` port. You can setup a [CoreOS Prometheus ServiceMonitor](https://github.com/coreos/prometheus-operator) using `kubed-operator` service. To change the port, use `--web.address` flag on Kubed operator.
> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Supported Notifiers
Kubed can send notifications via Email, SMS or Chat for various operations using [appscode/go-notify](https://github.com/appscode/go-notify) library. To connect to these services, you need to create a Secret with the appropriate keys. Then pass the secret name to Kubed by setting `notifierSecretName` field in Kubed cluster config.
> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Kubernetes Recycle Bin
Kubed provides a recycle bin for deleted and/or updated Kubernetes objects. Once activated, any deleted and/or updated object is stored in YAML format in folder mounted inside Kubed pod. This tutorial will show you how to use Kubed to setup a recycle bin for Kubernetes cluster objects.

> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Kubernetes Recycle Bin
Kubed provides a recycle bin for deleted and/or updated Kubernetes objects. Once activated, any deleted and/or updated object is stored in YAML format in folder mounted inside Kubed pod. This tutorial will show you how to use Kubed to setup a recycle bin for Kubernetes cluster objects.


