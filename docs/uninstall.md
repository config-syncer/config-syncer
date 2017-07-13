> New to Kubed? Please start [here](/docs/tutorial.md).

# Uninstall Kubed
Please follow the steps below to uninstall Kubed:

1. Delete the deployment and service used for Kubed operator.
```sh
$ kubectl delete deployment -l app=kubed -n <operator-namespace>
$ kubectl delete service -l app=kubed -n <operator-namespace>
```

2. Now, wait several seconds for Kubed to stop running. To confirm that Kubed operator pod(s) have stopped running, run:
```sh
$ kubectl get pods --all-namespaces -l app=kubed
```
