FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY kpulse /usr/local/bin/kpulse
ENTRYPOINT ["kpulse"]
