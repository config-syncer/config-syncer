---
title: API Server
description: API Server
menu:
  product_kubed_0.10.0:
    identifier: guides-apiserver
    name: API Server
    parent: guides
    weight: 30
product_name: kubed
menu_name: product_kubed_0.10.0
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

| ClusterRole             | Aggregates To       | Desription                                                                                                 |
| ----------------------- | ------------------- | ---------------------------------------------------------------------------------------------------------- |
| appscode:kubed:view     | admin, edit, view   | Allows read-only access to Kubed resources, intended to be granted within a namespace using a RoleBinding. |

These user facing roles supports [ClusterRole Aggregation](https://kubernetes.io/docs/admin/authorization/rbac/#aggregated-clusterroles) feature in Kubernetes 1.9 or later clusters.

### Search Kubernetes objects
To search for Kubernetes objects, use `kubectl get searchresult <search-term>`.

```console
$ kubectl get searchresult dashboard -n kube-system -o yaml

apiVersion: kubed.appscode.com/v1alpha1
hits:
- object:
    apiVersion: v1
    kind: Service
    metadata:
      annotations: {}
      creationTimestamp: 2018-03-03T00:41:14Z
      labels:
        addonmanager.kubernetes.io/mode: Reconcile
        app: kubernetes-dashboard
        kubernetes.io/minikube-addons: dashboard
        kubernetes.io/minikube-addons-endpoint: dashboard
      name: kubernetes-dashboard
      namespace: kube-system
      resourceVersion: "97"
      selfLink: /api/v1/namespaces/kube-system/services/kubernetes-dashboard
      uid: 93e060a2-1e7b-11e8-b100-0800274b2cc3
    spec:
      clusterIP: 10.97.199.12
      externalTrafficPolicy: Cluster
      ports:
      - nodePort: 30000
        port: 80
        protocol: TCP
        targetPort: 9090
      selector:
        app: kubernetes-dashboard
      sessionAffinity: None
      type: NodePort
    status:
      loadBalancer: {}
  score: 0.4400261
- object:
    apiVersion: v1
    data:
      priv: null
      pub: null
    kind: Secret
    metadata:
      creationTimestamp: 2018-03-03T00:41:18Z
      name: kubernetes-dashboard-key-holder
      namespace: kube-system
      resourceVersion: "154"
      selfLink: /api/v1/namespaces/kube-system/secrets/kubernetes-dashboard-key-holder
      uid: 9649eca5-1e7b-11e8-b100-0800274b2cc3
    type: Opaque
  score: 0.40519527
- object:
    apiVersion: extensions/v1beta1
    kind: ReplicaSet
    metadata:
      annotations:
        deployment.kubernetes.io/desired-replicas: "1"
        deployment.kubernetes.io/max-replicas: "2"
        deployment.kubernetes.io/revision: "1"
      creationTimestamp: 2018-03-03T00:41:14Z
      generation: 1
      labels:
        addonmanager.kubernetes.io/mode: Reconcile
        app: kubernetes-dashboard
        pod-template-hash: "3384654141"
        version: v1.8.1
      name: kubernetes-dashboard-77d8b98585
      namespace: kube-system
      ownerReferences:
      - apiVersion: extensions/v1beta1
        blockOwnerDeletion: true
        controller: true
        kind: Deployment
        name: kubernetes-dashboard
        uid: 93cb3dc4-1e7b-11e8-b100-0800274b2cc3
      resourceVersion: "158"
      selfLink: /apis/extensions/v1beta1/namespaces/kube-system/replicasets/kubernetes-dashboard-77d8b98585
      uid: 93cbe990-1e7b-11e8-b100-0800274b2cc3
    spec:
      replicas: 1
      selector:
        matchLabels:
          addonmanager.kubernetes.io/mode: Reconcile
          app: kubernetes-dashboard
          pod-template-hash: "3384654141"
          version: v1.8.1
      template:
        metadata:
          creationTimestamp: null
          labels:
            addonmanager.kubernetes.io/mode: Reconcile
            app: kubernetes-dashboard
            pod-template-hash: "3384654141"
            version: v1.8.1
        spec:
          containers:
          - image: k8s.gcr.io/kubernetes-dashboard-amd64:v1.8.1
            imagePullPolicy: IfNotPresent
            livenessProbe:
              failureThreshold: 3
              httpGet:
                path: /
                port: 9090
                scheme: HTTP
              initialDelaySeconds: 30
              periodSeconds: 10
              successThreshold: 1
              timeoutSeconds: 30
            name: kubernetes-dashboard
            ports:
            - containerPort: 9090
              protocol: TCP
            resources: {}
            terminationMessagePath: /dev/termination-log
            terminationMessagePolicy: File
          dnsPolicy: ClusterFirst
          restartPolicy: Always
          schedulerName: default-scheduler
          securityContext: {}
          terminationGracePeriodSeconds: 30
    status:
      availableReplicas: 1
      fullyLabeledReplicas: 1
      observedGeneration: 1
      readyReplicas: 1
      replicas: 1
  score: 0.3066834
- object:
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      annotations:
        deployment.kubernetes.io/revision: "1"
      creationTimestamp: 2018-03-03T00:41:14Z
      generation: 1
      labels:
        addonmanager.kubernetes.io/mode: Reconcile
        kubernetes.io/minikube-addons: dashboard
        version: v1.8.1
      name: kubernetes-dashboard
      namespace: kube-system
      resourceVersion: "159"
      selfLink: /apis/apps/v1beta1/namespaces/kube-system/deployments/kubernetes-dashboard
      uid: 93cb3dc4-1e7b-11e8-b100-0800274b2cc3
    spec:
      progressDeadlineSeconds: 600
      replicas: 1
      revisionHistoryLimit: 10
      selector:
        matchLabels:
          addonmanager.kubernetes.io/mode: Reconcile
          app: kubernetes-dashboard
          version: v1.8.1
      strategy:
        rollingUpdate:
          maxSurge: 25%
          maxUnavailable: 25%
        type: RollingUpdate
      template:
        metadata:
          creationTimestamp: null
          labels:
            addonmanager.kubernetes.io/mode: Reconcile
            app: kubernetes-dashboard
            version: v1.8.1
        spec:
          containers:
          - image: k8s.gcr.io/kubernetes-dashboard-amd64:v1.8.1
            imagePullPolicy: IfNotPresent
            livenessProbe:
              failureThreshold: 3
              httpGet:
                path: /
                port: 9090
                scheme: HTTP
              initialDelaySeconds: 30
              periodSeconds: 10
              successThreshold: 1
              timeoutSeconds: 30
            name: kubernetes-dashboard
            ports:
            - containerPort: 9090
              protocol: TCP
            resources: {}
            terminationMessagePath: /dev/termination-log
            terminationMessagePolicy: File
          dnsPolicy: ClusterFirst
          restartPolicy: Always
          schedulerName: default-scheduler
          securityContext: {}
          terminationGracePeriodSeconds: 30
    status:
      availableReplicas: 1
      conditions:
      - lastTransitionTime: 2018-03-03T00:41:18Z
        lastUpdateTime: 2018-03-03T00:41:18Z
        message: Deployment has minimum availability.
        reason: MinimumReplicasAvailable
        status: "True"
        type: Available
      - lastTransitionTime: 2018-03-03T00:41:14Z
        lastUpdateTime: 2018-03-03T00:41:18Z
        message: ReplicaSet "kubernetes-dashboard-77d8b98585" has successfully progressed.
        reason: NewReplicaSetAvailable
        status: "True"
        type: Progressing
      observedGeneration: 1
      readyReplicas: 1
      replicas: 1
      updatedReplicas: 1
  score: 0.25559625
kind: SearchResult
maxScore: 0.4400261
metadata:
  creationTimestamp: null
  name: dashboard
  namespace: kube-system
  selfLink: /apis/kubed.appscode.com/v1alpha1/namespaces/kube-system/searchresults/dashboard
took: 121.049µs
totalHits: 4
```

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
