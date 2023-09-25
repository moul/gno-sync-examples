all: test install

install:
	go install .

test:
	go test -v .
