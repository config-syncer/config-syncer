package kubedb

import (
	_ "k8s.io/apimachinery/pkg/apimachinery"
	_ "k8s.io/apimachinery/pkg/apimachinery/announced"
	_ "k8s.io/apimachinery/pkg/apimachinery/registered"
)

// GroupName is the group name use in this package
const GroupName = "kubedb.com"
