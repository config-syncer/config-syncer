package trashcan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/kubed/pkg/config"
	"github.com/ghodss/yaml"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TrashCan struct {
	Spec   config.TrashCanSpec
	Loader envconfig.LoaderFunc
}

func (c *TrashCan) Update(t metav1.TypeMeta, meta metav1.ObjectMeta, old, new interface{}) error {
	p := filepath.Join(c.Spec.Path, meta.SelfLink)
	dir := filepath.Dir(p)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	name := filepath.Base(p)
	fn := fmt.Sprintf("%s.%s.yaml", name, meta.CreationTimestamp.UTC().Format(config.TimestampFormat))

	fullPath := filepath.Join(dir, fn)
	bytes, err := yaml.Marshal(new)
	if err != nil {
		return err
	}

	if c.Spec.NotifyVia != "" {
		sub := fmt.Sprintf("%s %s %s/%s updated", t.APIVersion, t.Kind, meta.Namespace, meta.Name)
		if notifier, err := unified.LoadVia(c.Spec.NotifyVia, c.Loader); err == nil {
			switch n := notifier.(type) {
			case notify.ByEmail:
				if diff, err := prepareDiff(old, new); err == nil {
					n.WithSubject(sub).WithBody(diff).WithNoTracking().Send()
				} else {
					n.WithSubject(sub).WithBody(string(bytes)).WithNoTracking().Send()
				}
			case notify.BySMS:
				n.WithBody(sub).Send()
			case notify.ByChat:
				n.WithBody(sub).Send()
			}
		}
	}

	return ioutil.WriteFile(fullPath, bytes, 0644)
}

func (c *TrashCan) Delete(t metav1.TypeMeta, meta metav1.ObjectMeta, v interface{}) error {
	p := filepath.Join(c.Spec.Path, meta.SelfLink)
	dir := filepath.Dir(p)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	name := filepath.Base(p)
	fn := fmt.Sprintf("%s.%s.yaml", name, meta.CreationTimestamp.UTC().Format(config.TimestampFormat))

	fullPath := filepath.Join(dir, fn)
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	if c.Spec.NotifyVia != "" {
		sub := fmt.Sprintf("%s %s %s/%s deleted", t.APIVersion, t.Kind, meta.Namespace, meta.Name)
		if notifier, err := unified.LoadVia(c.Spec.NotifyVia, c.Loader); err == nil {
			switch n := notifier.(type) {
			case notify.ByEmail:
				n.WithSubject(sub).WithBody(string(bytes)).WithNoTracking().Send()
			case notify.BySMS:
				n.WithBody(sub).Send()
			case notify.ByChat:
				n.WithBody(sub).Send()
			}
		}
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

func prepareDiff(old, new interface{}) (string, error) {
	oldBytes, err := json.Marshal(old)
	if err != nil {
		return "", err
	}

	newBytes, err := json.Marshal(new)
	if err != nil {
		return "", err
	}

	// Then, compare them
	differ := diff.New()
	d, err := differ.Compare(oldBytes, newBytes)
	if err != nil {
		return "", err
	}

	var aJson map[string]interface{}
	if err := json.Unmarshal(oldBytes, &aJson); err != nil {
		return "", err
	}

	format := formatter.NewAsciiFormatter(aJson, formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       false,
	})
	return format.Format(d)
}
