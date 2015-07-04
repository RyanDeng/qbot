
all:
	cd go; go install -ldflags  -v ./...

install: all
	@echo

test:
	cd go; CGO_ENABLED=0 go test ./...

clean:
	cd go; go clean -i ./...



