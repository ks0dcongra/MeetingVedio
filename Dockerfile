# 第一階段: 建置階段
FROM golang:1.20.6-alpine3.18 as builder
LABEL maintainer="Shehomebow"
RUN mkdir -p /app
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go build -o meeting-video-helper

# 第二階段: 運行階段
FROM alpine:3.18
RUN apk update && apk add --no-cache git ffmpeg
COPY --from=builder /app /app
WORKDIR /app
EXPOSE 9528
ENTRYPOINT ["/app/meeting-video-helper"]
