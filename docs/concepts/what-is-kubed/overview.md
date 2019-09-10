---
title: Kubed Overview
description: Kubed Overview
menu:
  product_kubed_0.11.0:
    identifier: overview-concepts
    name: Overview
    parent: what-is-kubed
    weight: 10
product_name: kubed
menu_name: product_kubed_0.11.0
section_menu_id: concepts
---

# Kubed

Kubed (pronounced Cube-Dee) by AppsCode is a Kubernetes Cluster Operator Daemon. Kubed can do the following things for you:

 - Kubed can protect your Kubernetes cluster from [various disasters scenarios](/docs/guides/disaster-recovery/).
 - Kubed can keep [ConfigMaps and Secrets synchronized across Namespaces](/docs/guides/config-syncer/).
 - Kubed can [forward cluster events](/docs/guides/cluster-events/) to various destinations.
 - Kubed can setup [janitors](/docs/guides/janitors.md) for Elasticsearch and InfluxDB.
 - Kubed can [send notifications](/docs/guides/cluster-events/notifiers.md) via Email, SMS or Chat.
 - Kubed has a built-in [search engine](/docs/guides/apiserver.md) for your cluster objects using [bleve](https://github.com/blevesearch/bleve).
