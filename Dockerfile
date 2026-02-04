FROM golang:1.21 AS build

ARG COMMIT_TXT
ARG BUILD_DATE
ARG BUILD_ENV
ARG SERVER_HOST
ARG JOB_TOKEN

RUN git config --global url.https://gitlab-ci-token:${JOB_TOKEN}@${SERVER_HOST}.insteadOf https://${SERVER_HOST}
RUN export GOPRIVATE=${SERVER_HOST}

WORKDIR $GOPATH/src/
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X 'main.Commit=${COMMIT_TXT}' -X 'main.BuildTime=${BUILD_DATE}' -X 'main.BuildEnv=${BUILD_ENV}'" -o /ws cmd/main.go

FROM alpine:3.18

ENV USER=pdaccess
ENV GROUPNAME=$USER
ENV UID=10001
ENV GID=10001

RUN addgroup \
    --gid "$GID" \
    "$GROUPNAME" \
&&  adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$GROUPNAME" \
    --no-create-home \
    --uid "$UID" \
    $USER

WORKDIR /

COPY --from=build /ws /ws

RUN apk add --no-cache libc6-compat tcpdump gcompat
RUN chmod 755 /ws

EXPOSE 8080

USER pdaccess:pdaccess

ENTRYPOINT ["/ws"]