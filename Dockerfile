# Build the manager binary
FROM golang:1.19.1 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/
COPY cmd/ cmd/

# Build
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o cnvrg-operator cmd/operator/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o copctl cmd/copctl/*.go


FROM registry.access.redhat.com/ubi9/ubi:latest
LABEL name="cnvrg-operator" \
      vendor="cnvrg.io" \
      version="4.3.16" \
      release="4.3.16" \
      summary="Cnvrg K8s Operator" \
      description="cnvrg.io AIOS patform"
RUN dnf update -y
USER 1000
WORKDIR /opt/app-root
COPY license /licenses
COPY --from=builder /workspace/cnvrg-operator .
COPY --from=builder /workspace/copctl .

