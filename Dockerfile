FROM golang:1.23.2-alpine3.20 AS builder

COPY . /github.com/dbulyk/auth/source
WORKDIR /github.com/dbulyk/auth/source

RUN go mod download
RUN go build -o ./bin/auth cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/dbulyk/auth/source/bin/auth /root/

ENTRYPOINT ["./auth"]
CMD []