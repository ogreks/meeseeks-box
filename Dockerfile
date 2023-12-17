FROM golang:1.21.0 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

## build
FROM scratch as meeseeks-box

WORKDIR /app

COPY --from=builder /app/meeseeks-box .
COPY --from=builder /app/config.yml.dev ./config.yml

EXPOSE 80

ENTRYPOINT ["./meeseeks-box", "server", "start"]

