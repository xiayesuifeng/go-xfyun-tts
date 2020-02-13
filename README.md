# go-xfyun-tts

xfyun online tts binding for Go (Golang)

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/xiayesuifeng/go-xfyun-tts)](https://goreportcard.com/report/gitlab.com/xiayesuifeng/go-xfyun-tts)
[![GoDoc](https://godoc.org/gitlab.com/xiayesuifeng/go-xfyun-tts?status.svg)](https://godoc.org/gitlab.com/firerainos/xiayesuifeng/go-xfyun-tts)
[![Sourcegraph](https://sourcegraph.com/gitlab.com/xiayesuifeng/go-xfyun-tts/-/badge.svg)](https://sourcegraph.com/gitlab.com/xiayesuifeng/go-xfyun-tts)

### Installation

    go get -u gitlab.com/xiayesuifeng/go-xfyun-tts
    
### Example Code
```go
package main

import (
    "gitlab.com/xiayesuifeng/go-xfyun-tts"
    "io/ioutil"
    "log"
)

func main() { 
    client := NewClient(APP_ID, API_key, API_SECRET)
    b := NewBusiness("xiaoyan")
    b.Bgs = 1
    d,err := client.GetAudio(b,"这是一段合成的语音")
    if err != nil {
    	log.Println(err)
    }
    
    ioutil.WriteFile("audio.pac",d.Bytes(),0664)
}
```

## License
This package is released under [GPLv3](LICENSE).