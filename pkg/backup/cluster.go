package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/appscode/go-term"
	"github.com/appscode/log"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	kapi "k8s.io/client-go/pkg/api"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
)

type ItemList struct {
	Items []map[string]interface{} `json:"items,omitempty"`
}

func SnapshotCluster(kubeConfig *rest.Config, backupDir string, sanitize bool) error {
	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(kubeConfig)
	rs, err := discoveryClient.ServerResources()
	if err != nil {
		return err
	}

	err = os.MkdirAll(backupDir, 0755)
	if err != nil {
		return err
	}
	resBytes, err := yaml.Marshal(rs)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(backupDir, "api_resources.yaml"), resBytes, 0755)
	if err != nil {
		return err
	}

	for _, v := range rs {
		gv, err := schema.ParseGroupVersion(v.GroupVersion)
		if err != nil {
			continue
		}
		for _, rss := range v.APIResources {
			log.Infoln("Taking backup for", rss.Name, "groupversion =", v.GroupVersion)
			if err := rest.SetKubernetesDefaults(kubeConfig); err != nil {
				return err
			}
			kubeConfig.ContentConfig = dynamic.ContentConfig()
			kubeConfig.GroupVersion = &schema.GroupVersion{Group: gv.Group, Version: gv.Version}
			kubeConfig.APIPath = "/apis"
			if gv.Group == kapi.GroupName {
				kubeConfig.APIPath = "/api"
			}
			restClient, err := rest.RESTClientFor(kubeConfig)
			if err != nil {
				return err
			}
			request := restClient.Get().Resource(rss.Name).Param("pretty", "true")
			b, err := request.DoRaw()
			if err != nil {
				log.Errorln(err)
				continue
			}
			list := &ItemList{}
			err = yaml.Unmarshal(b, &list)
			if err != nil {
				log.Errorln(err)
				continue
			}
			if len(list.Items) > 1000 {
				ok := term.Ask(fmt.Sprintf("Too many objects (%v). Want to take backup ?", len(list.Items)), true)
				if !ok {
					continue
				}
			}
			for _, ob := range list.Items {
				var selfLink string
				ob["apiVersion"] = v.GroupVersion
				ob["kind"] = rss.Kind
				i, ok := ob["metadata"]
				if ok {
					selfLink = getSelfLinkFromMetadata(i)
				} else {
					log.Errorln("Metadata not found")
					continue
				}
				if sanitize {
					cleanUpObjectMeta(i)
					spec, ok := ob["spec"].(map[string]interface{})
					if ok {
						if rss.Kind == "Pod" {
							spec = cleanUpPodSpec(spec)
						}
						template, ok := spec["template"].(map[string]interface{})
						if ok {
							podSpec, ok := template["spec"].(map[string]interface{})
							if ok {
								template["spec"] = cleanUpPodSpec(podSpec)
							}
						}
					}
					delete(ob, "status")
				}
				b, err := yaml.Marshal(ob)
				if err != nil {
					log.Errorln(err)
					break
				}
				path := filepath.Dir(filepath.Join(backupDir, selfLink))
				obName := filepath.Base(selfLink)
				err = os.MkdirAll(path, 0777)
				if err != nil {
					log.Errorln(err)
					break
				}
				fileName := filepath.Join(path, obName+".yaml")
				if err = ioutil.WriteFile(fileName, b, 0644); err != nil {
					log.Errorln(err)
					continue
				}
			}
		}
	}
	return nil
}

func cleanUpObjectMeta(i interface{}) {
	meta, ok := i.(map[string]interface{})
	if !ok {
		return
	}
	delete(meta, "creationTimestamp")
	delete(meta, "resourceVersion")
	delete(meta, "selfLink")
	delete(meta, "uid")
	delete(meta, "generateName")
	delete(meta, "generation")
	annotation, ok := meta["annotations"]
	if !ok {
		return
	}
	annotations, ok := annotation.(map[string]string)
	if !ok {
		return
	}
	cleanUpDecorators(annotations)
}

func cleanUpDecorators(i interface{}) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return
	}
	delete(m, "controller-uid")
	delete(m, "deployment.kubernetes.io/desired-replicas")
	delete(m, "deployment.kubernetes.io/max-replicas")
	delete(m, "deployment.kubernetes.io/revision")
	delete(m, "pod-template-hash")
	delete(m, "pv.kubernetes.io/bind-completed")
	delete(m, "pv.kubernetes.io/bound-by-controller")
}

func cleanUpPodSpec(podSpec map[string]interface{}) map[string]interface{} {
	b, err := yaml.Marshal(podSpec)
	if err != nil {
		term.Errorln(err)
		return podSpec
	}
	p := &kapi.PodSpec{}
	err = yaml.Unmarshal(b, p)
	if err != nil {
		term.Errorln(err)
		return podSpec // Not a podspec
	}
	p.DNSPolicy = kapi.DNSPolicy("")
	p.NodeName = ""
	if p.ServiceAccountName == "default" {
		p.ServiceAccountName = ""
	}
	p.TerminationGracePeriodSeconds = nil
	for i, c := range p.Containers {
		c.TerminationMessagePath = ""
		p.Containers[i] = c
	}
	for i, c := range p.InitContainers {
		c.TerminationMessagePath = ""
		p.InitContainers[i] = c
	}
	b, err = yaml.Marshal(p)
	if err != nil {
		term.Errorln(err)
		return nil
	}
	var spec map[string]interface{}
	err = yaml.Unmarshal(b, &spec)
	if err != nil {
		term.Errorln(err)
		return spec
	}
	return spec
}

func getSelfLinkFromMetadata(i interface{}) string {
	meta, ok := i.(map[string]interface{})
	if ok {
		return meta["selfLink"].(string)
	}
	return ""
}
