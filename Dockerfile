FROM golang:1.21-alpine as build

WORKDIR /app

ADD go.* ./
RUN go mod download
COPY bot.go .
RUN go build -o /bot


FROM alpine:3.17
WORKDIR /app
COPY --from=build /bot .
RUN apk add chromium
CMD [ "/app/bot", "localhost" ]