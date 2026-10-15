//go:build android
// +build android

package main

func main() {
	// Android mode: runs automatically without CLI args and shares device storage
	StartServer(true)
}