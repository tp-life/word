#FROM golang:1.14.4-alpine3.12 AS builder
#
#WORKDIR /build
#
#ENV GOPROXY https://goproxy.cn
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
#  apk add --no-cache upx ca-certificates tzdata
#
#COPY go.mod .
#COPY go.sum .
#RUN go mod download
#
#COPY . .
#RUN CGO_ENABLED=0 go  build -o helloword -ldflags "-X 'word/pkg/app.GinMode=debug' -s -w" -tags doc cmd/main.go && \
#      upx --best helloword -o _upx_server && \
#      mv -f _upx_server helloword

#FROM alpine:3.12 as runner
#ENV GIN_MODE release
#
#WORKDIR /app
#COPY --from=builder /build/helloword /app
#COPY --from=builder /build/configs /app/configs
#COPY --from=builder /build/locales /app/locales
#COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#
#EXPOSE 8080
#CMD ["./helloword", "api", "start", "--daemon=false"]


FROM alpine:3.12 as runner
ENV GIN_MODE release
ENV PORT 8080

WORKDIR /app

COPY  ./helloword /app
COPY  ./configs /app/configs
COPY  ./locales /app/locales
COPY  ./Shanghai /etc/localtime
COPY  ./ca-certificates.crt /etc/ssl/certs/


EXPOSE 8080
CMD ["./helloword", "api", "start", "--daemon=false"]