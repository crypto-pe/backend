GITTAG           ?= $(shell git describe --exact-match --tags HEAD 2>/dev/null || :)
GITBRANCH        ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || :)
LONGVERSION      ?= $(shell git describe --tags --long --abbrev=8 --always HEAD)$(echo -$GITBRANCH | tr / - | grep -v '\-master' || :)
VERSION          ?= $(if $(GITTAG),$(GITTAG),$(LONGVERSION))
GITCOMMIT        ?= $(shell git log -1 --date=iso --pretty=format:%H)
GITCOMMITDATE    ?= $(shell git log -1 --date=iso --pretty=format:%cd)
GITCOMMITAUTHOR  ?= $(shell git log -1 --date=iso --pretty="format:%an")
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))


.PHONY: postgres createdb dropdb migrate_up migrate_down migrate_create sqlc_generate run proto

postgres:
	docker run --name skyweaverpostgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb:
	docker exec -it skyweaverpostgres createdb --username=postgres --owner=postgres cryptope

dropdb:
	docker exec -it skyweaverpostgres dropdb -U postgres cryptope

migrate_up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -source file://data/migrations -database postgres://postgres:postgres@localhost:5432/cryptope?sslmode=disable up

migrate_down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate  -source file://data/migrations -database postgres://postgres:postgres@localhost:5432/cryptope?sslmode=disable down 1

migrate_create:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir data/migrations -seq $(FILE_NAME)

sqlc_generate:
	docker run --rm -v $(ROOT_DIR):/src -w /src kjconroy/sqlc generate

run:
	@go run github.com/VojtechVitek/rerun/cmd/rerun -watch ./ -run \
		sh -c 'make build-api && ./bin/api-server --config=./etc/example.conf'

make run_bin:
	./bin/api-server --config=./etc/example.conf

proto:
	go generate ./proto/...

build-api:
	GOGC=off \
	go build -tags='$(BUILDTAGS)' \
	-o $(ROOT_DIR)/bin/api-server \
	-gcflags='-e' \
	-ldflags='-X "github.com/crypto-pe/backend.VERSION=$(VERSION)" -X "github.com/crypto-pe/backend.GITBRANCH=$(GITBRANCH)" -X "github.com/crypto-pe/backend.GITCOMMIT=$(GITCOMMIT)" -X "github.com/crypto-pe/backend.GITCOMMITDATE=$(GITCOMMITDATE)" -X "github.com/crypto-pe/backend.GITCOMMITAUTHOR=$(GITCOMMITAUTHOR)"' \
	./cmd/api-server

build-docker:
	sudo docker buildx build --platform=linux/amd64 -f Dockerfile -t asia-south1-docker.pkg.dev/crypto-pe-351511/cryptope/api:${VERSION} --build-arg VERSION=${VERSION} --build-arg GITCOMMIT=${GITCOMMIT} --build-arg GITBRANCH=${GITBRANCH} .
	sudo docker push asia-south1-docker.pkg.dev/crypto-pe-351511/cryptope/api:${VERSION}
