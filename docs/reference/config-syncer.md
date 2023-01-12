---
title: Config-Syncer
menu:
  product_kubed_{{ .version }}:
    identifier: config-syncer
    name: Config-Syncer
    parent: reference
    weight: 0

product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: reference
aliases:
  - /products/kubed/{{ .version }}/reference/

---
## config-syncer

Config Syncer by AppsCode - A Kubernetes Cluster Operator Daemon

### Synopsis

Config Syncer is a Kubernetes daemon to perform cluster management tasks. For more information, visit here: https://github.com/kubeops/config-syncer/tree/master/docs

### Options

```
  -h, --help                             help for config-syncer
      --use-kubeapiserver-fqdn-for-aks   if true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)
```

### SEE ALSO

* [config-syncer run](/docs/reference/config-syncer_run.md)	 - Launch Kubernetes Cluster Daemon
* [config-syncer version](/docs/reference/config-syncer_version.md)	 - Prints binary version number.

