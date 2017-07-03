package notifier

import (
	"errors"
	"fmt"

	"github.com/appscode/kubed/pkg/notifier/extpoints"
)

const (
	Notify_Via string = "notify_via"
)

type Notifier struct {
	configuration map[string]string
}

func New(conf map[string]string) Notifier {
	return Notifier{
		configuration: conf,
	}
}

func (n Notifier) SendNotification(notification string) (string, error) {
	driver, err := n.notificationDriver()
	if err != nil {
		return "", err
	}
	return driver.Uid(), driver.Notify(notification)
}

func (n Notifier) notificationDriver() (extpoints.Driver, error) {
	via, ok := n.configuration[Notify_Via]
	if !ok {
		return nil, errors.New("No notifier set")
	}
	driver := extpoints.Drivers.Lookup(via)
	if driver == nil {
		return nil, errors.New(fmt.Sprintf("Notifier `%v` not found", via))
	}
	err := driver.SetOptions(n.configuration)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func EnsureRequiredKeys(mp map[string]string, keys []string) error {
	for _, k := range keys {
		if _, found := mp[k]; !found {
			return errors.New(fmt.Sprintf("%v not found", k))
		}
	}
	return nil
}
