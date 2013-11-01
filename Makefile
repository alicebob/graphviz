.PHONY: all graphviz test
all: graphviz test

graphviz:
	go build

test:
	go test
