[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/kubed)](https://goreportcard.com/report/github.com/appscode/kubed)

# Kubed
Kubed by AppsCode is a Kubernetes Cluster Operator Daemon. Kubed can do the following things for you:

 - Kubed can take periodic [snapshot](/docs/tutorials/cluster-snapshot.md) of a Kubernetes cluster objects.
 - Kubed provides a [recycle bin](/docs/tutorials/recycle-bin.md) for deleted and/or updated Kubernetes objects.
 - Kubed can [synchronize ConfigMap/Secret across Namespaces](/docs/tutorials/config-syncer.md).
 - Kubed can [forward Cluster Events](/docs/tutorials/event-forwarder.md) to various destinations.
 - Kubed can setup [janitors](/docs/tutorials/janitors.md) for Elasticsearch and InfluxDB.
 - Kubed can [send notifications](/docs/tutorials/notifiers.md) via Email, SMS or Chat.
 - Kubed has a built-in [search engine](/docs/tutorials/apiserver.md) for your cluster objects using [bleve](https://github.com/blevesearch/bleve).

## Supported Versions
Kubernetes 1.5+

## Installation
To install Kubed, please follow the guide [here](/docs/install.md).

## Using Kubed
Want to learn how to use Kubed? Please start [here](/docs/tutorials/README.md).

## Contribution guidelines
Want to help improve Kubed? Please start [here](/CONTRIBUTING.md).

## Project Status
Wondering what features are coming next? Please visit [here](/ROADMAP.md).

---

**The kubed operator collects anonymous usage statistics to help us learn how the software is being used and how we can improve it. To disable stats collection, run the operator with the flag** `--analytics=false`.

---

## Support
If you have any questions, you can reach out to us.
* [Slack](https://slack.appscode.com)
* [Twitter](https://twitter.com/AppsCodeHQ)
* [Website](https://appscode.com)
