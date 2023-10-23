FROM golang:1.18

WORKDIR /worker

COPY main.go go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

ENV GOPROXY https://goproxy.cn,direct

EXPOSE 8888

RUN go mod download
RUN go build main.go

ENTRYPOINT ./main -port $SERVER_PORT