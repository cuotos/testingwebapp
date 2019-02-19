FROM golang:1.11 as builder

COPY app /app

WORKDIR /app

RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /webtest

FROM alpine

COPY --from=builder /webtest /webtest

CMD ["/webtest"]