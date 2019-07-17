FROM golang as build

ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

ADD . /jiode

WORKDIR /jiode

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jiode

FROM alpine:3.7

ENV PORT=3000

RUN echo "http://mirrors.aliyun.com/alpine/v3.7/main/" > /etc/apk/repositories && \
    apk update && \
    apk add ca-certificates && \
    echo "hosts: files dns" > /etc/nsswitch.conf

WORKDIR /www

COPY --from=build /jiode/jiode /usr/bin/jiode

RUN chmod +x /usr/bin/jiode

ENTRYPOINT ["jiode"]