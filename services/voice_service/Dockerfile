FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY voice voice

CMD ["./voice"]
