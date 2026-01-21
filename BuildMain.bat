set CGO_ENABLED=1
set CC=gcc
set CXX=g++
set CGO_CFLAGS=-I%cd%/SDL2/include
set CGO_LDFLAGS=-L%cd%/SDL2/lib/ -lSDL2

Rem go clean -cache -modcache
go build -v ./VM-main
