FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /etc/frp

COPY frps /usr/local/bin/frps
COPY frps.toml /etc/frp/frps.toml

RUN chmod +x /usr/local/bin/frps

EXPOSE 7000
EXPOSE 80
EXPOSE 60002

ENTRYPOINT ["/usr/local/bin/frps", "-c", "/etc/frp/frps.toml"]