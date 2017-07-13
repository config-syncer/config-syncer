package recover

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
	Opt config.RecoverSpec
}

func (c *RecoverStuff) Save(meta metav1.ObjectMeta, v interface{}) error {
	p := filepath.Join(c.Opt.Path, meta.SelfLink)
	dir := filepath.Dir(p)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	name := filepath.Base(p)
	fn := fmt.Sprintf("%s.%d.yaml", name, time.Now().UTC().Unix())

	fullPath := filepath.Join(dir, fn)
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, bytes, 0644)
}
