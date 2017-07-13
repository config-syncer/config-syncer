package trashcan

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

type TrashCan struct {
	Spec config.TrashCanSpec
}

func (c *TrashCan) Update(meta metav1.ObjectMeta, old, new interface{}) error {
	p := filepath.Join(c.Spec.Path, meta.SelfLink)
	dir := filepath.Dir(p)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	name := filepath.Base(p)
	fn := fmt.Sprintf("%s.%s.yaml", name, meta.CreationTimestamp.UTC().Format(time.RFC3339))

	fullPath := filepath.Join(dir, fn)
	bytes, err := yaml.Marshal(new)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, bytes, 0644)
}

func (c *TrashCan) Delete(meta metav1.ObjectMeta, v interface{}) error {
	p := filepath.Join(c.Spec.Path, meta.SelfLink)
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

func (c *TrashCan) Cleanup() error {
	now := time.Now()
	return filepath.Walk(c.Spec.Path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if info.ModTime().Add(c.Spec.TTL.Duration).Before(now) {
				os.Remove(path)
			}
		}
		return nil
	})
}
