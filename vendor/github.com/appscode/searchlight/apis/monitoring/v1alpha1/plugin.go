package v1alpha1

import (
	"fmt"

	"github.com/appscode/kutil/meta"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindSearchlightPlugin     = "SearchlightPlugin"
	ResourcePluralSearchlightPlugin   = "searchlightplugins"
	ResourceSingularSearchlightPlugin = "searchlightplugin"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SearchlightPlugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the SearchlightPlugin.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec SearchlightPluginSpec `json:"spec,omitempty"`
}

// SearchlightPluginSpec describes the SearchlightPlugin the user wishes to create.
type SearchlightPluginSpec struct {
	// Check Command
	Command string `json:"command,omitempty"`

	// Webhook provides a reference to the service for this SearchlightPlugin.
	// It must communicate on port 80
	Webhook *WebhookServiceSpec `json:"webhook,omitempty"`

	// AlertKinds refers to supports Alert kinds for this plugin
	AlertKinds []string `json:"alertKinds"`
	// Supported arguments for SearchlightPlugin
	Arguments PluginArguments `json:"arguments,omitempty"`
	// Supported Icinga Service State
	State []string `json:"state"`
}

type WebhookServiceSpec struct {
	// Namespace is the namespace of the service
	Namespace string `json:"namespace,omitempty"`
	// Name is the name of the service
	Name string `json:"name"`
}

type VarType string

const (
	VarTypeInteger  VarType = "integer"
	VarTypeNumber   VarType = "number"
	VarTypeBoolean  VarType = "boolean"
	VarTypeString   VarType = "string"
	VarTypeDuration VarType = "duration"
)

type PluginVarItem struct {
	Description string  `json:"description,omitempty"`
	Type        VarType `json:"type"`
}

type PluginVars struct {
	Items    map[string]PluginVarItem `json:"Item"`
	Required []string                 `json:"required,omitempty"`
}

type PluginArguments struct {
	Vars *PluginVars       `json:"vars,omitempty"`
	Host map[string]string `json:"host,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SearchlightPluginList is a collection of SearchlightPlugin.
type SearchlightPluginList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of SearchlightPlugin.
	Items []SearchlightPlugin `json:"items"`
}

var (
	validateVarValue = map[VarType]meta.ParserFunc{}
)

func registerVarValueParser(key VarType, fn meta.ParserFunc) {
	validateVarValue[key] = fn
}

func init() {
	registerVarValueParser(VarTypeInteger, meta.GetInt)
	registerVarValueParser(VarTypeNumber, meta.GetFloat)
	registerVarValueParser(VarTypeBoolean, meta.GetBool)
	registerVarValueParser(VarTypeString, meta.GetString)
	registerVarValueParser(VarTypeDuration, meta.GetDuration)
}

func validateVariables(pluginVars *PluginVars, vars map[string]string) error {
	if pluginVars == nil {
		return nil
	}
	// Check if any invalid variable is provided
	var err error
	for k := range vars {
		p, found := pluginVars.Items[k]
		if !found {
			return fmt.Errorf("var '%s' is unsupported", k)
		}

		fn, found := validateVarValue[p.Type]
		if !found {
			return errors.Errorf(`type "%v" is not registered`, p.Type)
		}
		if _, err = fn(vars, k); err != nil {
			return errors.Wrapf(err, `validation failure: variable "%s" must be of type %v`, k, p.Type)
		}
	}
	for _, k := range pluginVars.Required {
		if _, ok := vars[k]; !ok {
			return fmt.Errorf("plugin variable '%s' is required", k)
		}
	}

	return nil
}
