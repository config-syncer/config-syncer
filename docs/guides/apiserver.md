---
title: API Server
description: API Server
menu:
  product_kubed_0.6.0-rc.0:
    identifier: guides-apiserver
    name: API Server
    parent: guides
    weight: 30
product_name: kubed
menu_name: product_kubed_0.6.0-rc.0
section_menu_id: guides
---

> New to Kubed? Please start [here](/docs/concepts/README.md).

# Kubed API Server

Kubed includes a search engine based on [blevesearch/bleve](https://github.com/blevesearch/bleve). This is exposed as a
Kubernetes [extended api server](https://kubernetes.io/docs/concepts/api-extension/apiserver-aggregation/). So, you can
just use `kubectl` to find any object by name in a namespace.


## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, deploy Kubed in your cluster following the steps [here](/docs/setup/install.md). Once the operator pod is running, go to the next section.


## Using Kubed API Server
In this section, we will show how you can use the kubed api server.

### Configuring RBAC
Kubed creates a custom resource: `SearchResult`. Kubed installer will create a user facing cluster role:

| ClusterRole           | Aggregates To     | Desription                            |
|-----------------------|-------------------|---------------------------------------|
| appscode:voyager:view | admin, edit, view | Allows read-only access to Kubed resources, intended to be granted within a namespace using a RoleBinding. |

These user facing roles supports [ClusterRole Aggregation](https://kubernetes.io/docs/admin/authorization/rbac/#aggregated-clusterroles) feature in Kubernetes 1.9 or later clusters.

### Search Kubernetes objects
To search for Kubernetes objects, use `kubectl get searchresult <search-term>`.

```console
$ kubectl get searchresult dashboard

$ kubectl get searchresult dashboard -n kube-system
NAME                              READY     STATUS    RESTARTS   AGE
kube-addon-manager-minikube       1/1       Running   0          33m
kube-dns-1301475494-hglm0         3/3       Running   0          33m
kubed-operator-3234987584-sbgrf   1/1       Running   0          19s
kubernetes-dashboard-l8vlj        1/1       Running   0          33m

```

Now, open the URL [http://127.0.0.1:8080/search?q=dashboard](http://127.0.0.1:8080/search?q=dashboard) in your browser.


### Programmtic Access to API Server

Kubed project includes an auto generated Go client for `SearchResult` resource. You can find it [here](https://github.com/appscode/kubed/blob/master/client/clientset/versioned/clientset.go).
This can be used like [k8s.io/client-go](https://github.com/kubernetes/client-go) library to access `SearchResult` programmatically.


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, uninstall Kubed following the steps [here](/docs/setup/uninstall.md).


## Next Steps
 - Learn how to use Kubed to protect your Kubernetes cluster from disasters [here](/docs/guides/disaster-recovery/).
 - Need to keep configmaps/secrets synchronized across namespaces or clusters? Try [Kubed config syncer](/docs/guides/config-syncer/).
 - Want to keep an eye on your cluster with automated notifications? Setup Kubed [event forwarder](/docs/guides/cluster-events/).
 - Out of disk space because of too much logs in Elasticsearch or metrics in InfluxDB? Configure [janitors](/docs/guides/janitors.md) to delete old data.
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
