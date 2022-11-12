FROM golang:1.13.15-alpine3.12 as builder

RUN apk add build-base

# Add project's files
WORKDIR /fuzzer
ADD . .
# Test project
#RUN go test ./...
# Build executable
RUN go build ./cmd/contractfuzzer


# Execute 
FROM alpine:3.12
WORKDIR /exec_env
COPY --from=builder /fuzzer/contractfuzzer .
COPY --from=builder /fuzzer/fuzzer_run.sh .
EXPOSE 8888
CMD ["/bin/sh"]
