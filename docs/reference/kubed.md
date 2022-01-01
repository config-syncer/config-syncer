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

Kubed by AppsCode - A Kubernetes Cluster Operator Daemon

### Synopsis

Kubed is a Kubernetes daemon to perform cluster management tasks. For more information, visit here: https://github.com/kubeops/config-syncer/tree/master/docs

### Options

```
      --alsologtostderr                  log to standard error as well as files
      --enable-analytics                 send usage events to Google Analytics (default true)
  -h, --help                             help for kubed
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

* [kubed run](/docs/reference/kubed_run.md)	 - Launch Kubernetes Cluster Daemon
* [kubed version](/docs/reference/kubed_version.md)	 - Prints binary version number.

