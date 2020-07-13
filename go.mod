module github.com/appscode/kubed

go 1.12

require (
	cloud.google.com/go v0.49.0 // indirect
	github.com/Azure/go-autorest/autorest v0.9.3-0.20191028180845-3492b2aff503 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.1-0.20191028180845-3492b2aff503 // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/appscode/go v0.0.0-20200323182826-54e98e09185a
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gophercloud/gophercloud v0.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.1 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/json-iterator/go v1.1.8
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1 // indirect
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.6 // indirect
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	go.uber.org/atomic v1.5.0 // indirect
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	golang.org/x/sys v0.0.0-20200509044756-6aff5f38e54f // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gomodules.xyz/cert v1.0.3
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	k8s.io/api v0.18.3
	k8s.io/apiextensions-apiserver v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/apiserver v0.18.3
	k8s.io/client-go v12.0.0+incompatible
	kmodules.xyz/client-go v0.0.0-20200713013004-81c7c59b7cfa
)

replace bitbucket.org/ww/goautoneg => gomodules.xyz/goautoneg v0.0.0-20120707110453-a547fc61f48d

replace git.apache.org/thrift.git => github.com/apache/thrift v0.13.0

replace github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v35.0.0+incompatible

replace github.com/Azure/go-ansiterm => github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.0.0+incompatible

replace github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.9.0

replace github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.5.0

replace github.com/Azure/go-autorest/autorest/azure/auth => github.com/Azure/go-autorest/autorest/azure/auth v0.2.0

replace github.com/Azure/go-autorest/autorest/date => github.com/Azure/go-autorest/autorest/date v0.1.0

replace github.com/Azure/go-autorest/autorest/mocks => github.com/Azure/go-autorest/autorest/mocks v0.2.0

replace github.com/Azure/go-autorest/autorest/to => github.com/Azure/go-autorest/autorest/to v0.2.0

replace github.com/Azure/go-autorest/autorest/validation => github.com/Azure/go-autorest/autorest/validation v0.1.0

replace github.com/Azure/go-autorest/logger => github.com/Azure/go-autorest/logger v0.1.0

replace github.com/Azure/go-autorest/tracing => github.com/Azure/go-autorest/tracing v0.5.0

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.5

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.0.0

replace go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace k8s.io/api => github.com/kmodules/api v0.18.4-0.20200524125823-c8bc107809b9

replace k8s.io/apimachinery => github.com/kmodules/apimachinery v0.19.0-alpha.0.0.20200520235721-10b58e57a423

replace k8s.io/apiserver => github.com/kmodules/apiserver v0.18.4-0.20200521000930-14c5f6df9625

replace k8s.io/client-go => k8s.io/client-go v0.18.3

replace k8s.io/kubernetes => github.com/kmodules/kubernetes v1.19.0-alpha.0.0.20200521033432-49d3646051ad
