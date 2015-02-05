Javascriptエラーをバックエンドで収集するサーバ
====================================

# 仕組み

Javascript側のwindow.onerrorですべてのエラーを収集し img.srcとしてバックエンドへ送信.

バックエンドはそのimg.srcでアクセスしてきたURLをパースしてfluentd(tcp/forward)へ渡す.

とてもシンプルなものです.

# 必須

- go1.4
- fluentd

# 準備

```bash
go get github.com/t-k/fluent-logger-golang/fluent
git clone https://github.com/PyYoshi/go-jserror-collector.git
```

# 実行

実行前に

Fluentdの接続設定 "FluentdConfig"

main関数内のJsecHandler.FluentTagとHTTPサーバポート

の確認を行ってください.


1. サーバの実行

```bash
cd go-jserror-collector
go run server.go
```

2. examples/index.html を開く

3. fluentdへログが流れていることを確認

```bash
tail -f /var/log/td-agent/td-agent.log
```

# 注意

このリポジトリはあくまで例です.

利用は自己責任でよろしくお願いします.

エラーが拾えていない場合は "Same-Origin Policy" が正しく設定されていない可能性があります.

サーバ側のヘッダー "Access-Control-Allow-Origin" 及び クライアント側の "script.crossorigin" を確認してください.

RequireJSを利用している場合はrequire.createNodeをオーバライドしてcrossorigin属性を定義するようにしてください.
