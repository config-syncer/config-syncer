## kubed run

Run daemon

### Synopsis


Run daemon

```
kubed run [flags]
```

### Options

```
      --address string                         The address of the Kubed API Server (default ":32600")
      --analytics                              Send analytical events to Google Analytics (default true)
      --enable-reverse-index                   Reverse indexing of pods to service and others (default true)
      --es-endpoint string                     Endpoint of elasticsearch
  -h, --help                                   help for run
      --indexer string                         Reverse indexing of pods to service and others (default "indexers.bleve")
      --influx-secret string                   Influxdb secret name (default "appscode-influx")
      --influx-secret-namespace string         Influxdb secret namespace (default "kube-system")
      --kubeconfig string                      Path to kubeconfig file with authorization information (the master location is set by the master flag).
      --kubed-config-secret-name string        Kubed configuration secret name (default "cluster-kubed-config")
      --kubed-config-secret-namespace string   Kubed configuration secret namespace (default "kube-system")
      --master string                          The address of the Kubernetes API server (overrides any value in kubeconfig)
      --notify-on-cert-expired                 If enabled notify cluster admin wheen cert expired soon. (default true)
      --notify-via string                      Default notification method (eg: hipchat, mailgun, smtp, twilio, slack, plivo) (default "plivo")
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO
* [kubed](kubed.md)	 - Kubed by AppsCode - Kubernetes Daemon


