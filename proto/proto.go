// Server
//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=api.ridl -target=go -pkg=proto -server -client -out=./api.gen.go

// Clients
//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=api.ridl -target=ts -pkg=api -client -out=./clients/api.gen.ts
//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=api.ridl -target=js -pkg=api -client -out=./clients/api.gen.js

package proto
