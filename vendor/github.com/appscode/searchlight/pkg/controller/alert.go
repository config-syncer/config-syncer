package controller

import (
	"fmt"
	"time"

	"github.com/appscode/errors"
	aci "github.com/appscode/k8s-addons/api"
	acs "github.com/appscode/k8s-addons/client/clientset"
	"github.com/appscode/k8s-addons/pkg/events"
	"github.com/appscode/k8s-addons/pkg/stash"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/controller/host"
	_ "github.com/appscode/searchlight/pkg/controller/host/localhost"
	_ "github.com/appscode/searchlight/pkg/controller/host/node"
	_ "github.com/appscode/searchlight/pkg/controller/host/pod"
	"github.com/appscode/searchlight/pkg/controller/types"
	kapi "k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/labels"
)

type IcingaController struct {
	ctx *types.Context
}

func New(kubeClient clientset.Interface,
	icingaClient *icinga.IcingaClient,
	appsCodeExtensionClient acs.AppsCodeExtensionInterface,
	storage *stash.Storage) *IcingaController {
	data, err := getIcingaDataMap()
	if err != nil {
		log.Errorln("Icinga data not found")
	}
	ctx := &types.Context{
		KubeClient:              kubeClient,
		AppsCodeExtensionClient: appsCodeExtensionClient,
		IcingaData:              data,
		IcingaClient:            icingaClient,
		Storage:                 storage,
	}
	return &IcingaController{ctx: ctx}
}

func (b *IcingaController) Handle(e *events.Event) error {
	var err error
	switch e.ResourceType {
	case events.Alert:
		err = b.handleAlert(e)
	case events.Pod:
		err = b.handlePod(e)
	case events.Node:
		err = b.handleNode(e)
	case events.Service:
		err = b.handleService(e)
	}

	if err != nil {
		log.Debugln(err)
	}

	return nil
}

func (b *IcingaController) handleAlert(e *events.Event) error {
	var alert []interface{}
	if e.ResourceType == events.Alert {
		alert = e.RuntimeObj
	}

	if e.EventType.IsAdded() {
		if len(alert) == 0 {
			return errors.New().WithMessage("Missing alert data").NotFound()
		}
		b.ctx.Resource = alert[0].(*aci.Alert)

		if err := b.IsObjectExists(); err != nil {
			return errors.New().WithCause(err).Internal()
		}
		return b.Create()
	} else if e.EventType.IsUpdated() {
		if len(alert) == 0 {
			return errors.New().WithMessage("Missing alert data").NotFound()
		}

		oldConfig := alert[0].(*aci.Alert)
		newConfig := alert[1].(*aci.Alert)

		if err := host.CheckAlertConfig(oldConfig, newConfig); err != nil {
			return errors.New().WithCause(err).BadRequest()
		}

		b.ctx.Resource = newConfig
		if err := b.IsObjectExists(); err != nil {
			return errors.New().WithCause(err).Internal()
		}

		return b.Update()
	} else if e.EventType.IsDeleted() {
		if len(alert) == 0 {
			return errors.New().WithMessage("Missing alert data").NotFound()
		}
		b.ctx.Resource = alert[0].(*aci.Alert)
		b.parseAlertOptions()
		return b.Delete()
	}
	return nil
}

func (b *IcingaController) handlePod(e *events.Event) error {
	if !(e.EventType.IsAdded() || e.EventType.IsDeleted()) {
		return nil
	}
	if host.IsIcingaApp(e.MetaData.Labels) {
		if e.EventType.IsAdded() {
			go b.handleIcingaPod()
		}
	} else {
		return b.handleRegularPod(e)
	}

	return nil
}

func (b *IcingaController) handleIcingaPod() {
	log.Debugln("Icinga pod is created...")
	then := time.Now()
	for {
		log.Debugln("Waiting for Icinga to UP")
		if b.checkIcingaAvailability() {
			break
		}
		now := time.Now()
		if now.Sub(then) > time.Minute*10 {
			log.Debugln("Icinga is down for more than 10 minutes..")
			return
		}
		time.Sleep(time.Second * 30)
	}

	namespaces, _ := b.ctx.KubeClient.Core().Namespaces().List(kapi.ListOptions{LabelSelector: labels.Everything()})
	icingaUp := false
	for _, ns := range namespaces.Items {
		alertList, err := b.ctx.AppsCodeExtensionClient.Alert(ns.Name).List(kapi.ListOptions{LabelSelector: labels.Everything()})
		if err != nil {
			log.Errorln(err)
			return
		}

		log.Debugln(fmt.Sprintf("Applying %v alert for namespace %s", len(alertList.Items), ns.Name))

		for _, alert := range alertList.Items {
			if !icingaUp && !b.checkIcingaAvailability() {
				log.Debugln("Icinga is down...")
				return
			}
			icingaUp = true

			fakeEvent := &events.Event{
				ResourceType: events.Alert,
				EventType:    events.Added,
				RuntimeObj:   make([]interface{}, 0),
			}
			fakeEvent.RuntimeObj = append(fakeEvent.RuntimeObj, &alert)

			if err := b.handleAlert(fakeEvent); err != nil {
				log.Debugln(err)
			}
		}
	}
	return
}

