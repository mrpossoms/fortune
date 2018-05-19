SRC=$(wildcard *.go)

all:
	go build

run:
	go run $(SRC)

clean:
	rm -f fortune
