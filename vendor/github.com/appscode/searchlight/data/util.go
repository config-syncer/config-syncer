package data

import (
	"encoding/json"

	"github.com/appscode/searchlight/data/files"
)

func LoadIcingaData() (ic IcingaData, err error) {
	bytes, err := files.Asset("icinga.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &ic)
	return
}
