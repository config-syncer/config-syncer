package v1alpha1

import (
	"fmt"

	"github.com/appscode/go/log"
	"github.com/appscode/searchlight/data"
)

type CheckPod string

const (
	CheckPodInfluxQuery CheckPod = "influx_query"
	CheckPodStatus      CheckPod = "pod_status"
	CheckPodVolume      CheckPod = "pod_volume"
	CheckPodExec        CheckPod = "pod_exec"
)

func (c CheckPod) IsValid() bool {
	_, ok := PodCommands[c]
	return ok
}

func ParseCheckPod(s string) (*CheckPod, error) {
	cmd := CheckPod(s)
	if _, ok := PodCommands[cmd]; !ok {
		return nil, fmt.Errorf("Invalid pod check command %s", s)
	}
	return &cmd, nil
}

type CheckNode string

const (
	CheckNodeInfluxQuery CheckNode = "influx_query"
	CheckNodeVolume      CheckNode = "node_volume"
	CheckNodeStatus      CheckNode = "node_status"
)

func (c CheckNode) IsValid() bool {
	_, ok := NodeCommands[c]
	return ok
}

func ParseCheckNode(s string) (*CheckNode, error) {
	cmd := CheckNode(s)
	if _, ok := NodeCommands[cmd]; !ok {
		return nil, fmt.Errorf("Invalid node check command %s", s)
	}
	return &cmd, nil
}

type CheckCluster string

const (
	CheckComponentStatus CheckCluster = "component_status"
	CheckJsonPath        CheckCluster = "json_path"
	CheckNodeExists      CheckCluster = "node_exists"
	CheckPodExists       CheckCluster = "pod_exists"
	CheckEvent           CheckCluster = "event"
	CheckCACert          CheckCluster = "ca_cert"
	CheckHttp            CheckCluster = "any_http"
	CheckEnv             CheckCluster = "env"
	CheckDummy           CheckCluster = "dummy"
	//CheckICMP          CheckCluster = "icmp"
	//CheckDIG           CheckCluster = "dig"
	//CheckDNS           CheckCluster = "dns"
)

func (c CheckCluster) IsValid() bool {
	_, ok := ClusterCommands[c]
	return ok
}

func ParseCheckCluster(s string) (*CheckCluster, error) {
	cmd := CheckCluster(s)
	if _, ok := ClusterCommands[cmd]; !ok {
		return nil, fmt.Errorf("invalid cluster check command %s", s)
	}
	return &cmd, nil
}

// +k8s:deepcopy-gen=false
// +k8s:gen-deepcopy=false
type IcingaCommand struct {
	Name   string
	Vars   map[string]data.CommandVar
	States []string
}

var (
	PodCommands     map[CheckPod]IcingaCommand
	NodeCommands    map[CheckNode]IcingaCommand
	ClusterCommands map[CheckCluster]IcingaCommand
)

func init() {
	ClusterCommands = map[CheckCluster]IcingaCommand{}
	clusterChecks, err := data.LoadClusterChecks()
	if err != nil {
		log.Fatal(err)
	}
	for _, cmd := range clusterChecks.Command {
		vars := make(map[string]data.CommandVar)
		for _, v := range cmd.Vars {
			vars[v.Name] = v
		}
		ClusterCommands[CheckCluster(cmd.Name)] = IcingaCommand{
			Name:   cmd.Name,
			Vars:   vars,
			States: cmd.States,
		}
	}

	NodeCommands = map[CheckNode]IcingaCommand{}
	nodeChecks, err := data.LoadNodeChecks()
	if err != nil {
		log.Fatal(err)
	}
	for _, cmd := range nodeChecks.Command {
		vars := make(map[string]data.CommandVar)
		for _, v := range cmd.Vars {
			vars[v.Name] = v
		}
		NodeCommands[CheckNode(cmd.Name)] = IcingaCommand{
			Name:   cmd.Name,
			Vars:   vars,
			States: cmd.States,
		}
	}

	PodCommands = map[CheckPod]IcingaCommand{}
	podChecks, err := data.LoadPodChecks()
	if err != nil {
		log.Fatal(err)
	}
	for _, cmd := range podChecks.Command {
		vars := make(map[string]data.CommandVar)
		for _, v := range cmd.Vars {
			vars[v.Name] = v
		}
		PodCommands[CheckPod(cmd.Name)] = IcingaCommand{
			Name:   cmd.Name,
			Vars:   vars,
			States: cmd.States,
		}
	}
}
