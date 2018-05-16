package stride

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/webhook/

const (
	// "conversation:updates" actions.
	ConversationUpdatesArchive   = "archive"
	ConversationUpdatesCreate    = "create"
	ConversationUpdatesDelete    = "delete"
	ConversationUpdatesUnarchive = "unarchive"
	ConversationUpdatesUpdate    = "update"
	// "roster:updates" actions.
	RosterUpdatesAdd    = "add"
	RosterUpdatesRemove = "remove"
)

type Webhook struct {
	Key   string `json:"key"`
	Event string `json:"event"`
	URL   string `json:"url"`
}

func (w *Webhook) Type() ModuleType {
	return ModuleWebhook
}

type ConversationMeta struct {
	AvatarURL  string `json:"avatarUrl"`
	Created    string `json:"created"`
	ID         string `json:"id"`
	IsArchived bool   `json:"isArchived"`
	Modified   string `json:"modified"`
	Name       string `json:"name"`
	Privacy    string `json:"privacy"`
	Topic      string `json:"topic"`
	Type       string `json:"type"`
}

type ConversationUpdate struct {
	Action       string           `json:"action"`
	CloudID      string           `json:"cloudId"`
	Conversation ConversationMeta `json:"conversation"`
	Type         string           `json:"type"`
	Initiator    ID               `json:"initiator"`
}

type RosterUpdate struct {
	Action       string           `json:"action"`
	CloudID      string           `json:"cloudId"`
	Conversation ConversationMeta `json:"conversation"`
	Type         string           `json:"type"`
	User         ID               `json:"initiator"`
}
