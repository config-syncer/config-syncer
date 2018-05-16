package stride

// Payload represents an outgoing message.
type Payload struct {
	Body *Document `json:"body"`
}

// Document is a representation of a rich object or message.
type Document struct {
	Version int           `json:"version"`
	Content []interface{} `json:"content"`
}

// Paragraph is a container for inline group nodes.
type Paragraph struct {
	Content []InlineGroupNode `json:"content"`
}

type InlineGroupNode interface {
	Type() inlineGroupNode
}

type inlineGroupNode string

type Text struct {
	Text  string `json:"text"`
	Marks []Mark `json:"marks,omitempty"`
}

type Lozenge struct {
	Text       string       `json:"text"`
	Appearance LozengeState `json:"appearance,omitempty"`
	Bold       bool         `json:"bold,omitempty"`
}

type Badge struct {
	Appearance BadgeState `json:"appearance,omitempty"`
	Max        int        `json:"max"`
	Theme      string     `json:"theme,omitempty"`
	Value      int        `json:"value"`
}

// Emoji is an inline group node.
type Emoji struct {
	Attrs EmojiAttrs `json:"attrs"`
}

type EmojiAttrs struct {
	ID        string `json:"id,omitempty"`
	ShortName string `json:"shortName"`
	Text      string `json:"text,omitempty"`
}

type MentionAttrs struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type ApplicationCard struct {
	Attrs *ApplicationCardAttrs `json:"attrs"`
}

type ApplicationCardAttrs struct {
	Background  *URL              `json:"background,omitempty"`
	Collapsible bool              `json:"collapsible"`
	Context     *CardAttrContext  `json:"context,omitempty"`
	Description *CardAttrText     `json:"description,omitempty"`
	Details     []CardAttrDetails `json:"details,omitempty"`
	Link        *URL              `json:"link,omitempty"`
	Preview     *URL              `json:"preview,omitempty"`
	Text        string            `json:"text"`
	TextURL     string            `json:"textUrl,omitempty"`
	Title       CardAttrTitle     `json:"title"`
}

type CardAttrContext struct {
	Icon Icon   `json:"icon"`
	Text string `json:"text"`
}

type CardAttrDetails struct {
	Badge   *Badge          `json:"badge,omitempty"`
	Icon    *Icon           `json:"icon,omitempty"`
	Lozenge *Lozenge        `json:"lozenge,omitempty"`
	Text    string          `json:"text,omitempty"`
	Title   string          `json:"title,omitempty"`
	Users   []*CardAttrUser `json:"users,omitempty"`
}

type CardAttrText struct {
	Text string `json:"text"`
}

type CardAttrTitle struct {
	Text string        `json:"text"`
	User *CardAttrUser `json:"user,omitempty"`
}

type CardAttrUser struct {
	ID   string `json:"id,omitempty"`
	Icon *Icon  `json:"icon,omitempty"`
}

// Mention is an inline group node.
type Mention struct {
	Attrs MentionAttrs `json:"attrs"`
}

// HardBreak is an inline group node.
type HardBreak struct {
}

// Key represents an object with a single field, Key.
type Key struct {
	Key string `json:"key"`
}

// URL repesents an object with a single field, URL.
type URL struct {
	URL string `json:"url"`
}

// ID repesents an object with a single field, URL.
type ID struct {
	ID string `json:"id"`
}

// Panel is an inline group node.
// https://developer.atlassian.com/cloud/stride/apis/document/nodes/panel/
type Panel struct {
	Content []interface{} `json:"content"`
	Attrs   PanelAttrs    `json:"attrs"`
}

type PanelAttrs struct {
	// Valid values: “info” | “note” | “tip” | “warning”
	PanelType string `json:"panelType"`
}
