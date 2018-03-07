FROM golang:1.9-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

WORKDIR /go/src/github.com/kord-network/go-kord
COPY . .
RUN go install -v ./cmd/kord

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/kord /usr/local/bin/

CMD ["kord"]
