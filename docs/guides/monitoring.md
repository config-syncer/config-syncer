---
title: Monitoring
description: Monitoring
menu:
  product_kubed_{{ .version }}:
    identifier: monitoring-guide
    name: Monitoring
    parent: guides
    weight: 35
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: guides
---

> New to Config Syncer? Please start [here](/docs/concepts/README.md).

# Monitoring Config Syncer

Config Syncer operator exposes Prometheus ready metrics via the following endpoints on port `:8443`:

- `/metrics`: Scrape this to monitor operator.

Follow the steps below to view the metrics:

1. Give `system:anonymous` user access to `/metrics` url. **This is not safe to do on a production cluster.**

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: appscode:system:metrics-collector
rules:
- nonResourceURLs: ["/metrics"]
  verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: appscode:system:metrics-collector
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: appscode:system:metrics-collector
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: system:anonymous
```

```console
$ kubectl auth reconcile -f docs/examples/monitoring/metrics-collector.yaml
clusterrole.rbac.authorization.k8s.io "appscode:system:metrics-collector" reconciled
clusterrolebinding.rbac.authorization.k8s.io "appscode:system:metrics-collector" reconciled
```

2. Now, forward the port `:8443` to your workstation.

```
$ kubectl get pods -n kube-system | grep voyager
voyager-operator-f89dcccdb-plvmt        1/1       Running   0          27m

$ kubectl port-forward -n kube-system voyager-operator-f89dcccdb-plvmt 8443
Forwarding from 127.0.0.1:8443 -> 8443
Forwarding from [::1]:8443 -> 8443
```

3. Now, visit the url: https://127.0.0.1:8443/metrics

![operator-metrics](/docs/images/monitoring/operator-metrics.png)

4. Once you are done, remove access to `system:anonymous` user.

```console
$ kubectl delete -f docs/examples/monitoring/metrics-collector.yaml
clusterrole.rbac.authorization.k8s.io "appscode:system:metrics-collector" deleted
clusterrolebinding.rbac.authorization.k8s.io "appscode:system:metrics-collector" deleted
```

## Next Steps
 - Need to keep configmaps/secrets synchronized across namespaces or clusters? Try [Config Syncer config syncer](/docs/guides/config-syncer/).
 - Want to hack on Config Syncer? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
