FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /eCommerce_go

EXPOSE 8080

CMD [ "/eCommerce_go" ]