---
title: Kubed Run
menu:
  product_kubed_0.3.1:
    identifier: kubed-run
    name: Kubed Run
    parent: reference
product_name: kubed
menu_name: product_kubed_0.3.1
section_menu_id: reference
---
## kubed run

Run daemon

### Synopsis

Run daemon

```
kubed run [flags]
```

### Options

```
      --api.address string       The address of the Kubed API Server (overrides any value in clusterconfig) (default ":8080")
      --clusterconfig string     Path to cluster config file (default "/srv/kubed/config.yaml")
  -h, --help                     help for run
      --kubeconfig string        Path to kubeconfig file with authorization information (the master location is set by the master flag).
      --master string            The address of the Kubernetes API server (overrides any value in kubeconfig)
      --resync-period duration   If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out. (default 5m0s)
      --scratch-dir emptyDir     Directory used to store temporary files. Use an emptyDir in Kubernetes. (default "/tmp")
      --web.address string       Address to listen on for web interface and telemetry. (default ":56790")
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

