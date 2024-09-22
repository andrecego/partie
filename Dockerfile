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

# install yt-dlp-youtube-oauth2 plugin
RUN mkdir /etc/yt-dlp
RUN mkdir /etc/yt-dlp/plugins
RUN curl -L https://github.com/coletdjnz/yt-dlp-youtube-oauth2/releases/download/v2024.9.14/yt-dlp-youtube-oauth2.zip -o /etc/yt-dlp/plugins/yt-dlp-youtube-oauth2.zip

RUN mkdir /server
WORKDIR /server
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
