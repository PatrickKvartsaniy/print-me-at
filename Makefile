export GO111MODULE=on
export GOSUMDB=off
export GOPROXY=direct

DOCKER_COMPOSE = docker-compose -f docker-compose.yml

.PHONY: all
all: deps gen build lint test dockerise

.PHONY: build
build:
	CGO_ENABLED=0 go build -mod=vendor -a -o artifacts/svc .

.PHONY: gen
gen:
	go get github.com/99designs/gqlgen
	go run -mod=mod github.com/99designs/gqlgen generate

.PHONY: deps
deps:
	go mod tidy
	go mod download
	go mod vendor

.PHONY: test
test:
	go test -mod=vendor -count=1 -cover -v  `go list ./...`

.PHONY: lint
lint:
	golangci-lint run

.PHONY: dockerise
dockerise:
	docker build -t "print-me-at:latest" .

.PHONY: service-up
service-up:
	$(DOCKER_COMPOSE) rm --force --stop -v
	$(DOCKER_COMPOSE) up -d --force-recreate --remove-orphans --build

.PHONY: service-down
service-down:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) rm --force --stop -v
