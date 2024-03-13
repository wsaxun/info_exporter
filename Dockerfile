FROM golang:1.20-alpine AS build-env

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /mfy

ADD . /mfy

RUN go mod tidy && \
    go build -o info_exporter

FROM alpine:3.15.6

WORKDIR /mfy

ENV TZ Asia/Shanghai

EXPOSE 8080

COPY --from=build-env /mfy/info_exporter .
COPY glibc-2.31-r0.apk glibc-bin-2.31-r0.apk /tmp/

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && apk add curl && apk add busybox-extras && \
    apk add --allow-untrusted /tmp/*.apk && \
    apk add tzdata && echo "${TZ}" > /etc/timezone && \
    ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime && \
    rm -rf /var/cache/apk/* /tmp/*.apk && \
    chmod +x ./info_exporter

CMD ["./info_exporter"]