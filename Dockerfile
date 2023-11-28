FROM golang:1.21.3-alpine3.18 as builder

RUN apk update && \
    apk add --no-cache git openssh tzdata build-base python3 net-tools

WORKDIR /app

COPY .env.example .env
COPY . .

RUN go install github.com/buu700/gin@latest
RUN go mod tidy

RUN make build

FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk --update --no-cache add tzdata && \
    apk --no-cache add curl && \
    apk --no-cache add chromium

ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

EXPOSE 8002

COPY --from=builder /app /app

ENTRYPOINT ["/app/invoice-service"]
