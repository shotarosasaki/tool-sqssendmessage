FROM golang:1.7.4

ADD . /go/src/xxxx/sqssendmessage

WORKDIR /go/src/xxxx/sqssendmessage
RUN go install xxxx/sqssendmessage

ARG ENVIRONMENT="local"
COPY config/${ENVIRONMENT}.toml /etc/sqssendmessage/sqssendmessage.toml
ENTRYPOINT /go/bin/sqssendmessage -f /etc/sqssendmessage/sqssendmessage.toml
