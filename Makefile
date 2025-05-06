CC=gcc
GO=go

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    INSTALL_MPG123 := sudo apt-get update && sudo apt-get install -y mpg123
else
    INSTALL_MPG123 := echo "Please install mpg123 manually from https://www.mpg123.de/download.shtml"
endif

all: build

build:
	GOOS=windows $(GO) build -o app.exe ./cmd/main.go

build-linux:
	GOOS=linux $(GO) build -o app ./cmd/main.go

run: build-linux check-mpg123
	./app

check-mpg123:
	@which mpg123 > /dev/null 2>&1 || (echo "Installing mpg123..." && $(INSTALL_MPG123))

clean:
	rm -f app.exe

.PHONY: all build clean run check-mpg123