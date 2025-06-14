.PHONY: all build update-readme-help

all: build update-readme-help

build:
	go build

update-readme-help:
	./script/update-readme-help.sh
