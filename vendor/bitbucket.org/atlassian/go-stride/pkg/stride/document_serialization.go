package stride

import (
	"bytes"
	"encoding/json"
	"io"
)

func (n *Emoji) Type() inlineGroupNode {
	return "emoji"
}

func (d *Document) MarshalJSON() ([]byte, error) {
	type Alias Document

	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "doc",
		Alias: (*Alias)(d),
	}
	return json.Marshal(aux)
}

func (p *Paragraph) MarshalJSON() ([]byte, error) {
	type Alias Paragraph

	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "paragraph",
		Alias: (*Alias)(p),
	}
	return json.Marshal(aux)
}

// Text is an inline group node.
func (n *Text) Type() inlineGroupNode {
	return "text"
}

func (n *Text) MarshalJSON() ([]byte, error) {
	type Alias Text

	aux := &struct {
		Type inlineGroupNode `json:"type"`
		*Alias
	}{
		Type:  n.Type(),
		Alias: (*Alias)(n),
	}
	return json.Marshal(aux)
}

func (n *Emoji) MarshalJSON() ([]byte, error) {
	type Alias Emoji

	aux := &struct {
		Type inlineGroupNode `json:"type"`
		*Alias
	}{
		Type:  n.Type(),
		Alias: (*Alias)(n),
	}
	return json.Marshal(aux)
}

func (n *HardBreak) Type() inlineGroupNode {
	return "hardBreak"
}

func (n *HardBreak) MarshalJSON() ([]byte, error) {
	type Alias HardBreak

	aux := &struct {
		Type inlineGroupNode `json:"type"`
		*Alias
	}{
		Type:  n.Type(),
		Alias: (*Alias)(n),
	}
	return json.Marshal(aux)
}

func (n *Mention) Type() inlineGroupNode {
	return "mention"
}

func (n *Mention) MarshalJSON() ([]byte, error) {
	type Alias Mention

	aux := &struct {
		Type inlineGroupNode `json:"type"`
		*Alias
	}{
		Type:  n.Type(),
		Alias: (*Alias)(n),
	}
	return json.Marshal(aux)
}

func (n *ApplicationCard) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCard

	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "applicationCard",
		Alias: (*Alias)(n),
	}
	return json.Marshal(aux)
}

// Panel is an inline group node.
func (n *Panel) Type() inlineGroupNode {
	return "panel"
}

func (n *Panel) MarshalJSON() ([]byte, error) {
	type Alias Panel

	aux := &struct {
		Type inlineGroupNode `json:"type"`
		*Alias
	}{
		Type:  n.Type(),
		Alias: (*Alias)(n),
	}
	return json.Marshal(aux)
}

func (payload *Payload) String() string {
	bs, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func (p *Payload) AsText() string {
	buf := &bytes.Buffer{}
	AsText(p.Body, buf)
	return buf.String()
}

func AsText(elem interface{}, buf io.Writer) {
	switch x := elem.(type) {
	case *Document:
		for _, c := range x.Content {
			AsText(c, buf)
		}
	case *ApplicationCard:
		io.WriteString(buf, x.Attrs.Title.Text)
		io.WriteString(buf, "\n")
		io.WriteString(buf, x.Attrs.Text)
		io.WriteString(buf, "\n")
		io.WriteString(buf, x.Attrs.Description.Text)
		io.WriteString(buf, "\n")
	case *Paragraph:
		b := false
		for _, c := range x.Content {
			if b {
				io.WriteString(buf, " ")
			}
			AsText(c, buf)
			b = true
		}
	case *Text:
		io.WriteString(buf, x.Text)
	case *HardBreak:
		io.WriteString(buf, "\n")
	default:
	}

}
