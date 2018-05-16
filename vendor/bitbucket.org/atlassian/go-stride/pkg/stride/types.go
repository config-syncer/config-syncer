package stride

import (
	"time"
)

// Message represents an incoming message.
type Message struct {
	CloudID      string `json:"cloudId"`
	Conversation struct {
		ID string `json:"id"`
	} `json:"conversation"`
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`
	Message struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"message"`
}

type Conversations struct {
	Values []*ConversationCommon `json:"values"`
	// API limits us to a certain number of results per request. The following attributes
	// allow us to determine whether we've hit that limit and how to fetch the next set.
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
	Size   int    `json:"size"`
}

type ConversationCommon struct {
	AvatarURL  string `json:"avatarUrl"`
	CloudID    string `json:"cloudId"`
	Created    time.Time
	ID         string `json:"id"`
	IsArchived bool   `json:"isArchived"`
	Modified   time.Time
	Name       string `json:"name"`
	Privacy    string `json:"privacy"`
	Topic      string `json:"topic"`
	Type       string `json:"type"`
}

type User struct {
	ID          string `json:"id"`
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`
	NickName    string `json:"nickName"`
	Name        struct {
		GivenName  string `json:"givenName"`
		Formatted  string `json:"formatted"`
		FamilyName string `json:"familyName"`
	} `json:"name"`
	Title    string `json:"title"`
	Active   bool   `json:"active"`
	Timezone string `json:"timezone"`
	Emails   []struct {
		Primary bool   `json:"primary"`
		Value   string `json:"value"`
	} `json:"emails"`
	Photos []struct {
		Primary bool   `json:"primary"`
		Value   string `json:"value"`
	} `json:"photos"`
	Meta struct {
		ResourceTyoe string `json:"resourceType"`
		LastModified string `json:"lastModified"`
		Created      string `json:"created"`
	} `json:"meta"`
}
