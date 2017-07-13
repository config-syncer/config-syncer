package recover

import (
	"github.com/appscode/kubed/pkg/config"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	"os"
	"fmt"
	"io/ioutil"
	"github.com/ghodss/yaml"
)

type RecoverStuff struct {
	Opt config.RecoverSpec
}

func (c * RecoverStuff) Save(v interface{}) error {
	meta := v.(metav1.ObjectMeta)
	p := filepath.Join(c.Opt.Path, meta.SelfLink)
	dir := filepath.Dir(p)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	name := filepath.Base(p)
	fn := fmt.Sprintf("%s.%d.yaml", name, 	meta.CreationTimestamp.Unix())

	fullPath := filepath.Join(dir, fn)
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, bytes, 0644)
}
