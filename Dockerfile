FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o search -mod=vendor -ldflags='-s -w'  -installsuffix cgo main.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /search
COPY --from=builder ./build/search ./cmd/


EXPOSE 80
ENTRYPOINT ["./cmd/search","-config=/configs/config.yml"]
