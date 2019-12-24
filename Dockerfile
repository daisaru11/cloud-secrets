FROM golang:1.13 as builder

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY injector/ injector/
COPY webhook/ webhook/
COPY decoder/ decoder/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o cloud-secrets main.go

FROM debian:buster-slim
WORKDIR /
COPY --from=builder /workspace/cloud-secrets /usr/local/bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER 65534

ENTRYPOINT ["/usr/local/bin/cloud-secrets"]
CMD ["controller"]
