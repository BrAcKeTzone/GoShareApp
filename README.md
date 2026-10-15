# GOSHARE

## Description

GOSHARE is a lightning-fast, zero-friction cross-platform file-sharing utility written in Go (Golang) with an embedded modern web interface. It turns any directory into a secure local web server, allowing you to seamlessly send, receive, preview, and download files over a local Wi-Fi network or hotspot without needing an active internet connection. Designed with a clean digital calculator aesthetic, it features instantaneous zero-argument startup for desktop environments.

---

## Installation & Uninstallation Steps

### Linux & macOS

- **How to Install:**

1. Extract the downloaded zip archive and navigate into the target build folder.
2. Run the installer script with root privileges:

```bash
sudo ./install.sh

```

- **How to Uninstall:**

1. Run the global uninstaller command from any terminal:

```bash
sudo goshare-uninstall

```

### Windows

- **How to Install:**

1. Extract the downloaded zip archive using File Explorer.
2. Right-click on **`install.bat`** and select **"Run as administrator"**.

- **How to Uninstall:**

1. Navigate to the installation directory (`C:\Program Files\GoShare`).
2. Right-click on **`uninstall.bat`** and select **"Run as administrator"**.

---

## How to Use & Stop

### Linux, macOS, & Windows

- **How to Use:**

1. Open your terminal (Command Prompt/PowerShell on Windows, or Terminal on Linux/macOS) and navigate to **any** directory or folder on your computer that you want to share.
2. Type the command:

```bash
goshare

```

3. The server will start up instantly, display an ASCII logo, and print a terminal-friendly QR code pointing to your local server URL.
4. Open any device connected to the same local Wi-Fi network, scan the QR code, or type the provided URL into your web browser to browse, upload, preview, or download files.

- **How to Stop:**

1. Return to the terminal window where `goshare` is running.
2. Press `Ctrl + C` on your keyboard to safely shut down the local server process.
