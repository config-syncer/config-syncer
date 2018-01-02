---
title: Concepts | Kubed
description: Overview of Concepts
menu:
  product_kubed_5.0.0:
    identifier: overview-concepts
    name: Overview
    parent: concepts
    weight: 10
product_name: kubed
menu_name: product_kubed_5.0.0
section_menu_id: concepts
---

[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/kubed)](https://goreportcard.com/report/github.com/appscode/kubed)

# Kubed
Kubed (pronounced Cube-Dee) by AppsCode is a Kubernetes Cluster Operator Daemon. Kubed can do the following things for you:

 - Kubed can take periodic [snapshot](/docs/guides/cluster-snapshot.md) of a Kubernetes cluster objects.
 - Kubed provides a [recycle bin](/docs/guides/recycle-bin.md) for deleted and/or updated Kubernetes objects.
 - Kubed can keep [ConfigMaps and Secrets synchronized across Namespaces](/docs/guides/config-syncer.md).
 - Kubed can [forward cluster events](/docs/guides/event-forwarder.md) to various destinations.
 - Kubed can setup [janitors](/docs/guides/janitors.md) for Elasticsearch and InfluxDB.
 - Kubed can [send notifications](/docs/guides/notifiers.md) via Email, SMS or Chat.
 - Kubed has a built-in [search engine](/docs/guides/apiserver.md) for your cluster objects using [bleve](https://github.com/blevesearch/bleve).