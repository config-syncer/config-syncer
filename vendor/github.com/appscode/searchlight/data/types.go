package data

type CommandVar struct {
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Type          string   `json:"type"`
	Format        string   `json:"format,omitempty"`
	Values        []string `json:"values,omitempty"`
	Parameterized bool     `json:"parameterized,omitempty"`
	Flag          struct {
		Long  string `json:"long"`
		Short string `json:"short"`
	} `json:"flag"`
	Optional bool `json:"optional"`
}

type IcingaCheckCommand struct {
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Envs         []string          `json:"envs"`
	ObjectToHost map[string]string `json:"object_to_host"`
	Vars         []CommandVar      `json:"vars,omitempty"`
	States       []string          `json:"states,omitempty"`
}

type IcingaData struct {
	Command []*IcingaCheckCommand `json:"check_command"`
}
