.PHONY: build run clean

all: build

build:
	go build -o app cmd/main.go

run: build
	./app

clean:
	rm -f app