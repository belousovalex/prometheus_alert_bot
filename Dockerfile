FROM golang:1.13
COPY . /src/
WORKDIR /src/
RUN go mod download
RUN mkdir /dst
RUN CGO_ENABLED=0 GOOS=linux go build -o /dst/bot cmd/bot.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /dst/bot .
CMD ./bot /etc/prometheus_alert_bot/config.yml