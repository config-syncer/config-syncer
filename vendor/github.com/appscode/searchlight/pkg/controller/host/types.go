package host

const (
	internalIP = "InternalIP"
)

type KubeObjectInfo struct {
	Name      string
	IP        string
	GroupName string
	GroupType string
}

type IcingaObject struct {
	Templates []string               `json:"templates,omitempty"`
	Attrs     map[string]interface{} `json:"attrs"`
}

type ResponseObject struct {
	Results []struct {
		Attrs struct {
			Name          string                 `json:"name"`
			CheckInterval float64                `json:"check_interval"`
			Vars          map[string]interface{} `json:"vars"`
		} `json:"attrs"`
		Name string `json:"name"`
	} `json:"results"`
}

func IVar(value string) string {
	return "vars." + value
}
