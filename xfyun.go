package xfyun

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/url"
	"time"
)

type Client struct {
	AppID     string
	ApiSecret string
	ApiKey    string
	HostUrl   string
}

type Common struct {
	AppID string `json:"app_id"`
}

type Business struct {
	// 引擎类型，可选值：aisound（普通效果）, intp65（中文）, intp65_en（英文）, mtts（小语种，需配合小语种发音人使用）, xtts（优化效果）, 默认为intp65
	Ent string `json:"ent"`

	// 音频编码，可选值：
	// raw：未压缩的pcm
	// speex-org-wb;7： 标准开源speex（for speex_wideband，即16k）数字代表指定压缩等级（默认等级为8）
	// speex-org-nb;7： 标准开源speex（for speex_narrowband，即8k）数字代表指定压缩等级（默认等级为8）
	// speex;7：压缩格式，压缩等级1~10，默认为7（8k讯飞定制speex）
	// speex-wb;7：压缩格式，压缩等级1~10，默认为7（16k讯飞定制speex）
	Aue string `json:"aue"`

	// 音频采样率，可选值：
	// audio/L16;rate=8000：合成8K 的音频
	// audio/L16;rate=16000：合成16K 的音频
	// auf不传值：合成16K 的音频
	Auf string `json:"auf"`

	// 发音人
	Vcn string `json:"vcn"`

	// 语速，可选值：[0-100]，默认为50
	Speed int `json:"speed"`

	// 音量，可选值：[0-100]，默认为50
	Volume int `json:"volume"`

	// 音高，可选值：[0-100]，默认为50
	Pitch int `json:"pitch"`

	// 合成音频的背景音
	// 0:无背景音（默认值）
	// 1:有背景音
	Bgs int `json:"bgs"`

	// 文本编码格式
	// GB2312
	// GBK
	// BIG5
	// UNICODE(小语种必须使用UNICODE编码)
	// GB18030
	// UTF8
	Tte string `json:"tte"`

	// 设置英文发音方式：
	// 0：自动判断处理，如果不确定将按照英文词语拼写处理（缺省）
	// 1：所有英文按字母发音
	// 2：自动判断处理，如果不确定将按照字母朗读
	// 默认按英文单词发音
	Reg string `json:"reg"`

	// 是否读出标点：
	// 0：不读出所有的标点符号（默认值）
	// 1：读出所有的标点符号
	Ram string `json:"ram"`

	// 合成音频数字发音方式
	// 0：自动判断（默认值）
	// 1：完全数值
	// 2：完全字符串
	// 3：字符串优先
	Rdn string `json:"rdn"`
}

type Data struct {
	Text   string `json:"text"`
	Status int    `json:"status"`
}

type TTSRequest struct {
	Common Common `json:"common"`

	Business Business `json:"business"`

	Data Data `json:"data"`
}

type TTSReturnData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Data    struct {
		Audio  string `json:"audio"`
		Ced    string `json:"ced"`
		Status int    `json:"status"`
	}
}

func NewClient(appID, key, secret string) *Client {
	return &Client{
		AppID:     appID,
		ApiKey:    key,
		ApiSecret: secret,
		HostUrl:   "wss://tts-api.xfyun.cn/v2/tts",
	}
}

func NewBusiness(Vcn string) Business {
	return Business{
		Ent:    "intp65",
		Aue:    "raw",
		Auf:    "audio/L16;rate=16000",
		Vcn:    Vcn,
		Speed:  50,
		Volume: 50,
		Pitch:  50,
		Tte:    "UTF8",
		Reg:    "0",
		Ram:    "0",
		Rdn:    "0",
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

func (client *Client) GetAudio(business Business, text string) (audio bytes.Buffer, err error) {
	ws, resp, err := websocket.DefaultDialer.Dial(client.getWebsocketUrl(), nil)
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		return audio, errors.New(string(b))
	}
	defer ws.Close()

	req := TTSRequest{
		Common: Common{
			AppID: client.AppID,
		},

		Business: business,

		Data: Data{Text: base64.StdEncoding.EncodeToString([]byte(text)), Status: 2},
	}

	done := make(chan error)
	defer close(done)

	go func() {
		data := TTSReturnData{}

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				done <- err
				break
			}

			if err := json.Unmarshal(msg, &data); err != nil {
				done <- err
				break
			}

			if data.Code != 0 {
				done <- errors.New(data.Message)
				break
			} else {
				bytes, err := base64.StdEncoding.DecodeString(data.Data.Audio)
				if err != nil {
					done <- err
					break
				}

				audio.Write(bytes)
				if data.Data.Status == 2 {
					done <- nil
					break
				}
			}
		}
	}()

	err = ws.WriteJSON(req)
	if err != nil {
		done <- err
	}

	err = <-done
	return audio, err
}
