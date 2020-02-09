package xfyun

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

type Client struct {
	AppID     string
	ApiSecret string
	ApiKey    string
	HostUrl   string
}

func NewClient(appID, key, secret string) *Client {
	return &Client{
		AppID:     appID,
		ApiKey:    key,
		ApiSecret: secret,
		HostUrl:   "wss://tts-api.xfyun.cn/v2/tts",
	}
}

func (client *Client) getWebsocketUrl() string {
	u, _ := url.Parse(client.HostUrl)

	date := time.Now().UTC().Format(time.RFC1123)

	h := hmac.New(sha256.New, []byte(client.ApiSecret))
	h.Write([]byte("host: " + u.Host))
	h.Write([]byte("\n"))
	h.Write([]byte("date: " + date))
	h.Write([]byte("\n"))
	h.Write([]byte("GET " + u.Path + " HTTP/1.1"))

	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))

	authUrl := fmt.Sprintf("api_key=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", client.ApiKey, "hmac-sha256", "host date request-line", sha)
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))
	v := url.Values{}
	v.Add("host", u.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	return client.HostUrl + "?" + v.Encode()
}
