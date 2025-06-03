ARG GOLANG_VERSION="1.24.3"

ARG PROJECT="go-gcp-pushover-notificationchannel"

ARG TARGETOS
ARG TARGETARCH

ARG COMMIT
ARG VERSION

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

ARG PROJECT
WORKDIR /${PROJECT}

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/server cmd/server
COPY pushover pushover
COPY function.go function.go
COPY incident_type.go incident_type.go
COPY incident_type_test.go incident_type_test.go
COPY message.go message.go
COPY template.go template.go

ARG TARGETOS
ARG TARGETARCH

ARG COMMIT
ARG VERSION

RUN BUILD_TIME=$(date +%s) && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -a \
    -installsuffix cgo \
    -ldflags "-X 'main.BuildTime=${BUILD_TIME}' -X 'main.GitCommit=${COMMIT}' -X 'main.OSVersion=${VERSION}'" \
    -o /bin/server \
    ./cmd/server


FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.source="https://github.com/DazWilkin/go-gcp-pushover-notificationchannel"

COPY --from=build /bin/server /

ENV PUSHOVER_USERKEY=""
ENV PUSHOVER_TOKEN=""

ENV PORT="8080"

ENTRYPOINT ["/server"]
