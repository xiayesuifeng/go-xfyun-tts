package xfyun

type Client struct {
	ApiSecret string
	ApiKey    string
	Host      string
}

func NewClient(key, secret, host string) *Client {
	return &Client{
		ApiKey:    key,
		ApiSecret: secret,
		Host:      host,
	}
}
