---
title: Janitors
description: Janitors
menu:
  product_kubed_0.3.1:
    identifier: tutorials-janitors
    name: janitors
    parent: tutorials
    weight: 40
product_name: kubed
menu_name: product_kubed_0.3.1
section_menu_id: tutorials
---

> New to Kubed? Please start [here](/docs/guides/README.md).


# Using Janitors
Kubed includes janitors for Elasticsearch and InfluxDB. These janitors can delete data older than a configured TTL. Kubernetes supports storing cluster logs in Elasticsearch and cluster metrics in InfluxDB. You use these janitors to clean up old data from Elasticsearch and InfluxDB before those fill up your node disks.

---

Please check your janitor configuration on test clusters before using in production. You have been forewarned! We [welcome contribution](https://github.com/appscode/kubed/issues/60) to support dryRun options for janitors.

---

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

## Deploy Kubed
To enable janitors, you need a cluster config like below.

```yaml
$ cat ./docs/examples/janitors/config.yaml

janitors:
- kind: Elasticsearch
  ttl: 2160h
  elasticsearch:
    endpoint: http://elasticsearch-logging.kube-system:9200
    logIndexPrefix: logstash-
- kind: InfluxDB
  ttl: 2160h
  influxdb:
    endpoint: https://monitoring-influxdb.kube-system:8086
```

| Key                            | Description                                                             |
|--------------------------------|-------------------------------------------------------------------------|
| `kind`                         | `Required`. Set this to either `Elasticseach` or `InfluxDB`.            |
| `ttl`                          | `Required`. Time-to-live for data (eg, 5h30m30s).                       |
| `elasticsearch.endpoint`       | Required for kind `Elasticsearch`. URL of Elasticsearch cluster.        |
| `elasticsearch.logIndexPrefix` | Required for kind `Elasticsearch`. Prefix of log index.                 |
| `elasticsearch.secretName`     | Optional for kind `Elasticsearch`. Name of secret used to pass Elasticsearch authentication . |
| `influxdb.endpoint`            | Required for kind `InfluxDB`. URL of [InfluxDB server](https://github.com/kubernetes/heapster/blob/master/docs/sink-configuration.md#influxdb). |
| `influxdb.username`            | Optional for kind `InfluxDB`. InfluxDB username (default: root)         |
| `influxdb.password`            | Optional for kind `InfluxDB`. InfluxDB password (default: root)         |


Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/janitors/config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  config.yaml: amFuaXRvcnM6Ci0gZWxhc3RpY3NlYXJjaDoKICAgIGVuZHBvaW50OiBodHRwOi8vZWxhc3RpY3NlYXJjaC1sb2dnaW5nLmt1YmUtc3lzdGVtOjkyMDAKICAgIGxvZ0luZGV4UHJlZml4OiBsb2dzdGFzaC0KICBraW5kOiBFbGFzdGljc2VhcmNoCiAgdHRsOiAyMTYwaAotIGluZmx1eGRiOgogICAgZW5kcG9pbnQ6IGh0dHBzOi8vbW9uaXRvcmluZy1pbmZsdXhkYi5rdWJlLXN5c3RlbTo4MDg2CiAga2luZDogSW5mbHV4REIKICB0dGw6IDIxNjBoCg==
kind: Secret
metadata:
  creationTimestamp: 2017-07-27T07:43:32Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "27760846"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 4a2eb872-729f-11e7-8b69-12f236046fba
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/setup/install.md). Once the operator pod is running, janitor operations are applied within one hour.


## Janitor Authentication
The following keys are supported for Secret passed via `elasticsearch.secretName`.

| Key                    | Description                                                                 |
-------------------------|-----------------------------------------------------------------------------|
| `CA_CERT_DATA`         | `Optional` PEM encoded CA certificate used to connect to Elasticsearch.     |
| `CLIENT_CERT_DATA`     | `Optional` PEM encoded Client certificate used to connect to Elasticsearch. |
| `CLIENT_KEY_DATA`      | `Optional` PEM encoded Client private key used to connect to Elasticsearch. |
| `INSECURE_SKIP_VERIFY` | `Optional` If set to `true`, skip certificate verification.                 |


## Disable Janitors
If you would like to disable this feature, remove the `janitors` portion of your Kubed cluster config. Then update the `kubed-config` Secret and restart Kubed operator pod(s).


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, uninstall Kubed operator following the steps [here](/docs/setup/uninstall.md).


## Next Steps
 - Learn how to use Kubed to take periodic snapshots of a Kubernetes cluster [here](/docs/guides/cluster-snapshot.md).
 - To setup a recycle bin for deleted and/or updated Kubernetes objects, please visit [here](/docs/guides/recycle-bin.md).
 - Need to keep some configuration synchronized across namespaces? Try [Kubed config syncer](/docs/guides/config-syncer.md).
 - Want to keep an eye on your cluster with automated notifications? Setup Kubed [event forwarder](/docs/guides/event-forwarder.md).
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/ROADMAP.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
