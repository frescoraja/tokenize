build:
	rm -rf ./bin && mkdir -p ./bin
	go build -o ./bin/tokenize .

default: build
