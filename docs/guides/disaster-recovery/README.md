---
title: Disaster Recovery | Kubernetes
description: Disaster Recovery for Kubernetes Clusters
menu:
  product_kubed_0.6.0-rc.0:
    identifier: readme-dr
    name: Overview
    parent: disaster-recovery
    weight: -1
product_name: kubed
menu_name: product_kubed_0.6.0-rc.0
section_menu_id: guides
url: /products/kubed/0.6.0-rc.0/guides/disaster-recovery/
aliases:
  - /products/kubed/0.6.0-rc.0/guides/disaster-recovery/README
---

# Disaster Recovery

This section contains guides on how to use Kubed to protect your Kubernetes cluster from disasters.

  - [Cluster Snapshots](/docs/guides/disaster-recovery/cluster-snapshot.md): This tutorial will show you how to use Kubed to take periodic snapshots of a Kubernetes cluster objects.
  - [Kubernetes Recycle Bin](/docs/guides/disaster-recovery/recycle-bin.md): Kubed provides a recycle bin for deleted and/or updated Kubernetes objects. This tutorial will show you how to use Kubed to setup a recycle bin for Kubernetes cluster objects.
  - [Backup & Restore Persistent Volumes](/docs/guides/disaster-recovery/stash.md). Use [Stash](https://appscode.com/products/stash) to backup & restore Persistent Volumes.
