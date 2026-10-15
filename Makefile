# Output directory
BUILD_DIR = builds

# Compiler flags for maximum size reduction
# -s: Omit symbol table and debug information
# -w: Omit DWARF debugging information
# -trimpath: Remove local file system paths from the binary
LDFLAGS = -ldflags="-s -w" -trimpath

.PHONY: all clean linux windows macos android

all: clean linux windows macos android
	@echo ""
	@echo "✨ All optimized builds completed inside the '$(BUILD_DIR)/' folder!"

clean:
	rm -rf $(BUILD_DIR)

# --- LINUX BUILDS ---
linux: linux-amd64 linux-arm64

linux-amd64:
	@echo "Building Linux (AMD64)..."
	mkdir -p $(BUILD_DIR)/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/linux-amd64/goshare .

linux-arm64:
	@echo "Building Linux (ARM64)..."
	mkdir -p $(BUILD_DIR)/linux-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/linux-arm64/goshare .


# --- WINDOWS BUILDS ---
windows: windows-amd64 windows-arm64

windows-amd64:
	@echo "Building Windows (AMD64)..."
	mkdir -p $(BUILD_DIR)/windows-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows-amd64/goshare.exe .

windows-arm64:
	@echo "Building Windows (ARM64)..."
	mkdir -p $(BUILD_DIR)/windows-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows-arm64/goshare.exe .


# --- MACOS BUILDS ---
macos: macos-amd64 macos-arm64

macos-amd64:
	@echo "Building macOS Intel (AMD64)..."
	mkdir -p $(BUILD_DIR)/macos-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/macos-amd64/goshare .

macos-arm64:
	@echo "Building macOS Apple Silicon (ARM64)..."
	mkdir -p $(BUILD_DIR)/macos-arm64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/macos-arm64/goshare .


# --- ANDROID BUILD ---
android: android-arm64

android-arm64:
	@echo "Building Android (ARM64)..."
	mkdir -p $(BUILD_DIR)/android-arm64
	CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/android-arm64/goshare .