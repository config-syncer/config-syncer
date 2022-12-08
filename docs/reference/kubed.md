---
title: Kubed
menu:
  product_kubed_{{ .version }}:
    identifier: kubed
    name: Kubed
    parent: reference
    weight: 0

product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: reference
aliases:
  - /products/kubed/{{ .version }}/reference/

---
## kubed

Config Syncer by AppsCode - A Kubernetes Cluster Operator Daemon

### Synopsis

Config Syncer is a Kubernetes daemon to perform cluster management tasks. For more information, visit here: https://github.com/kubeops/config-syncer/tree/master/docs

### Options

```
  -h, --help                             help for kubed
      --use-kubeapiserver-fqdn-for-aks   if true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)
```

### SEE ALSO

* [kubed run](/docs/reference/kubed_run.md)	 - Launch Kubernetes Cluster Daemon
* [kubed version](/docs/reference/kubed_version.md)	 - Prints binary version number.

