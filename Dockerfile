FROM golang:1.13.7
WORKDIR /go/src/wxapp
ADD . .
RUN export GOPROXY=https://goproxy.io
RUN go mod init app && \
go build -o run .
EXPOSE 8080
CMD ["/go/src/wxapp/run"]