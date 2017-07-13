package recyclebin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/appscode/kubed/pkg/config"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RecoverStuff struct {
	Opt config.RecycleBinSpec
}

func (c *RecoverStuff) Save(meta metav1.ObjectMeta, v interface{}) error {
	p := filepath.Join(c.Opt.Path, meta.SelfLink)
	dir := filepath.Dir(p)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	name := filepath.Base(p)
	fn := fmt.Sprintf("%s.%s.yaml", name, meta.CreationTimestamp.UTC().Format(time.RFC3339))

	fullPath := filepath.Join(dir, fn)
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, bytes, 0644)
}
