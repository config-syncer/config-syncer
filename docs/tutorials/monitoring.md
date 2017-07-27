> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Monitoring Kubed

KubeDB has native support for monitoring via Prometheus. KubeDB operator exposes Prometheus native monitoring data via `/metrics` endpoint on `:56790` port. You can setup a [CoreOS Prometheus ServiceMonitor](https://github.com/coreos/prometheus-operator) using `kubed-operator` service. To change the port, use `--web.address` flag on Kubed operator.

```console
$ kubectl get pods -n kube-system
NAME                              READY     STATUS    RESTARTS   AGE
kube-addon-manager-minikube       1/1       Running   0          33m
kube-dns-1301475494-hglm0         3/3       Running   0          33m
kubed-operator-3234987584-sbgrf   1/1       Running   0          19s
kubernetes-dashboard-l8vlj        1/1       Running   0          33m


$ kubectl port-forward $(kubectl get pods --all-namespaces -l app=kubed -o jsonpath={.items[0].metadata.name}) -n kube-system 56790
Forwarding from 127.0.0.1:56790 -> 56790
E0727 03:50:34.668103   22871 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:56790: bind: cannot assign requested address
Handling connection for 56790
```

Now, open the URL [http://127.0.0.1:56790/metrics](http://127.0.0.1:56790/metrics) in your browser.
