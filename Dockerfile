FROM golang:1.17-alpine3.15 as builder
WORKDIR /app

ADD go.mod ./
ADD go.sum ./
RUN go mod download -x

COPY . ./

RUN apk add build-base

RUN GOGC=off go build ./cmd/cryptope/main.go
EXPOSE 8000

FROM alpine:3.15
WORKDIR /app

COPY --from=builder /app/main . 
COPY --from=builder /app/configs/api.conf .

CMD ["./main", "-config", "api.conf"]