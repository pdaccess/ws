FROM golang:1.26-alpine AS build

ARG COMMIT_TXT
ARG BUILD_DATE
ARG BUILD_ENV

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Download and install the static tokenizer library
RUN wget -q https://github.com/daulet/tokenizers/releases/download/v1.26.0/libtokenizers.linux-amd64.tar.gz \
    && tar -xzf libtokenizers.linux-amd64.tar.gz \
    && mv libtokenizers.a /usr/lib/ \
    && rm libtokenizers.linux-amd64.tar.gz

# Download ONNX Runtime for linking
RUN wget -q https://github.com/microsoft/onnxruntime/releases/download/v1.24.4/onnxruntime-linux-x64-1.24.4.tgz \
    && tar -xzf onnxruntime-linux-x64-1.24.4.tgz \
    && cp onnxruntime-linux-x64-1.24.4/lib/libonnxruntime.so /usr/lib/ \
    && ldconfig

RUN GOOS=linux CGO_LDFLAGS="-L/usr/lib" go build -tags ORT -ldflags="-s -w -X 'main.Commit=${COMMIT_TXT}' -X 'main.BuildTime=${BUILD_DATE}' -X 'main.BuildEnv=${BUILD_ENV}'" -o /ws cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /usr/lib/libonnxruntime.so /usr/lib/
COPY --from=build /ws /ws

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/ws"]
