SRC=$(wildcard *.go)

build:
	go build $(SRC):

run:
	go run $(SRC)
