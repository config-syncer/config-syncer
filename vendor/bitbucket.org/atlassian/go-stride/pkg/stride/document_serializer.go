package stride

import (
	"encoding/json"
	"strings"

	log "github.com/rs/zerolog/log"
)

// Container objects with generic contents need explicit unmarshal handlers :(
//
// For the gist of how this technique works
//   see https://golang.org/pkg/encoding/json/#RawMessage
//   and http://gregtrowbridge.com/golang-json-serialization-with-interfaces/

func (d *Document) UnmarshalJSON(b []byte) (err error) {
	firstPass := struct {
		Version int               `json:"version"`
		Content []json.RawMessage `json:"content"`
	}{}
	err = json.Unmarshal(b, &firstPass)
	if err != nil {
		return
	}
	d.Version = firstPass.Version

	if len(firstPass.Content) == 0 {
		// make([]..., 0) != uninitialized slice
		return
	}
	d.Content = make([]interface{}, len(firstPass.Content))

	for i, c := range firstPass.Content {
		secondPass := struct {
			Type string `json:"type"`
		}{}

		// break loop on error?
		err = json.Unmarshal(c, &secondPass)
		if err != nil {
			log.Error().Msg(err.Error())
			continue
		}

		switch strings.ToLower(secondPass.Type) {
		case "paragraph":
			x := &Paragraph{}
			err = json.Unmarshal(c, &x)
			d.Content[i] = x
		case "applicationcard":
			x := &ApplicationCard{}
			err = json.Unmarshal(c, &x)
			d.Content[i] = x
		default:
			log.Error().Msg(err.Error())
			continue
		}
	}
	return
}

func (p *Paragraph) UnmarshalJSON(b []byte) (err error) {
	firstPass := struct {
		Content []json.RawMessage `json:"content"`
	}{}
	err = json.Unmarshal(b, &firstPass)
	if err != nil {
		return
	}
	if len(firstPass.Content) == 0 {
		// make([]..., 0) != uninitialized slice
		return
	}
	p.Content = make([]InlineGroupNode, len(firstPass.Content))

	for i, c := range firstPass.Content {
		secondPass := struct {
			Type string `json:"type"`
		}{}

		// break loop on error?
		err = json.Unmarshal(c, &secondPass)
		if err != nil {
			log.Error().Msg(err.Error())
			continue
		}

		switch strings.ToLower(secondPass.Type) {
		case "text":
			x := &Text{}
			err = json.Unmarshal(c, &x)
			p.Content[i] = x
		default:
			log.Error().Msg(err.Error())
			continue
		}
	}
	return
}

func (t *Text) UnmarshalJSON(b []byte) (err error) {
	firstPass := struct {
		Text  string            `json:"text"`
		Marks []json.RawMessage `json:"marks,omitempty"`
	}{}
	err = json.Unmarshal(b, &firstPass)
	if err != nil {
		return
	}
	t.Text = firstPass.Text
	if len(firstPass.Marks) == 0 {
		// make([]..., 0) != uninitialized slice
		return
	}
	t.Marks = make([]Mark, len(firstPass.Marks))

	for i, c := range firstPass.Marks {
		secondPass := struct {
			Type  string          `json:"type"`
			Attrs json.RawMessage `json:"attrs"`
		}{}

		// break loop on error?
		err = json.Unmarshal(c, &secondPass)
		if err != nil {
			log.Error().Msg(err.Error())
			continue
		}

		switch strings.ToLower(secondPass.Type) {
		case "link":
			x := &Link{}
			err = json.Unmarshal(secondPass.Attrs, &x)
			t.Marks[i] = x
		default:
			log.Error().Msg(err.Error())
			continue
		}
	}
	return
}

func (d *Panel) UnmarshalJSON(b []byte) error {
	return nil
}
