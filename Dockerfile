# Build the manager binary
FROM golang:1.16.5 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go ./
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/

# Build
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go
RUN go get github.com/markbates/pkger/cmd/pkger
RUN pkger && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o cnvrg-operator main.go pkged.go

FROM ubuntu:20.04
WORKDIR /opt/app-root
COPY --from=builder /workspace/cnvrg-operator .
