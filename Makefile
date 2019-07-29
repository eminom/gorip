all:
	go build -o rip
	GOOS=linux GOARCH=amd64 go build -o build/rip-linux_amd64

format:
	find . -type f -name "*.go" | xargs -i{} gofmt -w {}
	
.PHONY: clean
clean:
	find . -type f -name "*.log" | xargs -i{} rm {}
	rm -f gorip
	rm -f rip
	rm -f rip-linux_amd64
