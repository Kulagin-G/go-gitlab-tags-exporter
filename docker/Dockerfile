FROM golang:1.20 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIT_TERMINAL_PROMPT=1 \
    GO111MODULE=on

COPY cmd ${GOPATH}/go-gitlab-tags-exporter/cmd
COPY internal ${GOPATH}/go-gitlab-tags-exporter/internal
COPY go.mod ${GOPATH}/go-gitlab-tags-exporter/
COPY go.sum ${GOPATH}/go-gitlab-tags-exporter/
COPY config ${GOPATH}/go-gitlab-tags-exporter/config
WORKDIR ${GOPATH}/go-gitlab-tags-exporter
RUN go mod tidy
RUN go build -ldflags="-s -w" -o go-gitlab-tags-exporter ./**/*.go

FROM scratch
COPY --from=builder /go/go-gitlab-tags-exporter /go/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /go

ENTRYPOINT ["./go-gitlab-tags-exporter"]