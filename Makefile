.PHONY: all clean amd64 arm64 universal

BINARY_NAME=rst2md
BUILD_DIR=bin

all: universal

clean:
	rm -rf $(BUILD_DIR)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

amd64: $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_amd64 ./cmd/rst2md

arm64: $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_arm64 ./cmd/rst2md

universal: amd64 arm64
	lipo -create -output $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/$(BINARY_NAME)_amd64 $(BUILD_DIR)/$(BINARY_NAME)_arm64
	rm $(BUILD_DIR)/$(BINARY_NAME)_amd64 $(BUILD_DIR)/$(BINARY_NAME)_arm64
	lipo -info $(BUILD_DIR)/$(BINARY_NAME)
