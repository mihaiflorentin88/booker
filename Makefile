
# Go project-related variables
BINARY_NAME=booker
MAIN_GO_FILE=main.go
APP_ID=com.app.booker

# OS and architecture-specific binaries
BUILD_DIR=bin
WINDOWS_AMD64_BINARY=$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
WINDOWS_ARM64_BINARY=$(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe
LINUX_AMD64_BINARY=$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
LINUX_ARM64_BINARY=$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64
MAC_AMD64_BINARY=$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
MAC_ARM64_BINARY=$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64

# Default compile commands for each OS
compile:
	@echo "Compiling for all supported platforms..."
	# Windows
	#GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_AMD64_BINARY) $(MAIN_GO_FILE)
	#GOOS=windows GOARCH=arm64 go build -o $(WINDOWS_ARM64_BINARY)
	# Linux
	#GOOS=linux GOARCH=amd64 go build -o $(LINUX_AMD64_BINARY) $(MAIN_GO_FILE)
	#GOOS=linux GOARCH=arm64 go build -o $(LINUX_ARM64_BINARY) $(MAIN_GO_FILE)
	# macOS
	#GOOS=darwin GOARCH=arm64 go build -o $(MAC_ARM64_BINARY) $(MAIN_GO_FILE)
	GOOS=darwin GOARCH=amd64 go build -o $(MAC_AMD64_BINARY) $(MAIN_GO_FILE)

# Docker-based cross-compilation for all OS versions
docker-compile:
	@echo "Compiling using Docker for all supported platforms and architectures..."
	go install github.com/fyne-io/fyne-cross@latest

	fyne-cross windows -arch amd64 -app-id $(APP_ID) -output booker-windows-amd64.exe
	fyne-cross windows -arch arm64 -app-id $(APP_ID) -output booker-windows-arm64.exe
	#fyne-cross linux -arch amd64 -app-id $(APP_ID) -output booker-linux-amd64
	#fyne-cross linux -arch arm64 -app-id $(APP_ID) -output booker-linux-arm64
	#fyne-cross darwin -arch amd64 -app-id $(APP_ID) --dependencies "go install fyne.io/fyne/v2/cmd/fyne@latest" -output booker-darwin-amd64
	#fyne-cross darwin -arch arm64 -app-id $(APP_ID) --dependencies "go install fyne.io/fyne/v2/cmd/fyne@latest" -output booker-darwin-arm64

clean:
	rm -rf $(BUILD_DIR)

.PHONY: compile docker-compile clean