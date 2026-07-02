.PHONY: build run clean dev

APP_NAME := vibeshare
BUILD_DIR := bin
GOPROXY ?= https://goproxy.cn,direct

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

build: $(BUILD_DIR)
	GOPROXY=$(GOPROXY) go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run: build
	./$(BUILD_DIR)/$(APP_NAME)

dev: build
	VIBESHARE_DEV=1 ./$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR) data
