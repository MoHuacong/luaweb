.PHONY: build

build:
	go fmt . && go build -o luaweb cmd/cmd.go && chmod 777 luaweb && ./luaweb
