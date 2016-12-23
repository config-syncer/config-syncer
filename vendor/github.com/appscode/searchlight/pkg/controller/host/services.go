package host

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/appscode/errors"
	"github.com/appscode/searchlight/pkg/client/icinga"
)

const (
	HostTypeLocalhost = "localhost"
	HostTypeNode      = "node"
	HostTypePod       = "pod"
)

// createIcingaServiceForCluster
func CreateIcingaService(icingaClient *icinga.IcingaClient, mp map[string]interface{}, object *KubeObjectInfo, serviceName string) error {
	var obj IcingaObject
	obj.Templates = []string{"generic-service"}
	obj.Attrs = mp
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.New().WithCause(err).Internal()
	}

	resp := icingaClient.Objects().Service(object.Name).Create([]string{serviceName}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.New().WithCause(resp.Err).Internal()
	}

	if resp.Status == 200 {
		return nil
	}
	if strings.Contains(string(resp.ResponseBody), "already exists") {
		return nil
	}

	return errors.New().WithMessage("Can't create Icinga service").Failed()
}

func UpdateIcingaService(icingaClient *icinga.IcingaClient, mp map[string]interface{}, object *KubeObjectInfo, icignaService string) error {
	var obj IcingaObject
	obj.Templates = []string{"generic-service"}
	obj.Attrs = mp
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.New().WithCause(err).Internal()
	}
	resp := icingaClient.Objects().Service(object.Name).Update([]string{icignaService}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.New().WithCause(resp.Err).Internal()
	}

	if resp.Status != 200 {
		return errors.New().WithMessage("Can't update Icinga service").Failed()
	}
	return nil
}

func DeleteIcingaService(icingaClient *icinga.IcingaClient, objectList []*KubeObjectInfo, icingaServiceName string) error {
	param := map[string]string{
		"cascade": "1",
	}
	in := IcingaServiceSearchQuery(icingaServiceName, objectList)
	resp := icingaClient.Objects().Service("").Delete([]string{}, in).Params(param).Do()

	if resp.Err != nil {
		return errors.New().WithCause(resp.Err).Internal()
	}
	if resp.Status == 200 {
		return nil
	}
	return errors.New().WithMessage("Fail to delete service").Failed()
}

func CheckIcingaService(icingaClient *icinga.IcingaClient, icingaServiceName string, objectList []*KubeObjectInfo) (bool, error) {
	in := IcingaServiceSearchQuery(icingaServiceName, objectList)
	var respService ResponseObject

	if _, err := icingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
		return true, errors.New().WithMessage("can't check icinga service").Failed()
	}
	return len(respService.Results) > 0, nil
}

func IcingaServiceSearchQuery(icingaServiceName string, objectList []*KubeObjectInfo) string {
	matchHost := ""
	for id, object := range objectList {
		if id > 0 {
			matchHost = matchHost + "||"
		}
		matchHost = matchHost + fmt.Sprintf(`match(\"%s\",host.name)`, object.Name)
	}
	return fmt.Sprintf(`{"filter": "(%s)&&match(\"%s\",service.name)"}`, matchHost, icingaServiceName)
}
