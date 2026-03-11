.DEFAULT_GOAL := vet
.PHONY: format generator vet build
.EXPORT_ALL_VARIABLES:

LINUX_AMD64 := GOOS=linux GOARCH=amd64
GIT_BRANCH := $(shell git rev-parse  --abbrev-ref HEAD)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
COMMIT_TXT := ${GIT_BRANCH}/${GIT_COMMIT}
BUILD_DATE := $(shell date)
BUILD_ENV := $(shell uname -a)

LOG_LEVEL := DEBUG

format: 
	@go mod tidy -e
	@go vet ./...
	@gofmt -s -w .
	@go fix ./...

unit-tests: format
	@go test github.com/pdaccess/ws/internal/...

cicd-tests: format
	@LOG_LEVEL=${LOG_LEVEL} ginkgo -v cicd/tests
	
build: format tidy cmd/main.go
	docker build --no-cache --build-arg COMMIT_TXT="${COMMIT_TXT}" --build-arg BUILD_DATE="${BUILD_DATE}" --build-arg BUILD_ENV="${BUILD_ENV}" -t registry.h2hsecure.com/pda/core/ws:${GIT_COMMIT} -f Dockerfile .
	
ci-build: cmd/main.go
	docker build --no-cache --build-arg COMMIT_TXT="${COMMIT_TXT}" --build-arg BUILD_DATE="${BUILD_DATE}" --build-arg BUILD_ENV="${BUILD_ENV}" -t registry.h2hsecure.com/pda/core/ws:${GIT_COMMIT} -f Dockerfile .

ci-push:
	docker push registry.h2hsecure.com/pda/core/ws:${GIT_COMMIT}

internal-generator:
	@rm -r internal/handlers/external/server.gen.go || true
	@rm -r pkg/http/client.gen.go || true

	@oapi-codegen -config resources/api-config-client.yaml resources/corews-api.yaml
	@oapi-codegen -config resources/api-config-server.yaml resources/corews-api.yaml
	@cp resources/corews-api.yaml internal/handlers/custom/openapi.yaml

generate: internal-generator format

clean:
	@rm -f corews*

local:
	@go run cmd/main.go --debug --console