func (b *IcingaController) handleRegularPod(e *events.Event) error {
	namespace := e.MetaData.Namespace
	ancestors := b.getParentsForPod(e.RuntimeObj[0])
	icingaUp := false
	for _, ancestor := range ancestors {
		objectType := ancestor.Type
		for _, objectName := range ancestor.Names {
			lb, err := host.GetLabelSelector(objectType, objectName)
			if err != nil {
				return errors.New().WithCause(err).Internal()
			}

			alertList, err := b.ctx.AppsCodeExtensionClient.Alert(namespace).List(kapi.ListOptions{
				LabelSelector: lb,
			})
			if err != nil {
				return errors.New().WithCause(err).Internal()
			}

			for _, alert := range alertList.Items {
				if !icingaUp && !b.checkIcingaAvailability() {
					return errors.New("Icinga is down").External()
				}
				icingaUp = true

				if e.EventType.IsAdded() {
					// Waiting for POD IP to use as Icinga Host IP
					then := time.Now()
					for {
						hasPodIP, err := b.checkPodIPAvailability(e.MetaData.Name, namespace)
						if err != nil {
							return errors.New().WithCause(err).Internal()
						}
						if hasPodIP {
							break
						}
						log.Debugln("Waiting for pod IP")
						now := time.Now()
						if now.Sub(then) > time.Minute*2 {
							return errors.New("Pod IP is not available for 2 minutes").Internal()
						}
						time.Sleep(time.Second * 30)
					}

					b.ctx.Resource = &alert
					b.parseAlertOptions()
					b.Create(e.MetaData.Name)
				} else if e.EventType.IsDeleted() {
					b.ctx.Resource = &alert
					b.parseAlertOptions()
					b.Delete(e.MetaData.Name)
				}
			}
		}
	}
	return nil
}

func (b *IcingaController) handleNode(e *events.Event) error {
	if !(e.EventType.IsAdded() || e.EventType.IsDeleted()) {
		return nil
	}

	lb, err := host.GetLabelSelector(events.Cluster.String(), "")
	if err != nil {
		return errors.New().WithCause(err).Internal()
	}
	lb1, err := host.GetLabelSelector(events.Node.String(), e.MetaData.Name)
	if err != nil {
		return errors.New().WithCause(err).Internal()
	}

	requirements, _ := lb1.Requirements()
	lb.Add(requirements...)

	namespaces, _ := b.ctx.KubeClient.Core().Namespaces().List(kapi.ListOptions{LabelSelector: labels.Everything()})
	icingaUp := false
	for _, ns := range namespaces.Items {
		alertList, err := b.ctx.AppsCodeExtensionClient.Alert(ns.Name).List(kapi.ListOptions{
			LabelSelector: lb,
		})
		if err != nil {
			return errors.New().WithCause(err).Internal()
		}

		for _, alert := range alertList.Items {
			if !icingaUp && !b.checkIcingaAvailability() {
				return errors.New("Icinga is down").External()
			}
			icingaUp = true

			if e.EventType.IsAdded() {
				b.ctx.Resource = &alert
				b.parseAlertOptions()
				b.Create(e.MetaData.Name)
			} else if e.EventType.IsDeleted() {
				b.ctx.Resource = &alert
				b.parseAlertOptions()
				b.Delete(e.MetaData.Name)
			}
		}
	}

	return nil
}

func (b *IcingaController) handleService(e *events.Event) error {
	if e.EventType.IsAdded() {
		if checkIcingaService(e.MetaData.Name, e.MetaData.Namespace) {
			service, err := b.ctx.KubeClient.Core().Services(e.MetaData.Namespace).Get(e.MetaData.Name)
			if err != nil {
				return errors.New().WithCause(err).Internal()
			}
			endpoint := fmt.Sprintf("https://%v:5665/v1", service.Spec.ClusterIP)
			b.ctx.IcingaClient = b.ctx.IcingaClient.SetEndpoint(endpoint)
		}
	}
	return nil
}

func getIcingaDataMap() (map[string]*types.IcingaData, error) {
	icingaData, err := data.LoadIcingaData()
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	icingaDataMap := make(map[string]*types.IcingaData)
	for _, command := range icingaData.Command {
		varsMap := make(map[string]data.CommandVar)
		for _, v := range command.Vars {
			varsMap[v.Name] = v
		}

		icingaDataMap[command.Name] = &types.IcingaData{
			HostType: command.ObjectToHost,
			VarInfo:  varsMap,
		}
	}
	return icingaDataMap, nil
}
