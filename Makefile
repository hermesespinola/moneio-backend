.PHONY: run build clean

run:
	go run *.go

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/server

clean:
	rm -rf ./bin ./vendor Gopkg.lock

# deploy: clean build
# 	sls deploy --verbose
