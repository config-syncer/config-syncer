package main

import (
	gort "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/apis/kubed/install"
	"github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kutil/openapi"
	"github.com/go-openapi/spec"
	"github.com/golang/glog"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kube-openapi/pkg/common"
	"os"
	"path/filepath"
)

func generateSwaggerJson() {
	var (
		groupFactoryRegistry = make(announced.APIGroupFactoryRegistry)
		registry             = registered.NewOrDie("")
		Scheme               = runtime.NewScheme()
		Codecs               = serializer.NewCodecFactory(Scheme)
	)

	install.Install(groupFactoryRegistry, registry, Scheme)

	apispec, err := openapi.RenderOpenAPISpec(openapi.Config{
		Registry: registry,
		Scheme:   Scheme,
		Codecs:   Codecs,
		Info: spec.InfoProps{
			Title:   "Kubed",
			Version: "v0.7.0",
			Contact: &spec.ContactInfo{
				Name:  "AppsCode Inc.",
				URL:   "https://appscode.com",
				Email: "hello@appscode.com",
			},
			License: &spec.License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		OpenAPIDefinitions: []common.GetOpenAPIDefinitions{
			v1alpha1.GetOpenAPIDefinitions,
		},
		GetterResources: []schema.GroupVersionResource{
			v1alpha1.SchemeGroupVersion.WithResource("searchresults"),
		},
	})
	if err != nil {
		glog.Fatal(err)
	}

	filename := gort.GOPath() + "/src/github.com/appscode/kubed/api/openapi-spec/swagger.json"
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		glog.Fatal(err)
	}
	err = ioutil.WriteFile(filename, []byte(apispec), 0644)
	if err != nil {
		glog.Fatal(err)
	}
}

func main() {
	generateSwaggerJson()
}
