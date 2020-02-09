package xfyun

type Client struct {
	ApiSecret string
	ApiKey    string
	HostUrl   string
}

func NewClient(key, secret string) *Client {
	return &Client{
		ApiKey:    key,
		ApiSecret: secret,
		HostUrl:   "wss://tts-api.xfyun.cn/v2/tts",
	}
}
