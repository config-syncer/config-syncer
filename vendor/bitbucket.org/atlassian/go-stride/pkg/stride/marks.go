package stride

import "encoding/json"

type Mark interface {
	Type() mark
}

type mark string

// A basicmark is a mark with a single attribute (Type).
type basicmark string

// Convenience for marshaling marks whose only attribute is type.
func (m basicmark) MarshalJSON() ([]byte, error) {
	aux := &struct {
		Type string `json:"type"`
	}{
		Type: string(m),
	}
	return json.Marshal(aux)
}

func (m basicmark) Type() mark {
	return mark(m)
}

const (
	// Marks which only have a Type attribute.
	Code      = basicmark("code")
	Em        = basicmark("em")
	Strike    = basicmark("strike")
	Strong    = basicmark("strong")
	Underline = basicmark("underline")
)

// Link is a mark.
type Link struct {
	Href  string `json:"href"`
	Title string `json:"title,omitempty"`
}

func (m *Link) Type() mark {
	return "link"
}

func (m *Link) MarshalJSON() ([]byte, error) {
	aux := &struct {
		Type  mark `json:"type"`
		Attrs struct {
			Href  string `json:"href"`
			Title string `json:"title,omitempty"`
		} `json:"attrs"`
	}{
		Type: m.Type(),
		Attrs: struct {
			Href  string `json:"href"`
			Title string `json:"title,omitempty"`
		}{
			Href:  m.Href,
			Title: m.Title,
		},
	}
	return json.Marshal(aux)
}

// Subscript is a mark.
// Subscript and Superscript share the 'subsup' mark type and are mutually exclusive.
type Subscript struct{}

func (m *Subscript) Type() mark {
	return "subsup"
}

func (m *Subscript) MarshalJSON() ([]byte, error) {
	type Alias Subscript
	type attrs struct {
		Type string `json:"type"`
	}

	aux := &struct {
		Type  mark  `json:"type"`
		Attrs attrs `json:"attrs"`
	}{
		Type:  m.Type(),
		Attrs: attrs{"sub"},
	}
	return json.Marshal(aux)
}

// Superscript is a mark.
// Subscript and Superscript share the 'subsup' mark type and are mutually exclusive.
type Superscript struct{}

func (m *Superscript) Type() mark {
	return "subsup"
}

func (m *Superscript) MarshalJSON() ([]byte, error) {
	type attrs struct {
		Type string `json:"type"`
	}

	aux := &struct {
		Type  mark  `json:"type"`
		Attrs attrs `json:"attrs"`
	}{
		Type:  m.Type(),
		Attrs: attrs{"sup"},
	}
	return json.Marshal(aux)
}

// TextColor is a mark.
type TextColor struct {
	Color string
}

func (m *TextColor) Type() mark {
	return "textColor"
}

func (m *TextColor) MarshalJSON() ([]byte, error) {
	type attrs struct {
		Color string `json:"color"`
	}

	aux := &struct {
		Type  mark  `json:"type"`
		Attrs attrs `json:"attrs"`
	}{
		Type:  m.Type(),
		Attrs: attrs{m.Color},
	}
	return json.Marshal(aux)
}
