FROM golang:1.23.1

RUN apt-get update -qq && apt-get install -y \
  build-essential \
  ca-certificates \
  openssl \
  iputils-ping \
  ffmpeg \
  && update-ca-certificates

# install yt-dlp
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp

RUN mkdir /server
WORKDIR /server

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
