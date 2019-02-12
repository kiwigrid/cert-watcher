# Build the manager binary
FROM golang:1.11.3 as builder

# Copy in the go src
WORKDIR /go/src/github.com/kiwigrid/cert-watcher
COPY main.go .
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o watcher github.com/kiwigrid/cert-watcher

# Copy the controller-manager into a thin image
FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/kiwigrid/cert-watcher/watcher .
ENTRYPOINT ["./watcher"]
