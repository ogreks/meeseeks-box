FROM golang:1.21.0 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

## build
FROM scratch as meeseeks-box

WORKDIR /app

## set time zone
COPY --from=builder /usr/share/zoneinfo/UTC /etc/localtime

## set builder
COPY --from=builder /app/meeseeks-box .
COPY --from=builder /app/config.yml.dev /etc/meeseeks-box/config.yml

EXPOSE 8088

ENTRYPOINT ["./meeseeks-box", "server", "start", "-c", "/etc/meeseeks-box/config.yml"]

