FROM golang:1.26-alpine AS build

ARG COMMIT_TXT
ARG BUILD_DATE
ARG BUILD_ENV
ARG SERVER_HOST
ARG JOB_TOKEN

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X 'main.Commit=${COMMIT_TXT}' -X 'main.BuildTime=${BUILD_DATE}' -X 'main.BuildEnv=${BUILD_ENV}'" -o /ws cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /ws /ws

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/ws"]
