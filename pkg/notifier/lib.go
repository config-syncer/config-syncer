package notifier

import (
	"errors"
	"fmt"

	"github.com/appscode/kubed/pkg/notifier/extpoints"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

const (
	Notify_Via string = "notify_via"
)

func SendNotification(notification string) error {
	driver, _ := notificationDriver(nil, "", "")

	if err := driver.Notify(notification); err != nil {
		return err
	}
	return driver.Notify(notification)
}

func notificationDriver(client clientset.Interface, secretName string, secretNamespace string) (extpoints.Driver, error) {
	clusterConf, err := client.Core().
		Secrets(secretNamespace).
		Get(secretName, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return notifierDriverFromConfiguration(clusterConf.Data)
}

func notifierDriverFromConfiguration(cnfData map[string][]byte) (extpoints.Driver, error) {
	data := make(map[string]string, len(cnfData))
	for key, val := range cnfData {
		data[key] = string(val)
	}
	via, ok := data[Notify_Via]
	if !ok {
		return nil, errors.New("No notifier set")
	}
	driver := extpoints.Drivers.Lookup(via)
	if driver == nil {
		return nil, errors.New(fmt.Sprintf("Notifier `%v` not found", via))
	}
	err := driver.SetOptions(data)
	if err != nil {
		return nil, err
	}
	return driver, nil
}
