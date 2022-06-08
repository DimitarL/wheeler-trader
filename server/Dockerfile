FROM golang:1.18.3-alpine3.16

WORKDIR /app

COPY cmd ./cmd
COPY common ./common
COPY migrations ./migrations
COPY pkg ./pkg
COPY go.mod go.sum ./

RUN apk add gcc musl-dev
RUN set -ex
RUN go build -o ./main ./cmd/main.go

ENV APP_HOST=""
ENV APP_PORT=8080

EXPOSE 8080

ENTRYPOINT [ "./main" ]
