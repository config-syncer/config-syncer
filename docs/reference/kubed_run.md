---
title: Kubed Run
menu:
  product_kubed_{{ .version }}:
    identifier: kubed-run
    name: Kubed Run
    parent: reference
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: reference
---
## kubed run

Launch Kubernetes Cluster Daemon

### Synopsis

Launch Kubernetes Cluster Daemon

```
kubed run [flags]
```

### Options

```
      --audit-log-batch-buffer-size int                         The size of the buffer to store events before batching and writing. Only used in batch mode. (default 10000)
      --audit-log-batch-max-size int                            The maximum size of a batch. Only used in batch mode. (default 1)
      --audit-log-batch-max-wait duration                       The amount of time to wait before force writing the batch that hadn't reached the max size. Only used in batch mode.
      --audit-log-batch-throttle-burst int                      Maximum number of requests sent at the same moment if ThrottleQPS was not utilized before. Only used in batch mode.
      --audit-log-batch-throttle-enable                         Whether batching throttling is enabled. Only used in batch mode.
      --audit-log-batch-throttle-qps float32                    Maximum average number of batches per second. Only used in batch mode.
      --audit-log-compress                                      If set, the rotated log files will be compressed using gzip.
      --audit-log-format string                                 Format of saved audits. "legacy" indicates 1-line text format for each event. "json" indicates structured json format. Known formats are legacy,json. (default "json")
      --audit-log-maxage int                                    The maximum number of days to retain old audit log files based on the timestamp encoded in their filename.
      --audit-log-maxbackup int                                 The maximum number of old audit log files to retain.
      --audit-log-maxsize int                                   The maximum size in megabytes of the audit log file before it gets rotated.
      --audit-log-mode string                                   Strategy for sending audit events. Blocking indicates sending events should block server responses. Batch causes the backend to buffer and write events asynchronously. Known modes are batch,blocking,blocking-strict. (default "blocking")
      --audit-log-path string                                   If set, all requests coming to the apiserver will be logged to this file.  '-' means standard out.
      --audit-log-truncate-enabled                              Whether event and batch truncating is enabled.
      --audit-log-truncate-max-batch-size int                   Maximum size of the batch sent to the underlying backend. Actual serialized size can be several hundreds of bytes greater. If a batch exceeds this limit, it is split into several batches of smaller size. (default 10485760)
      --audit-log-truncate-max-event-size int                   Maximum size of the audit event sent to the underlying backend. If the size of an event is greater than this number, first request and response are removed, and if this doesn't reduce the size enough, event is discarded. (default 102400)
      --audit-log-version string                                API group and version used for serializing audit events written to log. (default "audit.k8s.io/v1")
      --audit-policy-file string                                Path to the file that defines the audit policy configuration.
      --audit-webhook-batch-buffer-size int                     The size of the buffer to store events before batching and writing. Only used in batch mode. (default 10000)
      --audit-webhook-batch-max-size int                        The maximum size of a batch. Only used in batch mode. (default 400)
      --audit-webhook-batch-max-wait duration                   The amount of time to wait before force writing the batch that hadn't reached the max size. Only used in batch mode. (default 30s)
      --audit-webhook-batch-throttle-burst int                  Maximum number of requests sent at the same moment if ThrottleQPS was not utilized before. Only used in batch mode. (default 15)
      --audit-webhook-batch-throttle-enable                     Whether batching throttling is enabled. Only used in batch mode. (default true)
      --audit-webhook-batch-throttle-qps float32                Maximum average number of batches per second. Only used in batch mode. (default 10)
      --audit-webhook-config-file string                        Path to a kubeconfig formatted file that defines the audit webhook configuration.
      --audit-webhook-initial-backoff duration                  The amount of time to wait before retrying the first failed request. (default 10s)
      --audit-webhook-mode string                               Strategy for sending audit events. Blocking indicates sending events should block server responses. Batch causes the backend to buffer and write events asynchronously. Known modes are batch,blocking,blocking-strict. (default "batch")
      --audit-webhook-truncate-enabled                          Whether event and batch truncating is enabled.
      --audit-webhook-truncate-max-batch-size int               Maximum size of the batch sent to the underlying backend. Actual serialized size can be several hundreds of bytes greater. If a batch exceeds this limit, it is split into several batches of smaller size. (default 10485760)
      --audit-webhook-truncate-max-event-size int               Maximum size of the audit event sent to the underlying backend. If the size of an event is greater than this number, first request and response are removed, and if this doesn't reduce the size enough, event is discarded. (default 102400)
      --audit-webhook-version string                            API group and version used for serializing audit events written to webhook. (default "audit.k8s.io/v1")
      --authentication-kubeconfig string                        kubeconfig file pointing at the 'core' kubernetes server with enough rights to create tokenreviews.authentication.k8s.io.
      --authentication-skip-lookup                              If false, the authentication-kubeconfig will be used to lookup missing authentication configuration from the cluster.
      --authentication-token-webhook-cache-ttl duration         The duration to cache responses from the webhook token authenticator. (default 10s)
      --authentication-tolerate-lookup-failure                  If true, failures to look up missing authentication configuration from the cluster are not considered fatal. Note that this can result in authentication that treats all requests as anonymous.
      --authorization-always-allow-paths strings                A list of HTTP paths to skip during authorization, i.e. these are authorized without contacting the 'core' kubernetes server. (default [/healthz,/readyz,/livez])
      --authorization-kubeconfig string                         kubeconfig file pointing at the 'core' kubernetes server with enough rights to create subjectaccessreviews.authorization.k8s.io.
      --authorization-webhook-cache-authorized-ttl duration     The duration to cache 'authorized' responses from the webhook authorizer. (default 10s)
      --authorization-webhook-cache-unauthorized-ttl duration   The duration to cache 'unauthorized' responses from the webhook authorizer. (default 10s)
      --bind-address ip                                         The IP address on which to listen for the --secure-port port. The associated interface(s) must be reachable by the rest of the cluster, and by CLI/web clients. If blank or an unspecified address (0.0.0.0 or ::), all interfaces will be used. (default 0.0.0.0)
      --burst int                                               The maximum burst for throttle (default 1000000)
      --cert-dir string                                         The directory where the TLS certs are located. If --tls-cert-file and --tls-private-key-file are provided, this flag will be ignored. (default "apiserver.local.config/certificates")
      --client-ca-file string                                   If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate.
      --cluster-name string                                     Name of cluster
      --config-source-namespace string                          Config source namespace
      --contention-profiling                                    Enable lock contention profiling, if profiling is enabled
      --egress-selector-config-file string                      File with apiserver egress selector configuration.
  -h, --help                                                    help for run
      --http2-max-streams-per-connection int                    The limit that the server gives to clients for the maximum number of streams in an HTTP/2 connection. Zero means to use golang's default. (default 1000)
      --kubeconfig string                                       kubeconfig file pointing at the 'core' kubernetes server.
      --kubeconfig-file string                                  kubeconfig file
      --permit-address-sharing                                  If true, SO_REUSEADDR will be used when binding the port. This allows binding to wildcard IPs like 0.0.0.0 and specific IPs in parallel, and it avoids waiting for the kernel to release sockets in TIME_WAIT state. [default=false]
      --permit-port-sharing                                     If true, SO_REUSEPORT will be used when binding the port, which allows more than one instance to bind on the same address and port. [default=false]
      --profiling                                               Enable profiling via web interface host:port/debug/pprof/ (default true)
      --qps float32                                             The maximum QPS to the master from this client (default 1e+06)
      --requestheader-allowed-names strings                     List of client certificate common names to allow to provide usernames in headers specified by --requestheader-username-headers. If empty, any client certificate validated by the authorities in --requestheader-client-ca-file is allowed.
      --requestheader-client-ca-file string                     Root certificate bundle to use to verify client certificates on incoming requests before trusting usernames in headers specified by --requestheader-username-headers. WARNING: generally do not depend on authorization being already done for incoming requests.
      --requestheader-extra-headers-prefix strings              List of request header prefixes to inspect. X-Remote-Extra- is suggested. (default [x-remote-extra-])
      --requestheader-group-headers strings                     List of request headers to inspect for groups. X-Remote-Group is suggested. (default [x-remote-group])
      --requestheader-username-headers strings                  List of request headers to inspect for usernames. X-Remote-User is common. (default [x-remote-user])
      --resync-period duration                                  If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out. (default 10m0s)
      --secure-port int                                         The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all. (default 443)
      --tls-cert-file string                                    File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert). If HTTPS serving is enabled, and --tls-cert-file and --tls-private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory specified by --cert-dir.
      --tls-cipher-suites strings                               Comma-separated list of cipher suites for the server. If omitted, the default Go cipher suites will be used. 
                                                                Preferred values: TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384, TLS_CHACHA20_POLY1305_SHA256, TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA, TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA, TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA, TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256, TLS_RSA_WITH_3DES_EDE_CBC_SHA, TLS_RSA_WITH_AES_128_CBC_SHA, TLS_RSA_WITH_AES_128_GCM_SHA256, TLS_RSA_WITH_AES_256_CBC_SHA, TLS_RSA_WITH_AES_256_GCM_SHA384. 
                                                                Insecure values: TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256, TLS_ECDHE_ECDSA_WITH_RC4_128_SHA, TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256, TLS_ECDHE_RSA_WITH_RC4_128_SHA, TLS_RSA_WITH_AES_128_CBC_SHA256, TLS_RSA_WITH_RC4_128_SHA.
      --tls-min-version string                                  Minimum TLS version supported. Possible values: VersionTLS10, VersionTLS11, VersionTLS12, VersionTLS13
      --tls-private-key-file string                             File containing the default x509 private key matching --tls-cert-file.
      --tls-sni-cert-key namedCertKey                           A pair of x509 certificate and private key file paths, optionally suffixed with a list of domain patterns which are fully qualified domain names, possibly with prefixed wildcard segments. The domain patterns also allow IP addresses, but IPs should only be used if the apiserver has visibility to the IP address requested by a client. If no domain patterns are provided, the names of the certificate are extracted. Non-wildcard matches trump over wildcard matches, explicit domain patterns trump over extracted names. For multiple key/certificate pairs, use the --tls-sni-cert-key multiple times. Examples: "example.crt,example.key" or "foo.crt,foo.key:*.foo.com,foo.com". (default [])
```

### Options inherited from parent commands

```
      --enable-analytics                 send usage events to Google Analytics (default true)
      --use-kubeapiserver-fqdn-for-aks   if true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)
```

### SEE ALSO

* [kubed](/docs/reference/kubed.md)	 - Config Syncer by AppsCode - A Kubernetes Cluster Operator Daemon

