package stride

// https://developer.atlassian.com/cloud/stride/apis/modules/about-stride-modules/

const (
	ModuleActionTarget  = ModuleType("chat:actionTarget")
	ModuleBot           = ModuleType("chat:bot")
	ModuleBotMessages   = ModuleType("chat:bot:messages")
	ModuleConfiguration = ModuleType("chat:configuration")
	ModuleDialog        = ModuleType("chat:dialog")
	ModuleExternalPage  = ModuleType("chat:externalPage")
	ModuleGlance        = ModuleType("chat:glance")
	ModuleInputAction   = ModuleType("chat:inputAction")
	ModuleMessageAction = ModuleType("chat:messageAction")
	ModuleSidebar       = ModuleType("chat:sidebar")
	ModuleWebhook       = ModuleType("chat:webhook")
)

type ModuleType string

type Module interface {
	Type() ModuleType
}

type Name struct {
	Value string `json:"value"`
	I18n  string `json:"i18n,omitempty"`
}

type Icon struct {
	URL string `json:"url"`
	// URL2x used for descriptor, glances, etc.
	URL2x string `json:"url@2x,omitempty"`
	// Label used in document nodes.
	Label string `json:"label,omitempty"`
}

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/bot/
type Bot struct {
	Key           string `json:"key"`
	Mention       *URL   `json:"mention,omitempty"`
	DirectMessage *URL   `json:"directMessage,omitempty"`
}

func (b *Bot) Type() ModuleType {
	return ModuleBot
}

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/bot-messages/
type BotMessages struct {
	Key     string `json:"key"`
	Pattern string `json:"pattern"`
	URL     string `json:"url"`
}

func (b *BotMessages) Type() ModuleType {
	return ModuleBotMessages
}

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/configuration/
type Configuration struct {
	Key          string             `json:"key"`
	Page         *ConfigurationPage `json:"page,omitempty"`
	ExternalPage *ConfigurationPage `json:"externalPage,omitempty"`
	State        struct {
		URL string `json:"url"`
	} `json:"state"`
	// Defaults to "jwt", alternative is "none"
	Authentication string `json:"authentication,omitempty"`
}

func (c *Configuration) Type() ModuleType {
	return ModuleConfiguration
}

type ConfigurationPage struct {
	// URL to load in the client (Page) or external browser (ExternalPage).
	URL    string `json:"url,omitempty"`
	Target string `json:"target,omitempty"`
}

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/dialog/
type Dialog struct {
	Key     string        `json:"key"`
	Title   Name          `json:"title"`
	Options DialogOptions `json:"options"`
	URL     string        `json:"url"`
	// Defaults to "jwt", alternative is "none"
	Authentication string `json:"authentication,omitempty"`
}

func (d *Dialog) Type() ModuleType {
	return ModuleDialog
}

type DialogOptions struct {
	Size             DialogSize     `json:"size"`
	PrimaryAction    DialogAction   `json:"primaryAction"`
	SecondaryActions []DialogAction `json:"secondaryActions"`
}

type DialogAction struct {
	Key  string `json:"key"`
	Name Name   `json:"name"`
}

type DialogSize struct {
	Width  string `json:"width"`
	Height string `json:"height"`
}

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/externalPage/
type ExternalPage struct {
	Name Name   `json:"name"`
	URL  string `json:"url"`
}

func (e *ExternalPage) Type() ModuleType {
	return ModuleExternalPage
}

// TODO: lookup authoritative definition of this object once stride docs updated with it.
type ActionTarget struct {
	Key               string `json:"key"`
	CallService       *URL   `json:"callService,omitempty"`
	OpenConfiguration *Key   `json:"openConfiguration,omitempty"`
	OpenDialog        *Key   `json:"openDialog,omitempty"`
	OpenExternalPage  *Key   `json:"openExternaPage,omitempty"`
	OpenSidebar       *Key   `json:"openSidebar,omitempty"`
}

func (t *ActionTarget) Type() ModuleType {
	return ModuleActionTarget
}

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/inputAction/
type InputAction struct {
	Key    string `json:"key"`
	Name   Name   `json:"name"`
	Target string `json:"target"`
	// https://developer.atlassian.com/cloud/stride/apis/modules/chat/condition/
	// Conditions // not supported
}

func (i *InputAction) Type() ModuleType {
	return ModuleInputAction
}

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/messageAction/
type MessageAction struct {
	Key    string `json:"key"`
	Name   Name   `json:"name"`
	Target string `json:"target"`
	// Default weight is 100, so we should omit this if 0.
	Weight int `json:"weight,omitempty"`
}

func (m *MessageAction) Type() ModuleType {
	return ModuleMessageAction
}

// https://developer.atlassian.com/cloud/stride/apis/modules/chat/sidebar/
type Sidebar struct {
	Key  string `json:"key"`
	Name Name   `json:"name"`
	URL  string `json:"url"`
	Icon *Icon  `json:"icon,omitempty"`
	// Defaults to "jwt", alternative is "none"
	Authentication string `json:"authentication,omitempty"`
}

func (s *Sidebar) Type() ModuleType {
	return ModuleSidebar
}
