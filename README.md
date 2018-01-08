[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/kubed)](https://goreportcard.com/report/github.com/appscode/kubed)

# Kubed
Kubed (pronounced Cube-Dee) by AppsCode is a Kubernetes Cluster Operator Daemon. Kubed can do the following things for you:

 - Kubed can protect your Kubernetes cluster from [various disasters scenarios](https://appscode.com/products/kubed/0.4.0/guides/disaster-recovery/).
 - Kubed can keep [ConfigMaps and Secrets synchronized across namespaces and/or clusters](https://appscode.com/products/kubed/0.4.0/guides/config-syncer/).
 - Kubed can [forward cluster events](https://appscode.com/products/kubed/0.4.0/guides/cluster-events/) to various destinations.
 - Kubed can setup [janitors](https://appscode.com/products/kubed/0.4.0/guides/janitors/) for Elasticsearch and InfluxDB.
 - Kubed can [send notifications](https://appscode.com/products/kubed/0.4.0/guides/cluster-events/notifiers/) via Email, SMS or Chat.
 - Kubed has a built-in [search engine](https://appscode.com/products/kubed/0.4.0/guides/apiserver/) for your cluster objects using [bleve](https://github.com/blevesearch/bleve).


## Supported Versions
Please pick a version of Kubed that matches your Kubernetes installation.

| Kubed Version                                                 | Docs                                                            | Kubernetes Version |
|---------------------------------------------------------------|-----------------------------------------------------------------|--------------------|
| [0.4.0](https://github.com/appscode/kubed/releases/tag/0.4.0) | [User Guide](https://appscode.com/products/kubed/0.4.0/)        | 1.7.x+             |
| [0.2.0](https://github.com/appscode/kubed/releases/tag/0.2.0) | [User Guide](https://github.com/appscode/kubed/tree/0.2.0/docs) | 1.5.x - 1.7.x      |

## Installation
To install Kubed, please follow the guide [here](https://appscode.com/products/kubed/0.4.0/setup/install/).

## Using Kubed
Want to learn how to use Kubed? Please start [here](https://appscode.com/products/kubed/0.4.0/).

## Contribution guidelines
Want to help improve Kubed? Please start [here](https://appscode.com/products/kubed/0.4.0/welcome/contributing/).

---

**Kubed binaries collects anonymous usage statistics to help us learn how the software is being used and how we can improve it. To disable stats collection, run the operator with the flag** `--analytics=false`.

---

## Support
We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [AppsCode Slack team](https://appscode.slack.com/messages/C6HSHCKBL/details/) channel `#kubed`. To sign up, use our [Slack inviter](https://slack.appscode.com/).

If you have found a bug with Searchlight or want to request for new features, please [file an issue](https://github.com/appscode/kubed/issues/new).
