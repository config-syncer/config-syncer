---
title: Kubed Snapshot
menu:
  product_kubed_0.3.1:
    identifier: kubed-snapshot
    name: Kubed Snapshot
    parent: reference
product_name: kubed
menu_name: product_kubed_0.3.1
section_menu_id: reference
---
## kubed snapshot

Takes a snapshot of Kubernetes api objects

### Synopsis

Takes a snapshot of Kubernetes api objects

```
kubed snapshot [flags]
```

### Options

```
      --backup-dir string   Directory where YAML files will be stored
      --context string      Name of the kubeconfig context to use
  -h, --help                help for snapshot
      --sanitize             Sanitize fields in YAML
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --analytics                        Send analytical events to Google Analytics (default true)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [kubed](/docs/reference/kubed.md)	 - Kubed by AppsCode - A Kubernetes Cluster Operator Daemon

