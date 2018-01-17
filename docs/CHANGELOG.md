---
title: Changelog | Kubed
description: Changelog
menu:
  product_kubed_0.5.0:
    identifier: changelog-kubed
    name: Changelog
    parent: welcome
    weight: 10
product_name: kubed
menu_name: product_kubed_0.5.0
section_menu_id: welcome
url: /products/kubed/0.5.0/welcome/changelog/
aliases:
  - /products/kubed/0.5.0/CHANGELOG/
---

# Change Log

## [0.5.0](https://github.com/appscode/kubed/releases/tag/0.5.0) (2018-01-16)
Kubed 0.5.0 can send notifications to Telegram and fixes various issues related to janitors and cluster backup.

__Changelog__

- Document valid time units for janitor TTL [\#188](https://github.com/appscode/kubed/pull/188)
- Reset shard duration for influx janitor [\#187](https://github.com/appscode/kubed/pull/187)
- Set min retention policy for kubed influx janitor [\#186](https://github.com/appscode/kubed/pull/186)
- Log influx janitor result [\#185](https://github.com/appscode/kubed/pull/185)
- Update github.com/influxdata/influxdb to v1.3.3 [\#184](https://github.com/appscode/kubed/pull/184)
- Increase burst and qps for kube client [\#183](https://github.com/appscode/kubed/pull/183)
- Update github.com/influxdata/influxdb to v1.1.1 [\#182](https://github.com/appscode/kubed/pull/182)
- Update Elasticsearch client to olivere/elastic.v5 [\#181](https://github.com/appscode/kubed/pull/181)
- Add Telegram as notifier [\#180](https://github.com/appscode/kubed/pull/180)
- Delete all older indices prior to a date [\#179](https://github.com/appscode/kubed/pull/179)
- Ensure bad backups are not used to overwrite last good backup [\#178](https://github.com/appscode/kubed/pull/178)


## [0.4.0](https://github.com/appscode/kubed/releases/tag/0.4.0) (2018-01-05)
Kubed 0.4.0 can sync confgimaps/secrets across clusters.

__Changelog__

- Reorganize docs for hosting on product site [\#173](https://github.com/appscode/kubed/pull/173)
- Add support for new DB types [\#172](https://github.com/appscode/kubed/pull/172)
- Update docs for syncer [\#170](https://github.com/appscode/kubed/pull/170)
- Fix analytics client-id detection [\#168](https://github.com/appscode/kubed/pull/168)
- Auto detect AWS bucket region [\#166](https://github.com/appscode/kubed/pull/166)
- Support hipchat server [\#165](https://github.com/appscode/kubed/pull/165)
- Write event for syncer origin conflict [\#164](https://github.com/appscode/kubed/pull/164)
- Fix Syncer [\#163](https://github.com/appscode/kubed/pull/163)
- Remove unnecessary IsPreferredAPIResource api calls [\#162](https://github.com/appscode/kubed/pull/162)
- Sync configmap/secret to selected namespaces/contexts [\#154](https://github.com/appscode/kubed/pull/154)


## [0.3.1](https://github.com/appscode/kubed/releases/tag/0.3.1) (2017-12-21)
Kubed 0.3.1 supports AWS S3 buckets in regions other than us-east-1.

__Changelog__

- Support region for s3 backend [\#159](https://github.com/appscode/kubed/issues/159)
- Avoid listing buckets [\#141](https://github.com/appscode/kubed/issues/141)


## [0.3.0](https://github.com/appscode/kubed/releases/tag/0.3.0) (2017-09-26)
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


## [0.2.0](https://github.com/appscode/kubed/releases/tag/0.2.0) (2017-08-04)
Kubed 0.2.0 can send push notifications.

__Notable Features__

 - Send [push notification](/docs/tutorials/notifiers.md#pushovernet) via [pushover.net](https://pushover.net/) Thanks [Sean Johnson](https://github.com/pirogoeth) !
 - Add `clusterName` in config. This cluster name will be prefixed to any notification sent via Email/SMS/Chat/Push so that you can identify the source easily.


## [0.1.0](https://github.com/appscode/kubed/releases/tag/0.1.0) (2017-08-01)
First public release of Kubed.

__Notable Features__

 - Takes periodic [snapshot](/docs/tutorials/cluster-snapshot.md) of a Kubernetes cluster objects.
 - Provides a [recycle bin](/docs/tutorials/recycle-bin.md) for deleted and/or updated Kubernetes objects.
 - Keep [ConfigMaps and Secrets synchronized across Namespaces](/docs/tutorials/config-syncer.md).
 - [Forward cluster events](/docs/tutorials/event-forwarder.md) to various destinations.
 - Setup [janitors](/docs/tutorials/janitors.md) for Elasticsearch and InfluxDB.
 - [Send notifications](/docs/tutorials/notifiers.md) via Email, SMS or Chat.
 - Includes a built-in [search engine](/docs/tutorials/apiserver.md) for your cluster objects using [bleve](https://github.com/blevesearch/bleve).
