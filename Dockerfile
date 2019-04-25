FROM golang:latest

RUN apt-get update

# Create project folder
RUN mkdir -p /go/src/github.com/hermesespinola/moneio
WORKDIR /go/src/github.com/hermesespinola/moneio

# Install dep
RUN go get -u github.com/golang/dep/cmd/dep

# Add dependencies
COPY Gopkg.toml Gopkg.lock  ./
RUN dep ensure -vendor-only

# Add project code
COPY *.go ./

# Run project
RUN GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -a -o bin/server *.go
CMD [ "bin/server" ]
EXPOSE 8080

LABEL "maintainer"="Hermes Espínola <hermes.espinola@gmail.com>"
LABEL version="1.0"
