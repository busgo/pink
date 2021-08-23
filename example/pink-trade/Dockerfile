FROM golang:1.16
MAINTAINER busgo "248434199@qq.com"

#env
ENV GOPROXY=https://goproxy.cn
ENV APP_NAME=pink-trade
ENV CONF_FILLE=trade-docker.yaml
ENV GOSUMDB=off

# copy
WORKDIR /build
COPY . /build
RUN go mod tidy
RUN go build -o ${APP_NAME}


# run
ENTRYPOINT ./$APP_NAME -conf=$CONF_FILLE