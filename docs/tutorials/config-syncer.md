> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Synchronize Configuration Across Namespaces
Sometimes you have a configuration data that you want to synchronize across all namespaces. Kubed can do that for you. If a ConfigMap or Secret has the label `kubed.appscode.com/sync: true`, Kubed will create a similar ConfigMap / Secret in all existing namespaces. Kubed will also create this ConfigMap/Secret, when you create a new namespace. If the data in the source ConfigMap/Secret is updated, all the copies will be updated.


Sync ConfigMaps and Secrets


kubectl create configmap special-config --from-literal=special.how=very --from-literal=special.type=charm
