---
title: Kubed Backup
menu:
  product_kubed_v0.11.0:
    identifier: kubed-backup
    name: Kubed Backup
    parent: reference
product_name: kubed
menu_name: product_kubed_v0.11.0
section_menu_id: reference
---
## kubed backup

Takes a backup of Kubernetes api objects

### Synopsis

Takes a backup of Kubernetes api objects

```
kubed backup [flags]
```

### Options

```
      --backup-dir string   Directory where YAML files will be stored
      --context string      Name of the kubeconfig context to use
  -h, --help                help for backup
      --kubeconfig string   kubeconfig file pointing at the 'core' kubernetes server
      --sanitize             Sanitize fields in YAML
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --enable-analytics                 send usage events to Google Analytics (default true)
      --log-flush-frequency duration     Maximum number of seconds between log flushes (default 5s)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr
      --use-kubeapiserver-fqdn-for-aks   if true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [kubed](/docs/reference/kubed.md)	 - Kubed by AppsCode - A Kubernetes Cluster Operator Daemon

