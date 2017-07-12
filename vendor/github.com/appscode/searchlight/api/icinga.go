package api

import (
	"fmt"

	"github.com/appscode/log"
	"github.com/appscode/searchlight/data"
)

type CheckPod string

const (
	CheckPodInfluxQuery      CheckPod = "influx_query"
	CheckPodStatus           CheckPod = "pod_status"
	CheckPodPrometheusMetric CheckPod = "prometheus_metric"
	CheckVolume              CheckPod = "volume"
	CheckPodExec             CheckPod = "kube_exec"
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
	CheckNodeInfluxQuery      CheckNode = "influx_query"
	CheckNodeDisk             CheckNode = "node_disk"
	CheckNodeStatus           CheckNode = "node_status"
	CheckNodePrometheusMetric CheckNode = "prometheus_metric"
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
	CheckHttp             CheckCluster = "any_http"
	CheckComponentStatus  CheckCluster = "component_status"
	CheckJsonPath         CheckCluster = "json_path"
	CheckNodeCount        CheckCluster = "node_count"
	CheckPodExists        CheckCluster = "pod_exists"
	CheckPrometheusMetric CheckCluster = "prometheus_metric"
	CheckClusterEvent     CheckCluster = "kube_event"
	CheckHelloIcinga      CheckCluster = "hello_icinga"
	CheckDIG              CheckCluster = "dig"
	CheckDNS              CheckCluster = "dns"
	CheckDummy            CheckCluster = "dummy"
	CheckICMP             CheckCluster = "icmp"
)

func (c CheckCluster) IsValid() bool {
	_, ok := ClusterCommands[c]
	return ok
}

func ParseCheckCluster(s string) (*CheckCluster, error) {
	cmd := CheckCluster(s)
	if _, ok := ClusterCommands[cmd]; !ok {
		return nil, fmt.Errorf("Invalid cluster check command %s", s)
	}
	return &cmd, nil
}

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
	PodCommands = map[CheckPod]IcingaCommand{}
	NodeCommands = map[CheckNode]IcingaCommand{}
	ClusterCommands = map[CheckCluster]IcingaCommand{}

	cmdList, err := data.LoadIcingaData()
	if err != nil {
		log.Fatal(err)
	}
	for _, cmd := range cmdList.Command {
		vars := make(map[string]data.CommandVar)
		for _, v := range cmd.Vars {
			vars[v.Name] = v
		}
		c := IcingaCommand{
			Name:   cmd.Name,
			Vars:   vars,
			States: cmd.States,
		}
		if c.Name == "influx_query" ||
			c.Name == "pod_status" ||
			c.Name == "prometheus_metric" ||
			c.Name == "volume" ||
			c.Name == "kube_exec" {
			PodCommands[CheckPod(c.Name)] = c
		}
		if c.Name == "influx_query" ||
			c.Name == "node_disk" ||
			c.Name == "node_status" ||
			c.Name == "prometheus_metric" {
			NodeCommands[CheckNode(c.Name)] = c
		}
		if c.Name == "any_http" ||
			c.Name == "component_status" ||
			c.Name == "json_path" ||
			c.Name == "node_exists" ||
			c.Name == "pod_exists" ||
			c.Name == "kube_event" ||
			c.Name == "hello_icinga" ||
			c.Name == "dig" ||
			c.Name == "dns" ||
			c.Name == "dummy" ||
			c.Name == "icmp" {
			ClusterCommands[CheckCluster(c.Name)] = c
		}
	}
}
