FROM golang

RUN go env GOPATH
RUN mkdir -p /go/src/github.com/AITestingOrg/notification-service
WORKDIR /go/src/github.com/AITestingOrg/notification-service

COPY . .

RUN go get -v ./cmd/notification-service
RUN go build ./cmd/notification-service
RUN export PATH=$PATH:/go/bin
EXPOSE 32700

CMD ["./notification-service"]