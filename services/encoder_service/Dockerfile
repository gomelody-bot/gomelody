FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY encoder encoder

CMD ["./encoder"]
