FROM golang:1.20
COPY . /usr/local/app
WORKDIR /usr/local/app
RUN go mod tidy && go build -o IPCityServer main.go
