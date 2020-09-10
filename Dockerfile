FROM golang:1.14

WORKDIR /go/src/marzipango
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["marzipango", "-hostname=0.0.0.0"]