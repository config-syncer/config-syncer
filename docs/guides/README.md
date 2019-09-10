---
title: Overview
description: Overview of guides
menu:
  product_kubed_v0.11.0:
    identifier: guides-overview
    name: Overview
    parent: guides
    weight: -1
product_name: kubed
menu_name: product_kubed_v0.11.0
section_menu_id: guides
url: /products/kubed/v0.11.0/guides/
aliases:
  - /products/kubed/v0.11.0/guides/README/
---

# Guides

This section contains guides on how to use Kubed. Please visit the links below to learn more:

- Disaster Recovery
  - [Cluster Snapshots](/docs/guides/disaster-recovery/cluster-snapshot.md): This tutorial will show you how to use Kubed to take periodic snapshots of a Kubernetes cluster objects.
  - [Kubernetes Recycle Bin](/docs/guides/disaster-recovery/recycle-bin.md): Kubed provides a recycle bin for deleted and/or updated Kubernetes objects. This tutorial will show you how to use Kubed to setup a recycle bin for Kubernetes cluster objects.
  - [Backup & Restore Persistent Volumes](/docs/guides/disaster-recovery/stash.md). Use [Stash](https://appscode.com/products/stash) to backup & restore Persistent Volumes.
- Configuration Syncer
  - [Synchronize Configuration across Namespaces](/docs/guides/config-syncer/intra-cluster.md): This tutorial will show you how Kubed can sync ConfigMaps/Secrets across Kubernetes namespaces.
  - [Synchronize Configuration across Clusters](/docs/guides/config-syncer/inter-cluster.md): This tutorial will show you how Kubed can sync ConfigMaps/Secrets across Kubernetes cluster.
- Cluster Events
  - [Forward Cluster Events](/docs/guides/cluster-events/event-forwarder.md): This tutorial will show you how to use Kubed to send notifications via Email, SMS or Chat for various cluster events.
  - [Supported Notifiers](/docs/guides/cluster-events/notifiers.md): This article documents how to configure Kubed to send notifications via Email, SMS or Chat
- [Using Janitors](/docs/guides/janitors.md): Kubernetes supports storing cluster logs in Elasticsearch and cluster metrics in InfluxDB. This tutorial will show you how to use kubed janitors to delete old data from Elasticsearch and InfluxDB.
- [Kubed API Server](/docs/guides/apiserver.md): This article documents various aspects of Kubed api server.
- [Monitoring Kubed](/docs/guides/monitoring.md): This article described the various metrics exposed by Kubed operator.
