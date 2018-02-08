FROM golang:1.9-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

WORKDIR /go/src/github.com/meta-network/go-meta
COPY . .
RUN go install -v ./cmd/meta

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/meta /usr/local/bin/

CMD ["meta"]
