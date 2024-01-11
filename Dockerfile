FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go get .
RUN go mod tidy

RUN go build -o /go-build

CMD [ "/go-build" ]