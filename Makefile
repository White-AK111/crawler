APP = crawler.app
PROJECT = github.com/t0pep0/GB_best_go1
VERSION = "1.0.0"
COMMIT = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

all: run

run: clean build
	@echo "+ $@"
	./${APP}

clean:
	@rm -f ./${APP}

build: test lint
	@echo "+ $@"
	@go build -ldflags $(LDFLAGS) -o $(APP) $(PROJECT)

test:
	@echo "+ $@"
	@go test -v -cover ./...

lint: bootstrap
	@echo "+ $@"
	@golangci-lint run ./...

HAS_LINT := $(shell command -v golangci-lint;)

bootstrap:
ifndef HAS_LINT
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
endif

.PHONY: all \
	run \
	build \
	test \
	lint \
	clean \
	bootstrap
