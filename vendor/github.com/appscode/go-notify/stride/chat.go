package stride

import (
	"errors"
	"net/http"

	"bitbucket.org/atlassian/go-stride/pkg/stride"
	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
)

const (
	UID = "stride"
)

type Options struct {
	CloudID      string   `envconfig:"CLOUD_ID" required:"true"`
	RoomToken    string   `envconfig:"ROOM_TOKEN"`
	ClientID     string   `envconfig:"CLIENT_ID"`
	ClientSecret string   `envconfig:"CLIENT_SECRET"`
	To           []string `envconfig:"TO"`
}

type client struct {
	opt  Options
	body string
}

var _ notify.ByChat = &client{}

func New(opt Options) *client {
	return &client{opt: opt}
}

func Default() (*client, error) {
	var opt Options
	err := envconfig.Process(UID, &opt)
	if err != nil {
		return nil, err
	}
	return New(opt), nil
}

func Load(loader envconfig.LoaderFunc) (*client, error) {
	var opt Options
	err := envconfig.Load(UID, &opt, loader)
	if err != nil {
		return nil, err
	}
	return New(opt), nil
}

func (c client) UID() string {
	return UID
}

func (c client) WithBody(body string) notify.ByChat {
	c.body = body
	return &c
}

func (c client) To(to string, cc ...string) notify.ByChat {
	c.opt.To = append([]string{to}, cc...)
	return &c
}

func (c *client) Send() error {
	if len(c.opt.To) == 0 {
		return errors.New("missing to")
	}

	if c.opt.RoomToken == "" && (c.opt.ClientID == "" || c.opt.ClientSecret == "") {
		return errors.New("missing auth")
	}

	var sc stride.Client
	if c.opt.RoomToken != "" {
		if len(c.opt.To) > 1 {
			return errors.New(`multiple "to" with room_token is not supported`)
		}
		sc = stride.NewRoomClient(c.opt.RoomToken, http.DefaultClient)
	} else {
		sc = stride.New(c.opt.ClientID, c.opt.ClientSecret)
	}

	for _, to := range c.opt.To {
		conversation, err := sc.GetConversationByName(c.opt.CloudID, to)
		if err != nil {
			return err
		}
		if err := stride.SendText(sc, conversation.CloudID, conversation.ID, c.body); err != nil {
			return err
		}
	}
	return nil
}
