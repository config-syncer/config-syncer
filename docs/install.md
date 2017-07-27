> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Installation Guide










```yaml
apiServer:
  address: :8080
  enableReverseIndex: true
  enableSearchIndex: true
enableConfigSyncer: true
eventForwarder:
  ingressAdded:
    handle: true
  nodeAdded:
    handle: true
  receiver:
    notifier: mailgun
    to:
    - ops@example.com
  storageAdded:
    handle: true
  warningEvents:
    handle: true
    namespaces:
    - kube-system
janitors:
- elasticsearch:
    endpoint: http://elasticsearch-logging.kube-system:9200
    logIndexPrefix: logstash-
  kind: Elasticsearch
  ttl: 2160h0m0s
- influxdb:
    endpoint: https://monitoring-influxdb.kube-system:8086
  kind: InfluxDB
  ttl: 2160h0m0s
notifierSecretName: kubed-notifier
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  receiver:
    notifier: mailgun
    to:
    - ops@example.com
  ttl: 168h0m0s
snapshotter:
  Storage:
    gcs:
      bucket: restic
      prefix: minikube
    storageSecretName: snap-secret
  sanitize: true
  schedule: '@every 6h'
```


























## Using YAML
Kubed can be installed using YAML files includes in the [/hack/deploy](/hack/deploy) folder.

```console
# Install without RBAC roles
$ curl https://raw.githubusercontent.com/appscode/kubed/0.1.0/hack/deploy/without-rbac.yaml \
  | kubectl apply -f -


# Install with RBAC roles
$ curl https://raw.githubusercontent.com/appscode/kubed/0.1.0/hack/deploy/with-rbac.yaml \
  | kubectl apply -f -
```

## Using Helm
Kubed can be installed via [Helm](https://helm.sh/) using the [chart](/chart/kubed) included in this repository. To install the chart with the release name `my-release`:
```bash
$ helm install chart/kubed --name my-release
```
To see the detailed configuration options, visit [here](/chart/kubed/README.md).


## Verify installation
To check if Kubed operator pods have started, run the following command:
```console
$ kubectl get pods --all-namespaces -l app=kubed --watch
```

Once the operator pods are running, you can cancel the above command by typing `Ctrl+C`.

Now, you are ready to [start managing your cluster](/docs/tutorials/README.md) using Kubed.
