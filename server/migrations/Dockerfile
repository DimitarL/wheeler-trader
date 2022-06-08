FROM golang:1.18.3-alpine3.16

WORKDIR /migrations

COPY cmd ./cmd
COPY sqls ./sqls
COPY go.mod go.sum ./

RUN apk add gcc musl-dev
RUN set -ex
RUN go build -o ./main ./cmd/main.go

ENTRYPOINT [ "./main" ]
