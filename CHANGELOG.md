# Change Log

## [Unreleased](https://github.com/appscode/kubed/tree/HEAD)

[Full Changelog](https://github.com/appscode/kubed/compare/4.0.0-alpha.0...HEAD)

**Implemented enhancements:**

- Support TLS for elasticsearch connection [\#126](https://github.com/appscode/kubed/pull/126) ([aerokite](https://github.com/aerokite))

**Fixed bugs:**

- Installing kubed fails due to missing service account [\#121](https://github.com/appscode/kubed/issues/121)
- Cleanup search index when a namespace is deleted. [\#109](https://github.com/appscode/kubed/issues/109)

**Closed issues:**

- Vault Integration [\#119](https://github.com/appscode/kubed/issues/119)
- Support auth for Elasticsearch janitor [\#64](https://github.com/appscode/kubed/issues/64)

**Merged pull requests:**

- Revendor for generator clients. [\#124](https://github.com/appscode/kubed/pull/124) ([tamalsaha](https://github.com/tamalsaha))
- Update chart to match recent convention [\#123](https://github.com/appscode/kubed/pull/123) ([tamalsaha](https://github.com/tamalsaha))
- Use correct service account for RBAC installer [\#122](https://github.com/appscode/kubed/pull/122) ([tamalsaha](https://github.com/tamalsaha))
- Fix command in Developer-guide doc [\#120](https://github.com/appscode/kubed/pull/120) ([the-redback](https://github.com/the-redback))

## [4.0.0-alpha.0](https://github.com/appscode/kubed/tree/4.0.0-alpha.0) (2017-09-05)
[Full Changelog](https://github.com/appscode/kubed/compare/0.2.0...4.0.0-alpha.0)

**Implemented enhancements:**

- Annotate replicated objects indicating they are a replica and the source [\#112](https://github.com/appscode/kubed/issues/112)

**Closed issues:**

- Notify about new CSR requests [\#73](https://github.com/appscode/kubed/issues/73)
- Support CRD [\#53](https://github.com/appscode/kubed/issues/53)

**Merged pull requests:**

- Forward CSR approved/denied events [\#117](https://github.com/appscode/kubed/pull/117) ([tamalsaha](https://github.com/tamalsaha))
- Use kutil package for utils [\#116](https://github.com/appscode/kubed/pull/116) ([tamalsaha](https://github.com/tamalsaha))
- Annotate copied configmaps & secrets with kubed.appscode.com/origin [\#115](https://github.com/appscode/kubed/pull/115) ([tamalsaha](https://github.com/tamalsaha))
- Use client-go 4.0.0 [\#114](https://github.com/appscode/kubed/pull/114) ([tamalsaha](https://github.com/tamalsaha))
- Fix config object. [\#105](https://github.com/appscode/kubed/pull/105) ([tamalsaha](https://github.com/tamalsaha))

## [0.2.0](https://github.com/appscode/kubed/tree/0.2.0) (2017-08-04)
Kubed 0.2.0 can send push notifications. To install, please visit [here](/docs/install.md). If you are existing user, please update the image tag.

__Notable Features__
 - Send [push notification](/docs/tutorials/notifiers.md#pushovernet) via [pushover.net](https://pushover.net/) Thanks [Sean Johnson](https://github.com/pirogoeth) !
 - Add `clusterName` in config. This cluster name will be prefixed to any notification sent via Email/SMS/Chat/Push so that you can identify the source easily.


## [0.1.0](https://github.com/appscode/kubed/tree/0.1.0) (2017-08-01)
First public release of Kubed. To install, please visit [here](/docs/install.md).

__Notable Features__
 - Takes periodic [snapshot](/docs/tutorials/cluster-snapshot.md) of a Kubernetes cluster objects.
 - Provides a [recycle bin](/docs/tutorials/recycle-bin.md) for deleted and/or updated Kubernetes objects.
 - Keep [ConfigMaps and Secrets synchronized across Namespaces](/docs/tutorials/config-syncer.md).
 - [Forward cluster events](/docs/tutorials/event-forwarder.md) to various destinations.
 - Setup [janitors](/docs/tutorials/janitors.md) for Elasticsearch and InfluxDB.
 - [Send notifications](/docs/tutorials/notifiers.md) via Email, SMS or Chat.
 - Includes a built-in [search engine](/docs/tutorials/apiserver.md) for your cluster objects using [bleve](https://github.com/blevesearch/bleve).
