package recyclebin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/ghodss/yaml"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/cache"
)

type RecycleBin struct {
	clusterName string
	spec        *api.RecycleBinSpec

	lock sync.RWMutex
}

var _ cache.ResourceEventHandler = &RecycleBin{}

func (c *RecycleBin) Configure(clusterName string, spec *api.RecycleBinSpec) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.clusterName = clusterName
	c.spec = spec

	return nil
}

func (c *RecycleBin) OnAdd(obj interface{}) {}

func (c *RecycleBin) OnUpdate(oldObj, newObj interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.spec == nil || !c.spec.HandleUpdates {
		return
	}

	if err := c.update(oldObj, newObj); err != nil {
		log.Errorln(err)
	}
}

func (c *RecycleBin) OnDelete(obj interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.spec == nil {
		return
	}

	if err := c.delete(obj); err != nil {
		log.Errorln(err)
	}
}

func (c *RecycleBin) update(oldObj, newObj interface{}) error {
	om, err := meta.Accessor(oldObj)
	if err != nil {
		return err
	}

	p := filepath.Join(c.spec.Path, om.GetSelfLink())
	dir := filepath.Dir(p)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	name := filepath.Base(p)
	fn := fmt.Sprintf("%s.%s.yaml", name, om.GetCreationTimestamp().UTC().Format(api.TimestampFormat))

	fullPath := filepath.Join(dir, fn)
	data, err := yaml.Marshal(oldObj)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fullPath, data, 0644)
}

func (c *RecycleBin) delete(obj interface{}) error {
	om, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	p := filepath.Join(c.spec.Path, om.GetSelfLink())
	dir := filepath.Dir(p)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	name := filepath.Base(p)
	fn := fmt.Sprintf("%s.%s.yaml", name, om.GetCreationTimestamp().UTC().Format(api.TimestampFormat))

	fullPath := filepath.Join(dir, fn)
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fullPath, data, 0644)
}

func (c *RecycleBin) Cleanup() error {
	now := time.Now()
	return filepath.Walk(c.spec.Path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if info.ModTime().Add(c.spec.TTL.Duration).Before(now) {
				os.Remove(path)
			}
		}
		return nil
	})
}
