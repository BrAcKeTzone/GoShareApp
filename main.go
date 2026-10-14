package main

import (
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

func main() {
	if len(os.Args) < 2 || os.Args[1] != "goshare" {
		fmt.Println("Usage: goshare goshare")
		os.Exit(1)
	}

	port := "6060"
	ip := getLocalIP()
	url := fmt.Sprintf("http://%s:%s", ip, port)

	// Print ASCII logo, URL, then QR Code
	fmt.Println("\n" + asciiLogo + "\n")
	fmt.Printf("🚀 Server Active!\n")
	fmt.Printf("👉 URL: %s\n\n", url)

	q, err := qrcode.New(url, qrcode.Medium)
	if err == nil {
		fmt.Println(q.ToSmallString(false))
	}

	statFiles, _ := fs.Sub(content, "static")
	fileServer := http.FileServer(http.FS(statFiles))
	http.Handle("/", fileServer)

	http.HandleFunc("/api/list", handleListFiles)
	http.HandleFunc("/api/download", handleDownloadFile)
	http.HandleFunc("/api/view", handleViewFile)
	http.HandleFunc("/api/upload", handleUploadFile)
	http.HandleFunc("/api/mkdir", handleMkdir)
	http.HandleFunc("/api/icon", handleAppIcon)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

// Dynamically sets root base directory based on OS platform
func getBaseDirectory() string {
	if runtime.GOOS == "android" {
		androidStorage := os.Getenv("EXTERNAL_STORAGE")
		if androidStorage == "" {
			androidStorage = "/storage/emulated/0"
		}
		return androidStorage
	}

	// Desktop fallback to Current Working Directory
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
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
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
			if ip != nil {
				return ip.String()
			}
		}
	}
	return "127.0.0.1"
}

func handleListFiles(w http.ResponseWriter, r *http.Request) {
	subPath := r.URL.Query().Get("path")
	baseDir := getBaseDirectory()
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

func handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	baseDir := getBaseDirectory()
	filePath := filepath.Join(baseDir, filename)

	baseName := filepath.Base(filePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", baseName))
	http.ServeFile(w, r, filePath)
}

func handleViewFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	baseDir := getBaseDirectory()
	filePath := filepath.Join(baseDir, filename)
	http.ServeFile(w, r, filePath)
}

func handleMkdir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MkdirReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	baseDir := getBaseDirectory()
	newDirPath := filepath.Join(baseDir, req.Path, req.Name)
	err := os.MkdirAll(newDirPath, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(32 << 20)
	subPath := r.FormValue("path")
	files := r.MultipartForm.File["files"]
	baseDir := getBaseDirectory()
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
