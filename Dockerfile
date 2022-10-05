# build the server binary
FROM golang:1.14.4-alpine AS base-builder
WORKDIR /go/src/github.com/ceperapl/tasks
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server -tags release ./cmd

# copy the server binary from builder stage; run the server binary
FROM alpine:latest
WORKDIR /bin
COPY --from=base-builder /go/src/github.com/ceperapl/tasks/bin/server .
ENTRYPOINT ["server"]
