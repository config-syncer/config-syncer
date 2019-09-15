---
title: Kubed Uninstall
description: Kubed Uninstall
menu:
  product_kubed_{{ .version }}:
    identifier: kubed-uninstall
    name: Uninstall
    parent: setup
    weight: 20
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: setup
---

> New to Kubed? Please start [here](/docs/concepts/README.md).

# Uninstall Kubed
Please follow the steps below to uninstall Kubed:

- Delete the various objects created for Kubed operator.

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/kubed/{{< param "info.version" >}}/hack/deploy/kubed.sh \
    | bash -s -- --uninstall [--namespace=NAMESPACE]

+ kubectl delete deployment -l app=kubed -n kube-system
deployment "kubed-operator" deleted
+ kubectl delete service -l app=kubed -n kube-system
service "kubed-operator" deleted
+ kubectl delete secret -l app=kubed -n kube-system
secret "azure-secret" deleted
secret "kubed-config" deleted
+ kubectl delete serviceaccount -l app=kubed -n kube-system
No resources found
+ kubectl delete clusterrolebindings -l app=kubed -n kube-system
No resources found
+ kubectl delete clusterrole -l app=kubed -n kube-system
No resources found
```

- Now, wait several seconds for Kubed to stop running. To confirm that Kubed operator pod(s) have stopped running, run:

```console
$ kubectl get pods --all-namespaces -l app=kubed
```
