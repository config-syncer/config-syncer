package recyclebin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	stringz "github.com/appscode/go/strings"
	"github.com/appscode/kubed/pkg/config"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/ghodss/yaml"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/cache"
)

type RecycleBin struct {
	clusterName  string
	spec         *config.RecycleBinSpec
	notifierCred envconfig.LoaderFunc

	lock sync.RWMutex
}

var _ cache.ResourceEventHandler = &RecycleBin{}

func (c *RecycleBin) Configure(clusterName string, spec *config.RecycleBinSpec, notifierCred envconfig.LoaderFunc) {
	c.clusterName = clusterName
	c.spec = spec
	c.notifierCred = notifierCred
}

func (c *RecycleBin) OnAdd(obj interface{}) {}

func (c *RecycleBin) OnUpdate(oldObj, newObj interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.spec == nil {
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
	om, err := meta.Accessor(newObj)
	if err != nil {
		return err
	}

	tm, err := meta.TypeAccessor(newObj)
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
	fn := fmt.Sprintf("%s.%s.yaml", name, om.GetCreationTimestamp().UTC().Format(config.TimestampFormat))

	fullPath := filepath.Join(dir, fn)
	bytes, err := yaml.Marshal(newObj)
	if err != nil {
		return err
	}

	for _, receiver := range c.spec.Receivers {
		if len(receiver.To) > 0 {
			sub := fmt.Sprintf("[%s]: %s %s %s/%s updated", stringz.Val(c.clusterName, "?"), tm.GetAPIVersion(), tm.GetKind(), om.GetNamespace(), om.GetName())
			if notifier, err := unified.LoadVia(strings.ToLower(receiver.Notifier), c.notifierCred); err == nil {
				switch n := notifier.(type) {
				case notify.ByEmail:
					n = n.To(receiver.To[0], receiver.To[1:]...)
					if diff, err := meta_util.JsonDiff(oldObj, newObj); err == nil {
						n.WithSubject(sub).WithBody(diff).WithNoTracking().Send()
					} else {
						n.WithSubject(sub).WithBody(string(bytes)).WithNoTracking().Send()
					}
				case notify.BySMS:
					n.To(receiver.To[0], receiver.To[1:]...).
						WithBody(sub).
						Send()
				case notify.ByChat:
					n.To(receiver.To[0], receiver.To[1:]...).
						WithBody(sub).
						Send()
				case notify.ByPush:
					n.To(receiver.To...).
						WithBody(sub).
						Send()
				}
			}
		}
	}

	return ioutil.WriteFile(fullPath, bytes, 0644)
}

func (c *RecycleBin) delete(obj interface{}) error {
	om, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	tm, err := meta.TypeAccessor(obj)
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
	fn := fmt.Sprintf("%s.%s.yaml", name, om.GetCreationTimestamp().UTC().Format(config.TimestampFormat))

	fullPath := filepath.Join(dir, fn)
	bytes, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	for _, receiver := range c.spec.Receivers {
		if len(receiver.To) > 0 {
			sub := fmt.Sprintf("[%s]: %s %s %s/%s deleted", stringz.Val(c.clusterName, "?"), tm.GetAPIVersion(), tm.GetKind(), om.GetNamespace(), om.GetName())
			if notifier, err := unified.LoadVia(strings.ToLower(receiver.Notifier), c.notifierCred); err == nil {
				switch n := notifier.(type) {
				case notify.ByEmail:
					n.To(receiver.To[0], receiver.To[1:]...).
						WithSubject(sub).
						WithBody(string(bytes)).
						WithNoTracking().
						Send()
				case notify.BySMS:
					n.To(receiver.To[0], receiver.To[1:]...).
						WithBody(sub).
						Send()
				case notify.ByChat:
					n.To(receiver.To[0], receiver.To[1:]...).
						WithBody(sub).
						Send()
				case notify.ByPush:
					n.To(receiver.To...).
						WithBody(sub).
						Send()
				}
			}
		}
	}

	return ioutil.WriteFile(fullPath, bytes, 0644)
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
