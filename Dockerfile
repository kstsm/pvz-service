FROM golang:1.24.2

WORKDIR ${GOPATH}/pvz-service/
COPY . ${GOPATH}/pvz-service/

RUN go build -o /build ./ && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]