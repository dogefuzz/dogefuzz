FROM golang:1.19-alpine3.17 as builder

RUN apk add build-base

WORKDIR /app

ADD . .

RUN go build ./cmd/dogefuzz


FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/dogefuzz .

EXPOSE 3456

CMD ["/app/dogefuzz"]
