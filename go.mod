module github.com/appscode/kubed

go 1.12

require (
	github.com/RoaringBitmap/roaring v0.0.0-20180103163510-cefad6e4f79d // indirect
	github.com/Smerity/govarint v0.0.0-20150407073650-7265e41f48f1 // indirect
	github.com/appscode/go v0.0.0-20190621064509-6b292c9166e3
	github.com/appscode/osm v0.12.0
	github.com/appscode/searchlight v0.0.0-20190604163604-8a6c4c21504d
	github.com/appscode/voyager v0.0.0-20190717194632-664837a7f865
	github.com/aws/aws-sdk-go v1.20.20
	github.com/blevesearch/bleve v0.7.0
	github.com/blevesearch/blevex v0.0.0-20180227211930-4b158bb555a3 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.2 // indirect
	github.com/blevesearch/segment v0.0.0-20160105220820-db70c57796cc // indirect
	github.com/boltdb/bolt v0.0.0-20161028193645-4b1ebc1869ad // indirect
	github.com/codeskyblue/go-sh v0.0.0-20190412065543-76bd3d59ff27
	github.com/coreos/prometheus-operator v0.31.1
	github.com/couchbase/vellum v0.0.0-20190328134517-462e86d8716b // indirect
	github.com/cznic/b v0.0.0-20181122101859-a26611c4d92d // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/cznic/strutil v0.0.0-20181122101858-275e90344537 // indirect
	github.com/edsrzf/mmap-go v0.0.0-20160512033002-935e0e8a636c // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/glycerine/go-unsnap-stream v0.0.0-20171127062821-62a9a9eb44fd // indirect
	github.com/glycerine/goconvey v0.0.0-20190410193231-58a59202ab31 // indirect
	github.com/go-openapi/spec v0.19.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gophercloud/gophercloud v0.0.0-20190516144603-ad4210895ed0 // indirect
	github.com/influxdata/influxdb v1.5.3
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/json-iterator/go v1.1.6
	github.com/kubedb/apimachinery v0.0.0-20190717225257-5f6c0b43c6e3
	github.com/mschoch/smat v0.0.0-20160514031455-90eadee771ae // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/philhofer/fwd v0.0.0-20170616204054-1612a2981176 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/common v0.4.1
	github.com/remyoudompheng/bigfft v0.0.0-20190512091148-babf20351dd7 // indirect
	github.com/robfig/cron/v3 v3.0.0
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/steveyen/gtreap v0.0.0-20150807155958-0abe01ef9be2 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tecbot/gorocksdb v0.0.0-20181010114359-8752a9433481 // indirect
	github.com/tinylib/msgp v0.0.0-20160803062324-ad0ff2e232ad // indirect
	github.com/willf/bitset v0.0.0-20160225150313-2e6e8094ef47 // indirect
	gomodules.xyz/cert v1.0.0
	gomodules.xyz/envconfig v1.3.1-0.20190308184047-426f31af0d45
	gomodules.xyz/notify v0.0.0-20190424183923-af47cb5a07a4
	gomodules.xyz/stow v0.2.0
	gopkg.in/olivere/elastic.v5 v5.0.61
	k8s.io/api v0.0.0-20190515023547-db5a9d1c40eb
	k8s.io/apiextensions-apiserver v0.0.0-20190516231611-bf6753f2aa24
	k8s.io/apimachinery v0.0.0-20190515023456-b74e4c97951f
	k8s.io/apiserver v0.0.0-20190516230822-f89599b3f645
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/kube-aggregator v0.0.0-20190515024249-81a6edcf70be
	k8s.io/kube-openapi v0.0.0-20190510232812-a01b7d5d6c22
	kmodules.xyz/client-go v0.0.0-20190715080709-7162a6c90b04
	kmodules.xyz/monitoring-agent-api v0.0.0-20190513065523-186af167f817
	kmodules.xyz/objectstore-api v0.0.0-20190715131036-28ef070833ef
	kmodules.xyz/webhook-runtime v0.0.0-20190715115250-a84fbf77dd30
	stash.appscode.dev/stash v0.0.0-20190717091949-1c51a7c58f08
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest/autorest v0.5.0
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.4
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20190508045248-a52a97a7a2bf
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20190508082252-8397d761d4b5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/utils => k8s.io/utils v0.0.0-20190221042446-c2654d5206da
)
