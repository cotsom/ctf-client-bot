FROM golang:1.21-alpine as build

WORKDIR /app

ADD go.* ./
COPY config config
RUN go mod download
COPY bot.go .
RUN go build -o /bot


FROM alpine:3.17
WORKDIR /app
COPY --from=build /bot .
COPY config.yml config.yml
RUN apk add chromium
CMD [ "/app/bot"]