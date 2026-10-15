package main

import (
	"context"
	"encoding/json"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/skip2/go-qrcode"
)

//go:embed static/*
var content embed.FS

type FileInfo struct {
	Name  string `json:"Name"`
	Size  int64  `json:"Size"`
	IsDir bool   `json:"IsDir"`
}

type MkdirReq struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

const asciiLogo = `  ██████╗  ██████╗ ███████╗██╗  ██╗ █████╗ ██████╗ ███████╗
 ██╔════╝ ██╔═══██║██╔════╝██║  ██║██╔══██╗██╔══██╗██╔════╝
 ██║  ███╗██║   ██║███████║███████║███████║███████╔█████╗  
 ██║   ██║██║   ██║╚════██║██╔══██║██╔══██║██╔══██╗██╔══╝  
 ╚██████╔╝╚██████╔╝███████║██║  ██║██║  ██║██║  ██║███████╗
  ╚═════╝  ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝`

var serverInstance *http.Server
var serverRunning = false

func StartServer(isAndroid bool) {
	port := "6060"
	ip := getLocalIP()
	url := fmt.Sprintf("http://%s:%s", ip, port)

	fmt.Println("\n" + asciiLogo + "\n")
	fmt.Printf("🚀 GOSHARE Server Active!\n")
	fmt.Printf("👉 URL: %s\n\n", url)

	q, err := qrcode.New(url, qrcode.Medium)
	if err == nil {
		fmt.Println(q.ToSmallString(false))
	}

	mux := http.NewServeMux()

	statFiles, _ := fs.Sub(content, "static")
	mux.Handle("/", http.FileServer(http.FS(statFiles)))

	mux.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) { handleListFiles(w, r, isAndroid) })
	mux.HandleFunc("/api/download", func(w http.ResponseWriter, r *http.Request) { handleDownloadFile(w, r, isAndroid) })
	mux.HandleFunc("/api/view", func(w http.ResponseWriter, r *http.Request) { handleViewFile(w, r, isAndroid) })
	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) { handleUploadFile(w, r, isAndroid) })
	mux.HandleFunc("/api/mkdir", func(w http.ResponseWriter, r *http.Request) { handleMkdir(w, r, isAndroid) })
	mux.HandleFunc("/api/icon", handleAppIcon)
	
	// Status endpoint for checking if server is active
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"running": serverRunning})
	})

	// Dynamic Stop endpoint for toggle control
	mux.HandleFunc("/api/stop", func(w http.ResponseWriter, r *http.Request) {
		serverRunning = false
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server stopped"))
		go func() {
			if serverInstance != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				serverInstance.Shutdown(ctx)
			}
		}()
	})

	serverInstance = &http.Server{Addr: ":" + port, Handler: mux}
	serverRunning = true

	if err := serverInstance.ListenAndServe(); err != http.ErrServerClosed {
		serverRunning = false
		fmt.Printf("Server error: %v\n", err)
	}
}

func getBaseDirectory(isAndroid bool) string {
	if isAndroid || runtime.GOOS == "android" {
		androidStorage := os.Getenv("EXTERNAL_STORAGE")
		if androidStorage == "" {
			androidStorage = "/storage/emulated/0"
		}
		return androidStorage
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

func getLocalIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}

	var candidateIP string

	for _, iface := range interfaces {
		// Skip down interfaces or loopbacks
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		name := iface.Name
		// Skip virtual adapters to prevent choosing WSL/Docker/VPN IPs
		if isVirtualAdapter(name) {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			ipStr := ip.String()

			// Prioritize standard local Wi-Fi / router ranges
			if hasPrefix(ipStr, "192.168.") || hasPrefix(ipStr, "10.") {
				return ipStr
			}

			// Keep as fallback if it's another local network format (e.g. 172.x)
			if candidateIP == "" {
				candidateIP = ipStr
			}
		}
	}

	if candidateIP != "" {
		return candidateIP
	}
	return "127.0.0.1"
}

func isVirtualAdapter(name string) bool {
	virtualKeywords := []string{
		"vEthernet", "VirtualBox", "VMware", "WSL", 
		"docker", "tun", "tap", "VPN", "Bluetooth",
	}
	for _, kw := range virtualKeywords {
		if len(name) >= len(kw) && hasSubstring(name, kw) {
			return true
		}
	}
	return false
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func hasSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func handleListFiles(w http.ResponseWriter, r *http.Request, isAndroid bool) {
	subPath := r.URL.Query().Get("path")
	baseDir := getBaseDirectory(isAndroid)
	targetDir := filepath.Join(baseDir, subPath)

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		name := entry.Name()
		if name == ".git" || name == "goshare" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:  name,
			Size:  info.Size(),
			IsDir: entry.IsDir(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func handleDownloadFile(w http.ResponseWriter, r *http.Request, isAndroid bool) {
	filename := r.URL.Query().Get("file")
	baseDir := getBaseDirectory(isAndroid)
	filePath := filepath.Join(baseDir, filename)

	baseName := filepath.Base(filePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", baseName))
	http.ServeFile(w, r, filePath)
}

func handleViewFile(w http.ResponseWriter, r *http.Request, isAndroid bool) {
	filename := r.URL.Query().Get("file")
	baseDir := getBaseDirectory(isAndroid)
	filePath := filepath.Join(baseDir, filename)
	http.ServeFile(w, r, filePath)
}

func handleMkdir(w http.ResponseWriter, r *http.Request, isAndroid bool) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req MkdirReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	baseDir := getBaseDirectory(isAndroid)
	newDirPath := filepath.Join(baseDir, req.Path, req.Name)
	if err := os.MkdirAll(newDirPath, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleUploadFile(w http.ResponseWriter, r *http.Request, isAndroid bool) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(32 << 20)
	subPath := r.FormValue("path")
	files := r.MultipartForm.File["files"]
	baseDir := getBaseDirectory(isAndroid)
	targetDir := filepath.Join(baseDir, subPath)

	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		dstPath := filepath.Join(targetDir, header.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			continue
		}
		defer dst.Close()

		io.Copy(dst, file)
	}
	w.WriteHeader(http.StatusOK)
}

func handleAppIcon(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Getwd()
	iconPath := filepath.Join(cwd, "static", "icon.png")
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		http.Error(w, "Icon not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, iconPath)
}