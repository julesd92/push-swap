all: checker push-swap

checker:
	go build -o checker ./checker

push-swap:
	go build -o push-swap ./push-swap

test:
	go vet ./...
	go test ./...

clean:
	rm -f checker push-swap

.PHONY: all checker push-swap test clean
