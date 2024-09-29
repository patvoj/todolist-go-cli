# Variables
BINARY_NAME = go-todolist
SOURCE_FILE = main.go

run:
	@go run $(SOURCE_FILE)

build:
	@go build -o $(BINARY_NAME) $(SOURCE_FILE)

clean:
	@rm -f $(BINARY_NAME)

start: build
	@./$(BINARY_NAME)

fmt:
	@go fmt ./...
