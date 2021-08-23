FROM golang:1.16
MAINTAINER busgo "248434199@qq.com"

#env
ENV GOPROXY=https://goproxy.cn
ENV HTTP_PORT=8001
ENV APP_NAME=pink
ENV CONF_FILLE=app-docker.yaml

# copy
WORKDIR /build
COPY . /build
RUN go mod tidy
RUN go build -o ${APP_NAME}

# expose
EXPOSE ${HTTP_PORT}

# run
ENTRYPOINT ./$APP_NAME -conf=$CONF_FILLE