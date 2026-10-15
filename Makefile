.PHONY: all clean linux windows macos android

all: clean linux windows macos android
	@echo "✨ All platform builds completed successfully inside the 'builds/' folder!"

clean:
	rm -rf builds/

linux:
	mkdir -p builds/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o builds/linux-amd64/goshare .

windows:
	mkdir -p builds/windows-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o builds/windows-amd64/goshare.exe .

macos:
	mkdir -p builds/macos-amd64 builds/macos-arm64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o builds/macos-amd64/goshare .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o builds/macos-arm64/goshare .

android:
	mkdir -p builds/android-arm64
	CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -o builds/android-arm64/goshare .