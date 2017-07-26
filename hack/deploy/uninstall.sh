#!/bin/bash
set -x

kubectl delete deployment -l app=kubed -n kube-system
kubectl delete service -l app=kubed -n kube-system
kubectl delete secret -l app=kubed -n kube-system

# Delete RBAC objects, if --rbac flag was used.
kubectl delete serviceaccount -l app=kubed -n kube-system
kubectl delete clusterrolebindings -l app=kubed -n kube-system
kubectl delete clusterrole -l app=kubed -n kube-system
