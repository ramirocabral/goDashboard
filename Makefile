BINARY_NAME = golang-system-monitor
BUILD_DIR = ./target
BUILD_PATH = $(BUILD_DIR)/$(BINARY_NAME)

.DEFAULT_GOAL := build

build:
	go build -o $(BUILD_PATH) ./cmd/server/main.go

run: build
	 ./$(BUILD_PATH)
	
clean:
	rm -rf $(BUILD_DIR)
