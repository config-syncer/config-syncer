---
title: Kubed Version
menu:
  product_kubed_{{ .version }}:
    identifier: kubed-version
    name: Kubed Version
    parent: reference
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: reference
---
## kubed version

Prints binary version number.

```
kubed version [flags]
```

### Options

```
      --check string   Check version constraint
  -h, --help           help for version
      --short          Print just the version number.
```

### Options inherited from parent commands

```
      --enable-analytics                 send usage events to Google Analytics (default true)
      --use-kubeapiserver-fqdn-for-aks   if true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)
```

### SEE ALSO

* [kubed](/docs/reference/kubed.md)	 - Config Syncer by AppsCode - A Kubernetes Cluster Operator Daemon

