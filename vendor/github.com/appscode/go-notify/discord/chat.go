package discord

import (
	"errors"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/bwmarrin/discordgo"
)

const UID = "discord"

type Options struct {
	AuthToken string   `envconfig:"AUTH_TOKEN" required:"true"`
	Channel   []string `envconfig:"CHANNEL"`
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
	c.opt.Channel = append([]string{to}, cc...)
	return &c
}

func (c *client) Send() error {
	if len(c.opt.Channel) == 0 {
		return errors.New("Missing to")
	}

	discord, err := discordgo.New("Bot " + c.opt.AuthToken)
	if err != nil {
		return err
	}

	if err = discord.Open(); err != nil {
		return err
	}
	defer discord.Close()

	for _, channel := range c.opt.Channel {
		if _, err := discord.ChannelMessageSend(channel, c.body); err != nil {
			return err
		}
	}
	return nil
}
