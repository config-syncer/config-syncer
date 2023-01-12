---
title: Config-Syncer Version
menu:
  product_kubed_{{ .version }}:
    identifier: config-syncer-version
    name: Config-Syncer Version
    parent: reference
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: reference
---
## config-syncer version

Prints binary version number.

```
config-syncer version [flags]
```

### Options

```
      --check string   Check version constraint
  -h, --help           help for version
      --short          Print just the version number.
```

### Options inherited from parent commands

```
      --use-kubeapiserver-fqdn-for-aks   if true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)
```

### SEE ALSO

* [config-syncer](/docs/reference/config-syncer.md)	 - Config Syncer by AppsCode - A Kubernetes Configuration Syncer

