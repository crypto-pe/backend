FROM golang:1.17-alpine3.15 as builder
WORKDIR /app

ARG GITBRANCH=""
ARG GITCOMMIT=""
ARG GITCOMMITDATE=""
ARG VERSION=""

RUN apk add --no-cache \
    bash \
    ca-certificates \
    make \
    git \
    build-base \
    gcc

ADD go.mod ./
ADD go.sum ./
RUN go mod download -x

COPY . ./

RUN apk update && apk upgrade && apk add --update alpine-sdk && \
    apk add --no-cache bash git openssh make cmake 

RUN apk add build-base

RUN GOGC=off \
	go build -tags='' \
	-o /app/main \
	-gcflags='-e' \
	-ldflags='-X "github.com/crypto-pe/backend.VERSION=$(VERSION)" -X "github.com/crypto-pe/backend.GITBRANCH=$(GITBRANCH)" -X "github.com/crypto-pe/backend.GITCOMMIT=$(GITCOMMIT)" -X "github.com/crypto-pe/backend.GITCOMMITDATE=$(GITCOMMITDATE)" -X "github.com/crypto-pe/backend.GITCOMMITAUTHOR=$(GITCOMMITAUTHOR)"' \
	./cmd/api-server

EXPOSE 8000

FROM alpine:3.15
WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main", "--config", "/etc/api.conf"]
