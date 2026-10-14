all:
	@echo "Building all platforms..."
	mkdir -p builds/linux-amd64 && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o builds/linux-amd64/goshare .
	mkdir -p builds/windows-amd64 && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o builds/windows-amd64/goshare.exe .
	mkdir -p builds/macos-amd64 && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o builds/macos-amd64/goshare .
	mkdir -p builds/macos-arm64 && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o builds/macos-arm64/goshare .
	mkdir -p builds/android-arm64 && CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -o builds/android-arm64/goshare .
	@echo "Done! Check the 'builds/' directory."