package bootstrap

import social "github.com/kkdai/line-login-sdk-go"

type LineEnv struct {
	ChannelID     string `env:"CHANNEL_ID"`
	ChannelSecret string `env:"CHANNEL_SECRET"`
}

func NewLineSocialClient(env *Env) *social.Client {
	client, _ := social.New(env.Line.ChannelID, env.Line.ChannelSecret)
	return client
}
