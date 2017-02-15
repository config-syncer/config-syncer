package event

import (
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/k8s-addons/pkg/events"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/controller/types"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

func CreateAlertEvent(kubeClient clientset.Interface, alert *aci.Alert, reason types.EventReason, additionalMessage ...string) {
	timestamp := unversioned.NewTime(time.Now().UTC())
	event := &kapi.Event{
		ObjectMeta: kapi.ObjectMeta{
			Name:      rand.WithUniqSuffix("alert"),
			Namespace: alert.Namespace,
		},
		InvolvedObject: kapi.ObjectReference{
			Kind:      events.ObjectKindAlert.String(),
			Namespace: alert.Namespace,
			Name:      alert.Name,
		},
		Source: kapi.EventSource{
			Component: "searchlight",
		},

		Count:          1,
		FirstTimestamp: timestamp,
		LastTimestamp:  timestamp,
	}

	switch reason {
	case types.CreatingIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`creating Icinga objects`)
		event.Type = kapi.EventTypeNormal
	case types.FailedToCreateIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to create Icinga objects. Error: %v`, additionalMessage)
		event.Type = kapi.EventTypeWarning
	case types.NoIcingaObjectCreated:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`no Icinga object is created. Reason: %v`, additionalMessage)
		event.Type = kapi.EventTypeNormal
	case types.CreatedIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully created Icinga objects`)
		event.Type = kapi.EventTypeNormal

	case types.UpdatingIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`updating Icinga objects`)
	case types.FailedToUpdateIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to update Icinga objects. Error: %v`, additionalMessage)
		event.Type = kapi.EventTypeWarning
	case types.UpdatedIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully updated Icinga objects.`)
		event.Type = kapi.EventTypeNormal

	case types.DeletingIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`deleting Icinga objects`)
		event.Type = kapi.EventTypeNormal
	case types.FailedToDeleteIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to delete Icinga objects. Error: %v`, additionalMessage)
		event.Type = kapi.EventTypeWarning
	case types.DeletedIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully deleted Icinga objects.`)
		event.Type = kapi.EventTypeNormal

	case types.SyncIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`synchronizing alert for %v.`, additionalMessage[0])
		event.Type = kapi.EventTypeNormal
	case types.FailedToSyncIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to synchronize alert for %v. Error: %v`, additionalMessage[0], additionalMessage[1])
		event.Type = kapi.EventTypeWarning
	case types.SyncedIcingaObjects:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully synchronized alert for %v.`, additionalMessage[0])
		event.Type = kapi.EventTypeNormal
	}

	if _, err := kubeClient.Core().Events(alert.Namespace).Create(event); err != nil {
		log.Debugln(err)
	}
}
