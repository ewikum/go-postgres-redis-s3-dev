FROM golang:1.9.2-alpine

ARG app_env
ENV APP_ENV $app_env

RUN apk add --no-cache git mercurial

# download go basic dependencies
RUN go get github.com/lib/pq
RUN go get github.com/minio/minio-go
RUN go get github.com/pilu/fresh
RUN go get github.com/garyburd/redigo/redis

COPY ./go-app /go/src/app
WORKDIR /go/src/app

RUN go get ./...

RUN go build

CMD fresh
EXPOSE 8080