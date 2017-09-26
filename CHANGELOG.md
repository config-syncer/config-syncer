# Change Log

## [0.3.0](https://github.com/appscode/kubed/tree/0.3.0) (2017-09-26)
Kubed 0.3.0 supports auth for Elasticsearch janitors and annotates copied configmaps & secrets.

__Changelog__
- Support auth for Elasticsearch janitor [\#64](https://github.com/appscode/kubed/issues/64)
- Install kubed as a critical addon [\#129](https://github.com/appscode/kubed/pull/129)
- Update chart to match recent convention [\#123](https://github.com/appscode/kubed/pull/123)
- Use correct service account for RBAC installer [\#122](https://github.com/appscode/kubed/pull/122)
- Forward CSR approved/denied events [\#117](https://github.com/appscode/kubed/pull/117)
- Use kutil package for utils [\#116](https://github.com/appscode/kubed/pull/116)
- Annotate copied configmaps & secrets with kubed.appscode.com/origin [\#115](https://github.com/appscode/kubed/pull/115)
- Use client-go 4.0.0 [\#114](https://github.com/appscode/kubed/pull/114)
- Fix config object. [\#105](https://github.com/appscode/kubed/pull/105)


## [0.2.0](https://github.com/appscode/kubed/tree/0.2.0) (2017-08-04)
Kubed 0.2.0 can send push notifications.

__Notable Features__
 - Send [push notification](/docs/tutorials/notifiers.md#pushovernet) via [pushover.net](https://pushover.net/) Thanks [Sean Johnson](https://github.com/pirogoeth) !
 - Add `clusterName` in config. This cluster name will be prefixed to any notification sent via Email/SMS/Chat/Push so that you can identify the source easily.


## [0.1.0](https://github.com/appscode/kubed/tree/0.1.0) (2017-08-01)
First public release of Kubed.

__Notable Features__
 - Takes periodic [snapshot](/docs/tutorials/cluster-snapshot.md) of a Kubernetes cluster objects.
 - Provides a [recycle bin](/docs/tutorials/recycle-bin.md) for deleted and/or updated Kubernetes objects.
 - Keep [ConfigMaps and Secrets synchronized across Namespaces](/docs/tutorials/config-syncer.md).
 - [Forward cluster events](/docs/tutorials/event-forwarder.md) to various destinations.
 - Setup [janitors](/docs/tutorials/janitors.md) for Elasticsearch and InfluxDB.
 - [Send notifications](/docs/tutorials/notifiers.md) via Email, SMS or Chat.
 - Includes a built-in [search engine](/docs/tutorials/apiserver.md) for your cluster objects using [bleve](https://github.com/blevesearch/bleve).
