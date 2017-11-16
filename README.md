# tool-sqssendmessage

# 概要

指定のペースで指定の件数分、SQSにメッセージを配信し続けます。

配信するコンテンツは「Hello, world」テキストコンテンツ１件だけです。

ローカルのElasticMQ宛てに配信するツールとして作りましたが、SQS宛でもいけるはずです。

# 実行時引数

## acid

ACID

## userID

ユーザID

## endpoint

SQSエンドポイント　※「/」で終わるようにしてください。

## queueName

SQSキュー名

## loopCount

ループカウント

## delay

１ループ毎の遅延時間（「500ms」のようなtime.Duration形式）

## entryCount

１ループ毎の送信メッセージ数

# コマンド実行例
go run main.go \
-acid 123xxx99999xxxxx \
-userID U9999999999aaaaaaaaaa9999999999bbbb \
-endpoint http://localhost:9324/queue/ \
-queueURL http://localhost:9324/queue/local_line_messages \
-loopCount 100 \
-delay 1s \
-entryCount 10

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
