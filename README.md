# tool-sqssendmessage

# 概要

指定のペースで指定の件数分、SQSにメッセージを配信します。

配信するコンテンツは「SQS投入負荷ツール：No.l-m-n」テキストコンテンツ１件だけです。

※上記の「l」には「0」から引数「duration」-1 までの値が入ります。
※上記の「m」には「0」からループカウント（引数「throughput」/10）までの値が入ります。
※上記の「n」には「0」から「9」までの値が入ります。

ローカルのElasticMQ宛てに配信するツールとして作りました。（一応、develop環境のSQS宛ても確認しました。）

# 実行時引数

## acid

ACID

#### [引数無し時のデフォルト値]

1234abcd5678efgh

## userID

ユーザID

#### [引数無し時のデフォルト値]

U1234567890abcdefgh1234567890abcd

## queueName

SQSキュー名

#### [引数無し時のデフォルト値]

local_line_messages

## throughput

１秒当たりのSQS投入メッセージ数　※必ず 10 の倍数を指定してください。

#### [引数無し時のデフォルト値]

10

## duration

継続時間（単位：秒）

#### [引数無し時のデフォルト値]

5

## endpoint ※ローカルでElasticMQに接続する場合用

SQSエンドポイント　※「/」で終わるようにしてください。

#### [引数無し時のデフォルト値]

http://localhost:9324/queue/

# コマンド実行例(秒間30メッセージ配信を10秒間継続)
go run main.go \
-acid 736afbdc388752eb \
-userID U9999999999aaaaaaaaaa9999999999bbbb \
-endpoint http://localhost:9324/queue/ \
-queueName local_line_messages \
-throughput 30 \
-duration 10
