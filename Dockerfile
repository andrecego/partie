FROM golang:1.16

RUN apt-get update -qq && apt-get install -y \
  build-essential \
  ca-certificates \
  openssl \
  iputils-ping \
  && update-ca-certificates

RUN mkdir /server
WORKDIR /server
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .