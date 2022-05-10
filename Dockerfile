FROM golang:1.17

COPY . /go/src/app

WORKDIR /go/src/app/cmd

RUN go build -o finstar-api main.go

EXPOSE 8080

CMD ["./finstar-api"]
