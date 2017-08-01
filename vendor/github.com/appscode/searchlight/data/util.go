package data

import (
	"encoding/json"

	"github.com/appscode/searchlight/data/files"
)

func LoadClusterChecks() (ic IcingaData, err error) {
	bytes, err := files.Asset("cluster_checks.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &ic)
	return
}

func LoadNodeChecks() (ic IcingaData, err error) {
	bytes, err := files.Asset("node_checks.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &ic)
	return
}

func LoadPodChecks() (ic IcingaData, err error) {
	bytes, err := files.Asset("pod_checks.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &ic)
	return
}
