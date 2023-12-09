FROM golang:1.21.0 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

## two build
FROM scratch as meeseeks-box

WORKDIR /app

COPY --from=builder /app/meeseeks-box .
COPY --from=builder /app/config.yaml.dev ./config.yaml

EXPOSE 80

ENTRYPOINT ["./app/meeseeks/meeseeks-box", "server", "start"]