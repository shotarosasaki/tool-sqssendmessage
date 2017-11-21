# tool-sqssendmessage

# 概要

指定のペースで指定の件数分、SQSにメッセージを配信し続けます。

配信するコンテンツは「SQS投入負荷ツール：No.n-n」テキストコンテンツ１件だけです。（左記の「n」には外側・内側それぞれのループカウンタが入ります）

ローカルのElasticMQ宛てに配信するツールとして作りましたが、SQS宛でもいけるはずです。

# 実行時引数

## acid

ACID

### 引数無し時のデフォルト値

1234abcd5678efgh

## userID

ユーザID

### 引数無し時のデフォルト値

U1234567890abcdefgh1234567890abcd

## queueName

SQSキュー名

### 引数無し時のデフォルト値

local_line_messages

## loopCount

ループカウント

### 引数無し時のデフォルト値

1

## delay

１ループ毎の遅延時間（「500ms」のようなtime.Duration形式）

### 引数無し時のデフォルト値

1s

## entryCount

１ループ毎の送信メッセージ数

### 引数無し時のデフォルト値

1

## endpoint ※ローカルでElasticMQに接続する場合用

SQSエンドポイント　※「/」で終わるようにしてください。

### 引数無し時のデフォルト値

http://localhost:9324/queue/

# コマンド実行例
go run main.go \
-acid 123xxx99999xxxxx \
-userID U9999999999aaaaaaaaaa9999999999bbbb \
-queueURL http://localhost:9324/queue/local_line_messages \
-loopCount 2 \
-delay 1s \
-entryCount 2 \
-endpoint http://localhost:9324/queue/ \

# 配信ペース調整パターン

## 秒間10メッセージ（※ 1ループあたり10メッセージ配信する場合、トータル1000メッセージ）

-loopCount 100 \
-delay 1s \

## 秒間20メッセージ（※ 1ループあたり10メッセージ配信する場合、トータル2000メッセージ）

-loopCount 200 \
-delay 500ms \

## 秒間40メッセージ（※ 1ループあたり10メッセージ配信する場合、トータル4000メッセージ）

-loopCount 400 \
-delay 250ms \

## 秒間100メッセージ（※ 1ループあたり10メッセージ配信する場合、トータル10000メッセージ）

-loopCount 1000 \
-delay 100ms \

## 秒間500メッセージ（※ 1ループあたり10メッセージ配信する場合、トータル50000メッセージ）

-loopCount 5000 \
-delay 20ms \
