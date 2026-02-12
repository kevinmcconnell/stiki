build:
	CGO_ENABLED=0 go build -o bin/stiki .

docker:
	docker build -t stiki .

clean:
	rm -rf bin

.PHONY: build docker clean
