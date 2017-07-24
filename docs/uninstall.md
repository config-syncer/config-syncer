> New to Kubed? Please start [here](/docs/tutorial.md).

# Uninstall Kubed
Please follow the steps below to uninstall Kubed:

1. Delete the deployment and service used for Kubed operator.
```console
$ kubectl delete deployment -l app=kubed -n <operator-namespace>
$ kubectl delete service -l app=kubed -n <operator-namespace>

# Delete RBAC objects, if --rbac flag was used.
$ kubectl delete serviceaccount -l app=kubed -n <operator-namespace>
$ kubectl delete clusterrolebindings -l app=kubed -n <operator-namespace>
$ kubectl delete clusterrole -l app=kubed -n <operator-namespace>
```

2. Now, wait several seconds for Kubed to stop running. To confirm that Kubed operator pod(s) have stopped running, run:
```console
$ kubectl get pods --all-namespaces -l app=kubed
```
