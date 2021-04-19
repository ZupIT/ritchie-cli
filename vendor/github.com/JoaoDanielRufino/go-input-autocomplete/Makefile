.PHONY:all
all: clean test build

.PHONY:build
build: go-fmt go-vet
	@echo "building the autocomplete executable"
	@echo "gopath being used is ${GOPATH}"
	@echo "goroot being used is ${GOROOT}"
	go build -x -v -o gia cmd/main.go
	@echo "executable created gia"

.PHONY:test
test:
	@echo "testing the autocomplete package"
	go test -cover -v -timeout 60s -coverprofile=coverage.out ./...
	# go cover works well when package is cloned in $GOPATH/src/ directory
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY:clean
clean:
	@echo "cleaning the autocomplete package"
	rm -f gia coverage.out coverage.html

.PHONY:go-fmt
go-fmt:
	@echo "formatting the go files using gofmt"
	gofmt -l -w  .

.PHONY:go-vet
go-vet:
	@echo "go vet of the package"
	go vet -v ./...