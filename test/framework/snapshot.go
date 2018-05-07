package framework

import (
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/pkg/storage"
	"github.com/graymeta/stow"
	. "github.com/onsi/gomega"
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

func (f *Invocation) EventuallyBackupSnapshot(backend api.Backend) GomegaAsyncAssertion {
	return Eventually(func() []stow.Item {
		loc, err := f.GetLocation(backend)
		Expect(err).NotTo(HaveOccurred())

		bucket, prefix, err := backend.GetBucketAndPrefix()
		Expect(err).NotTo(HaveOccurred())
		prefix = prefix + "/"

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
