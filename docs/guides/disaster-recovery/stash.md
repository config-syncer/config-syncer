---
title: Backup & Restore Persistent Volumes
description: Backup & Restore Persistent Volumes
menu:
  product_kubed_0.7.0-rc.0:
    identifier: stash-dr
    name: Backup Volumes
    parent: disaster-recovery
    weight: 20
product_name: kubed
menu_name: product_kubed_0.7.0-rc.0
section_menu_id: guides
---

> New to Kubed? Please start [here](/docs/concepts/README.md).

# Backup & Restore Persistent Volumes

[Stash](https://appscode.com/products/stash) by AppsCode is a Kubernetes operator to backup & restore Persistent Volumes. If you are running production workloads in Kubernetes, you might want to take backup of your disks. Traditional tools are too complex to setup and maintain in a dynamic compute environment like Kubernetes. `restic` is a backup program that is fast, efficient and secure with few moving parts. Stash is a CRD controller for Kubernetes built around `restic` to address these issues. Using Stash, you can backup Kubernetes volumes mounted in following types of workloads:

- Deployment
- DaemonSet
- ReplicaSet
- ReplicationController
- StatefulSet

## Features
 - Fast, secure, efficient backup of any kubernetes [volumes](https://kubernetes.io/docs/concepts/storage/volumes/).
 - Automates configuration of `restic` for periodic backup.
 - Store backed up files in various cloud storage provider, including S3, GCS, Azure, OpenStack Swift, DigitalOcean Spaces etc.
 - Restore backup easily in any cluster on any cloud provider.
 - Periodically check integrity of backed up data.
 - Take backup in offline mode.
 - Support workload initializer for faster backup.
 - Prometheus ready metrics for backup process.

## Next Steps
 - Learn how to use Kubed to take periodic snapshots of a Kubernetes cluster [here](/docs/guides/disaster-recovery/cluster-snapshot.md).
 - To setup a recycle bin for deleted and/or updated Kubernetes objects, please visit [here](/docs/guides/disaster-recovery/recycle-bin.md).
 - Need to keep configmaps/secrets synchronized across namespaces or clusters? Try [Kubed config syncer](/docs/guides/config-syncer/).
 - Want to keep an eye on your cluster with automated notifications? Setup Kubed [event forwarder](/docs/guides/cluster-events/).
 - Out of disk space because of too much logs in Elasticsearch or metrics in InfluxDB? Configure [janitors](/docs/guides/janitors.md) to delete old data.
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
