FROM golang:1.14.4 as builder

WORKDIR /app
COPY . ./

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o ./bin/svc

FROM scratch

# Copy our static executable
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /app/bin/svc /svc

# Port on which the service will be exposed.
EXPOSE 8080 8888

# Run the svc binary.
CMD ["./svc"]
