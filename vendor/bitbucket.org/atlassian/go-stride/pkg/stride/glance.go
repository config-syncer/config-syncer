package stride

import (
	"fmt"
)

type GlanceBody struct {
	Context  GlanceContext          `json:"context"`
	Label    GlanceLabel            `json:"label"`
	Metadata map[string]interface{} `json:"metadata"`
}

type GlanceContext struct {
	CloudID        string `json:"cloudId"`
	ConversationID string `json:"conversationId,omitempty"`
	UserID         string `json:"userId,omitempty"`
}

type GlanceLabel struct {
	Value string `json:"value"`
}

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/glance/
type Glance struct {
	Key      string `json:"key"`
	Name     Name   `json:"name"`
	Icon     Icon   `json:"icon"`
	Target   string `json:"target,omitempty"`
	QueryURL string `json:"queryUrl"`
	// Default weight is 100, so we should omit this if 0.
	Weight int `json:"weight,omitempty"`
	// Conditions // not supported
	Authentication string `json:"authentication,omitempty"`
}

func (g *Glance) Type() ModuleType {
	return ModuleGlance
}

func (c *clientImpl) UpdateGlanceState(cloudID, conversationID, glanceKey, stateTxt string) error {
	url := fmt.Sprintf("%s/app/module/chat/conversation/chat:glance/%s/state", c.apiBaseURL, glanceKey)
	update := &GlanceBody{
		Label: GlanceLabel{stateTxt},
		Context: GlanceContext{
			CloudID:        cloudID,
			ConversationID: conversationID,
		},
	}
	return c.postToAPI(url, update)
}
