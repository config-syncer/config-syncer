package pkg

import (
	"fmt"
	"os"
	"testing"
	"time"

	"appscode.com/kubed/pkg/watcher"
	"github.com/appscode/client"
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/flags"
	aci "github.com/appscode/k8s-addons/api"
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/mikespook/golib/signal"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/selection"
	"k8s.io/kubernetes/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/util/sets"
)

func TestExtIngressList(t *testing.T) {
	flags.SetLogLevel(5)
	k := testKubernetesWithoutApiClient()
	log.Infoln("running test")

	opts := api.ListOptions{}
	list, err := k.AppsCodeExtensionClient.Ingress(api.NamespaceAll).List(opts)
	fmt.Println("Error", err)
	for _, l := range list.Items {
		fmt.Println(l.Name, l.Kind, l.APIVersion, l.Labels)
	}

	l, err := k.AppsCodeExtensionClient.Ingress(api.NamespaceDefault).Get("appscode-ext-service")
	fmt.Println(err, l.Name)
}

func TestExtIngressCreate(t *testing.T) {
	flags.SetLogLevel(7)
	k := testKubernetesWithoutApiClient()
	log.Infoln("running test")

	eng := &aci.Ingress{
		ObjectMeta: api.ObjectMeta{
			Name:      "engress-from-" + rand.Characters(4),
			Namespace: "default",
			Labels: map[string]string{
				"owner": "sadlil",
			},
		},

		Spec: aci.ExtendedIngressSpec{
			Backend: &aci.ExtendedIngressBackend{
				ServiceName: "no-service",
				ServicePort: intstr.FromString("2323"),
				RewriteRule: []string{"a*b.com"},
			},
		},
	}

	k.AppsCodeExtensionClient.Ingress("default").Create(eng)
}

func TestDelete(t *testing.T) {
	flags.SetLogLevel(7)
	k := testKubernetesWithoutApiClient()
	log.Infoln("running test")

	fmt.Println(k.AppsCodeExtensionClient.Ingress("default").Delete("engress-from-ebb4", nil))
}

func TestDaemon(t *testing.T) {
	flags.SetLogLevel(6)
	k := testKubernetesWithoutApiClient()

	k.Run()

	signal.Bind(os.Interrupt, func() uint { return signal.BreakExit })
	signal.Wait()
}

func TestAlertD(t *testing.T) {
	flags.SetLogLevel(1)
	k := testKubernetesWithoutApiClient()

	log.Infoln("running test")

	icingaConfig, err := readConfig("/srv/icinga2/secrets/.env")
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	icingaClient := icinga.NewClient(fmt.Sprintf("https://%v:5665/v1", ""), icingaConfig[IcingaAPIUser], icingaConfig[IcingaAPIPass], nil)

	k.IcingaClient = icingaClient
	k.Run()

	signal.Bind(os.Interrupt, func() uint { return signal.BreakExit })
	signal.Wait()
}

func TestCodec(t *testing.T) {
	l := &api.ListOptions{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ListOptions",
			APIVersion: "v1",
		},
		LabelSelector: labels.NewSelector(),
	}

	fmt.Println(l.GroupVersionKind().GroupVersion())

	req, _ := labels.NewRequirement("hello", selection.Equals, sets.NewString("world").List())
	l.LabelSelector = l.LabelSelector.Add(*req)

	x, err := runtime.NewParameterCodec(api.Scheme).EncodeParameters(l, l.GroupVersionKind().GroupVersion())
	fmt.Println(x, err)
}

func TestDeleteAllAlert(t *testing.T) {
	flags.SetLogLevel(7)
	k := testKubernetesWithoutApiClient()
	log.Infoln("running test")

	ns, _ := k.Client.Core().Namespaces().List(api.ListOptions{LabelSelector: labels.Everything()})

	for _, n := range ns.Items {
		list, _ := k.AppsCodeExtensionClient.Alert(n.Name).List(api.ListOptions{LabelSelector: labels.Everything()})

		for _, l := range list.Items {
			fmt.Println("--- ", l.Name)
			k.AppsCodeExtensionClient.Alert(n.Name).Delete(l.Name, nil)
		}
	}
}

func TestIcinga(t *testing.T) {
	icingaConfig, err := readConfig("/srv/icinga2/secrets/.env")
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	icingaService := icingaConfig[IcingaService]
	if icingaService == "" {
		icingaService = "appscode-alert"
	}
	icingaClient := icinga.NewClient(fmt.Sprintf("https://%v:5665/v1", ""), icingaConfig[IcingaAPIUser], icingaConfig[IcingaAPIPass], nil)

	resp := icingaClient.Check().Get([]string{}).Do()

	fmt.Println("-- ", resp.Err, resp.Status)
}

func TestGetAlert(t *testing.T) {
	w := testKubernetesWithoutApiClient()

	l, err := w.AppsCodeExtensionClient.Alert("").List(api.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	for _, item := range l.Items {
		fmt.Println(item, item.Name)
	}

	item, err := w.AppsCodeExtensionClient.Alert("default").Get("cluster-pod-status")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(*item, "|", item.Name)
}

func testKubernetes() *watcher.KubedWatcher {
	config := &restclient.Config{
		Host:     "https://54.82.21.25:6443",
		Username: "admin@h-505-qacode.appscode.xyz",
		Password: "hkesTYFbs8DOkiTq",
	}
	c := clientset.NewForConfigOrDie(config)
	ops := client.NewOption("localhost:50051")
	ops.BearerAuth("appscode", "api-T3stT0ken")
	ac := acs.NewACExtensionsForConfigOrDie(config)

	kubeWatcher := &watcher.KubedWatcher{
		Watcher: acw.Watcher{
			Client:                  c,
			AppsCodeExtensionClient: ac,
			SyncPeriod:              time.Second * 10,
		},
		AppsCodeApiClientOptions: ops,
	}
	return kubeWatcher
}

func testKubernetesWithoutApiClient() *watcher.KubedWatcher {
	config := &restclient.Config{
		Host:     "https://54.82.21.25:6443/",
		Username: "admin@h-505-qacode.appscode.xyz",
		Password: "hkesTYFbs8DOkiTq",
		Insecure: true,
	}

	c := clientset.NewForConfigOrDie(config)
	ac := acs.NewACExtensionsForConfigOrDie(config)

	kubeWatcher := &watcher.KubedWatcher{
		Watcher: acw.Watcher{
			Client:                  c,
			AppsCodeExtensionClient: ac,
			SyncPeriod:              time.Second * 10,
		},
	}
	return kubeWatcher
}
