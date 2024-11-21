FROM golang:1.23-alpine AS builder

RUN go version
ENV GOPATH=/

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o voting_app ./main.go

CMD [ "./voting_app" ]