FROM golang:1.19-alpine3.17 as builder

RUN apk add build-base

WORKDIR /app

ADD . .

RUN go build ./cmd/dogefuzz


FROM alpine:3.17

RUN apk add curl

# https://stackoverflow.com/questions/66963068/docker-alpine-executable-binary-not-found-even-if-in-path
RUN apk add gcompat

WORKDIR /app

COPY --from=builder /app/dogefuzz .
COPY config.json .

RUN mkdir assets
COPY assets/ ./assets/

EXPOSE 3456

CMD ["/app/dogefuzz"]
