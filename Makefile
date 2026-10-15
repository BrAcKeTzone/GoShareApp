# Output directory
BUILD_DIR = builds

# Compiler flags for maximum size reduction
# -s: Omit symbol table and debug information
# -w: Omit DWARF debugging information
# -trimpath: Remove local file system paths from the binary
LDFLAGS = -ldflags="-s -w" -trimpath

.PHONY: all clean linux windows macos

all: clean linux windows macos
	@echo ""
	@echo "✨ All optimized platform builds and zip archives completed inside the '$(BUILD_DIR)/' folder!"

clean:
	rm -rf $(BUILD_DIR)

# --- LINUX BUILDS ---
linux: linux-amd64 linux-arm64

linux-amd64:
	@echo "Building Linux (AMD64)..."
	mkdir -p $(BUILD_DIR)/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/linux-amd64/goshare .
	cp installer/installer.sh $(BUILD_DIR)/linux-amd64/install.sh
	chmod +x $(BUILD_DIR)/linux-amd64/install.sh
	cd $(BUILD_DIR) && zip -r goshare-linux-amd64.zip linux-amd64

linux-arm64:
	@echo "Building Linux (ARM64)..."
	mkdir -p $(BUILD_DIR)/linux-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/linux-arm64/goshare .
	cp installer/installer.sh $(BUILD_DIR)/linux-arm64/install.sh
	chmod +x $(BUILD_DIR)/linux-arm64/install.sh
	cd $(BUILD_DIR) && zip -r goshare-linux-arm64.zip linux-arm64


# --- WINDOWS BUILDS ---
windows: windows-amd64 windows-arm64

windows-amd64:
	@echo "Building Windows (AMD64)..."
	mkdir -p $(BUILD_DIR)/windows-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows-amd64/goshare.exe .
	cp installer/installer.bat $(BUILD_DIR)/windows-amd64/install.bat
	cd $(BUILD_DIR) && zip -r goshare-windows-amd64.zip windows-amd64

windows-arm64:
	@echo "Building Windows (ARM64)..."
	mkdir -p $(BUILD_DIR)/windows-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows-arm64/goshare.exe .
	cp installer/installer.bat $(BUILD_DIR)/windows-arm64/install.bat
	cd $(BUILD_DIR) && zip -r goshare-windows-arm64.zip windows-arm64


# --- MACOS BUILDS ---
macos: macos-amd64 macos-arm64

macos-amd64:
	@echo "Building macOS Intel (AMD64)..."
	mkdir -p $(BUILD_DIR)/macos-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/macos-amd64/goshare .
	cp installer/installer.sh $(BUILD_DIR)/macos-amd64/install.sh
	chmod +x $(BUILD_DIR)/macos-amd64/install.sh
	cd $(BUILD_DIR) && zip -r goshare-macos-amd64.zip macos-amd64

macos-arm64:
	@echo "Building macOS Apple Silicon (ARM64)..."
	mkdir -p $(BUILD_DIR)/macos-arm64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/macos-arm64/goshare .
	cp installer/installer.sh $(BUILD_DIR)/macos-arm64/install.sh
	chmod +x $(BUILD_DIR)/macos-arm64/install.sh
	cd $(BUILD_DIR) && zip -r goshare-macos-arm64.zip macos-arm64