##
## Build 
##
FROM golang:1.16-alpine AS repostats_builder 
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.com.cn,direct
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o repostats

##
## Deploy
##
FROM alpine:latest
RUN apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk del tzdata
WORKDIR /app
COPY --from=repostats_builder /app/repostats .
EXPOSE 9103
ENTRYPOINT ["/app/repostats","repostats.ini"]

