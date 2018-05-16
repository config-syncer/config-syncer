package stride

type LifeCycle struct {
	Installed   string `json:"installed"`
	Uninstalled string `json:"uninstalled"`
}

// LifeCyclePayload is received by the app on a LifeCycle event.
// https://developer.atlassian.com/cloud/stride/blocks/app-lifecycle/

type LifeCyclePayload struct {
	Key           string `json:"key"`
	ProductType   string `json:"productType"`
	CloudID       string `json:"cloudId"`
	ResourceType  string `json:"resourceType"`
	ResourceID    string `json:"resourceId"`
	EventType     string `json:"eventType"`
	UserID        string `json:"userId"`
	OAuthClientID string `json:"oauthClientId"`
	Version       string `json:"version"`
}
