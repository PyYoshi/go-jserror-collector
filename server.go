package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/t-k/fluent-logger-golang/fluent"
)

// 1px x 1px gif
var ResIGif = []byte("GIF87a\x01\x00\x01\x00\x80\x00\x00\xff\xff\xff\xff\xff\xff,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02D\x01\x00;")

// FluentdConfig fluentd接続設定 詳細は https://github.com/t-k/fluent-logger-golang 参照
var FluentdConfig = fluent.Config{
	FluentHost:  "127.0.0.1",
	FluentPort:  24224,
	Timeout:     3 * time.Second,
	BufferLimit: 8 * 1024 * 1024,
	RetryWait:   500,
	MaxRetry:    13,
	TagPrefix:   "",
}

// JSErrorMessage 送られてきたJSONを格納するオブジェクト
type JSErrorMessage struct {
	Message       string `json:"m"`
	JavascriptURL string `json:"u"`
	LineNumber    int64  `json:"l"`
	ColumnNumber  int64  `json:"c"`
	RemoteAddr    string `json:"r"`
	UserAgent     string `json:"a"`
	Timestamp     string `json:"t"`
}

// JsecHandler 送られてきたJSONをfluentdへ飛ばすHTTPハンドラ
type JsecHandler struct {
	FLogger   *fluent.Fluent
	FluentTag string // NOTICE: Fluentdへ送信するときのタグ
}

// SendLog Fluentdへログを送信
func (jh *JsecHandler) SendLog(jsErrorMessage *JSErrorMessage) {
	jh.FLogger.Post(jh.FluentTag, &jsErrorMessage)
}

func (jh *JsecHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlObj, err := url.Parse(r.URL.String())
	if err == nil {
		query := urlObj.Query()
		if _, ok := query["r"]; ok && len(query["r"]) > 0 {
			var jsErrorMessage JSErrorMessage
			err = json.Unmarshal([]byte(query["r"][0]), &jsErrorMessage)
			if err == nil {
				utcTime := time.Now().UTC()
				timestamp := utcTime.Format("2006/01/02 15:04:05 MST")
				jsErrorMessage.Timestamp = timestamp
				jsErrorMessage.RemoteAddr = r.RemoteAddr
				jsErrorMessage.UserAgent = r.UserAgent()
				jh.SendLog(&jsErrorMessage)
			}
		}
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "image/gif")
	w.Write(ResIGif)
}

func main() {
	flogger, flerr := fluent.New(FluentdConfig)
	if flerr != nil {
		log.Fatal(flerr)
	}
	defer flogger.Close()

	mux := http.NewServeMux()
	jh := &JsecHandler{FLogger: flogger, FluentTag: "debug.jsec"}
	mux.Handle("/jsec", jh)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
