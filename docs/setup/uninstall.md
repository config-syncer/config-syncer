---
title: Kubed Uninstall
description: Kubed Uninstall
menu:
  product_kubed_0.3.1:
    identifier: kubed-uninstall
    name: Uninstall
    parent: welcome
    weight: 30
product_name: kubed
menu_name: product_kubed_0.3.1
section_menu_id: welcome
url: /products/kubed/0.3.1/welcome/uninstall/
aliases:
  - /products/kubed/0.3.1/uninstall/
---

> New to Kubed? Please start [here](/docs/guides/README.md).

# Uninstall Kubed
Please follow the steps below to uninstall Kubed:

1. Delete the various objects created for Kubed operator.
```console
$ ./hack/deploy/uninstall.sh
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

2. Now, wait several seconds for Kubed to stop running. To confirm that Kubed operator pod(s) have stopped running, run:
```console
$ kubectl get pods --all-namespaces -l app=kubed
```
