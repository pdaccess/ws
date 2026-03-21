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
	@go test github.com/pdaccess/ws/internal...

cicd-tests: format
	@LOG_LEVEL=${LOG_LEVEL} CGO_LDFLAGS="-L/usr/lib" ginkgo -tags ORT -v cicd/tests
	
build: format cmd/main.go
	docker build --no-cache --build-arg COMMIT_TXT="${COMMIT_TXT}" --build-arg BUILD_DATE="${BUILD_DATE}" --build-arg BUILD_ENV="${BUILD_ENV}" -t ghcr.io/pdaccess/ws:${GIT_COMMIT} -f Dockerfile .
	docker tag ghcr.io/pdaccess/ws:${GIT_COMMIT} ghcr.io/pdaccess/ws:latest
	
ci-build: cmd/main.go
	docker build --no-cache --build-arg COMMIT_TXT="${COMMIT_TXT}" --build-arg BUILD_DATE="${BUILD_DATE}" --build-arg BUILD_ENV="${BUILD_ENV}" -t ghcr.io/pdaccess/ws:${GIT_COMMIT} -f Dockerfile .

ci-push:
	docker push ghcr.io/pdaccess/ws:${GIT_COMMIT}
	

internal-generator:
	@rm -r internal/platform/handlers/external/server.gen.go || true
	@rm -r pkg/http/client.gen.go || true

	@go tool oapi-codegen -config resources/api-config-client.yaml resources/corews-api.yaml
	@go tool oapi-codegen -config resources/api-config-server.yaml resources/corews-api.yaml
	@cp resources/corews-api.yaml internal/platform/handlers/custom/openapi.yaml

generate: internal-generator format

clean:
	@rm -f corews*

local:
	@CGO_LDFLAGS="-L/usr/lib" go build -tags ORT -o corews cmd/main.go
	@./corews --debug --console

install-onnx:
	wget -q https://github.com/daulet/tokenizers/releases/download/v1.26.0/libtokenizers.linux-amd64.tar.gz \
    && tar -xzf libtokenizers.linux-amd64.tar.gz \
    && sudo mv libtokenizers.a /usr/lib/ \
    && rm libtokenizers.linux-amd64.tar.gz

# Download ONNX Runtime for linking
	wget -q https://github.com/microsoft/onnxruntime/releases/download/v1.24.4/onnxruntime-linux-x64-1.24.4.tgz \
    && tar -xzf onnxruntime-linux-x64-1.24.4.tgz \
    && sudo cp onnxruntime-linux-x64-1.24.4/lib/libonnxruntime.so /usr/lib/ \
    && sudo ldconfig \
	&& rm onnxruntime-linux-x64-1.24.4.tgz