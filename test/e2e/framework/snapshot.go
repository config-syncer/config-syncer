package framework

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/appscode/go/encoding/yaml"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/pkg/storage"
	"github.com/graymeta/stow"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	apps "k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	TEST_LOCAL_BACKUP_DIR = "/tmp/kubed/snapshot"
)

func NewMinioBackend(bucket, prefix, endpoint, secretName string) *api.Backend {
	return &api.Backend{
		S3: &api.S3Spec{
			Bucket:   bucket,
			Prefix:   prefix,
			Endpoint: endpoint,
		},
		StorageSecretName: secretName,
	}
}

func NewLocalBackend(dir string) *api.Backend {
	return &api.Backend{
		Local: &api.LocalSpec{
			Path: dir,
		},
	}
}

func (f *Invocation) EventuallyBackupSnapshot(backend api.Backend) GomegaAsyncAssertion {
	return Eventually(func() []stow.Item {
		loc, err := f.GetLocation(backend)
		Expect(err).NotTo(HaveOccurred())

		bucket, prefix, err := backend.GetBucketAndPrefix()
		Expect(err).NotTo(HaveOccurred())
		if backend.Local == nil {
			prefix = prefix + "/"
		}

		container, err := loc.Container(bucket)
		Expect(err).NotTo(HaveOccurred())

		items, _, err := container.Items(prefix, stow.CursorStart, 50)
		Expect(err).NotTo(HaveOccurred())

		return items
	})
}

func (f *Invocation) CreateBucketIfNotExist(backend api.Backend) error {
	err := storage.CheckBucketAccess(f.KubeClient, backend, f.namespace)
	if err != nil {
		if err.Error() == stow.ErrNotFound.Error() {
			loc, err := f.GetLocation(backend)
			if err != nil {
				return err
			}
			bucket, err := backend.Container()
			if err != nil {
				return err
			}
			_, err = loc.CreateContainer(bucket)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (f *Invocation) GetLocation(backend api.Backend) (stow.Location, error) {
	cfg, err := storage.NewOSMContext(f.KubeClient, backend, f.namespace)
	if err != nil {
		return nil, err
	}

	loc, err := stow.Dial(cfg.Provider, cfg.Config)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
// Ref: https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func Untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	if err != nil {
		return err
	}

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

			// return any other error
		case err != nil:
			return err

			// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case 5:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// if it's a file create it
		case 0:
			err := os.MkdirAll(filepath.Dir(target), 0755)
			if err != nil {
				return err
			}

			f, err := os.Create(target)
			if err != nil {
				return err
			}
			defer f.Close()

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

func ReadYaml(fileName string) (*apps.Deployment, error) {
	dpl := &apps.Deployment{}

	err := filepath.Walk(TEST_LOCAL_BACKUP_DIR, func(path string, info os.FileInfo, err error) error {
		Expect(err).NotTo(HaveOccurred())
		if !info.IsDir() && info.Name() == fileName {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			err = yaml.Unmarshal(data, dpl)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dpl, nil
}

func DeploymentSnapshotSanitized(dpl *apps.Deployment) error {

	if err := ObjectMetaCleaned(dpl.ObjectMeta); err != nil {
		return err
	}

	if err := PodSpecCleaned(dpl.Spec.Template.Spec); err != nil {
		return err
	}

	if dpl.Status.ObservedGeneration != 0 ||
		dpl.Status.ReadyReplicas != 0 ||
		dpl.Status.Replicas != 0 ||
		dpl.Status.AvailableReplicas != 0 ||
		dpl.Status.UnavailableReplicas != 0 ||
		dpl.Status.UpdatedReplicas != 0 ||
		dpl.Status.Conditions != nil ||
		dpl.Status.CollisionCount != nil {
		return errors.New(`"status" field not cleaned up`)
	}

	return nil
}

func ObjectMetaCleaned(meta metav1.ObjectMeta) error {
	if !meta.CreationTimestamp.IsZero() {
		return errors.New(`"creationTimestamp" not cleaned up`)
	}

	if meta.ResourceVersion != "" {
		return errors.New(`"resourceVersion" not cleaned up`)
	}

	if meta.UID != "" {
		return errors.New(`"uid" not cleaned up`)
	}

	if meta.GenerateName != "" {
		return errors.New(`"generateName" not cleaned up`)
	}

	if meta.Generation != 0 {
		return errors.New(`"generation" not cleaned up`)
	}

	if err := DecoratorCleaned(meta); err != nil {
		return err
	}
	return nil
}

func DecoratorCleaned(meta metav1.ObjectMeta) error {
	if metav1.HasAnnotation(meta, "controller-uid") {
		return errors.New(`"controller-uid" not cleaned up`)
	}

	if metav1.HasAnnotation(meta, "deployment.kubernetes.io/desired-replicas") {
		return errors.New(`"deployment.kubernetes.io/desired-replicas" not cleaned up`)
	}

	if metav1.HasAnnotation(meta, "deployment.kubernetes.io/max-replicas") {
		return errors.New(`"deployment.kubernetes.io/max-replicas" not cleaned up`)
	}

	// currently revision cleaned up not working //TODO: fix revision cleanup
	//if metav1.HasAnnotation(meta, "deployment.kubernetes.io/revision") {
	//	return errors.New(`"deployment.kubernetes.io/revision" not cleaned up`)
	//}

	if metav1.HasAnnotation(meta, "pod-template-hash") {
		return errors.New(`"pod-template-hash" not cleaned up`)
	}

	if metav1.HasAnnotation(meta, "pv.kubernetes.io/bind-completed") {
		return errors.New(`"pv.kubernetes.io/bind-completed" not cleaned up`)
	}

	if metav1.HasAnnotation(meta, "pv.kubernetes.io/bound-by-controller") {
		return errors.New(`"pv.kubernetes.io/bound-by-controller" not cleaned up`)
	}
	return nil
}

func PodSpecCleaned(spec v1.PodSpec) error {
	if spec.DNSPolicy != "" {
		return errors.New(`"dns policy" not cleaned up`)
	}

	if spec.NodeName != "" {
		return errors.New(`"nodeName" not cleaned up`)
	}

	if spec.ServiceAccountName == "default" {
		return errors.New(`"default" service account not removed`)
	}

	if spec.TerminationGracePeriodSeconds != nil {
		return errors.New(`"terminationGracePeriodSeconds" not cleaned up`)
	}

	for _, c := range spec.Containers {
		if c.TerminationMessagePath != "" {
			return errors.New(`"terminationMessagePath" not removed from container`)
		}
	}

	for _, c := range spec.InitContainers {
		if c.TerminationMessagePath != "" {
			return errors.New(`"terminationMessagePath" not removed from init container`)
		}
	}

	return nil
}
