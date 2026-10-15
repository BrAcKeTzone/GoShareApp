#!/bin/bash
set -e

BINARY_NAME="goshare"
INSTALL_DIR="/usr/local/bin"

# Ensure script is run with sudo/root privileges
if [ "$EUID" -ne 0 ]; then
  echo "❌ Please run this installer as root (e.g., sudo ./install.sh)"
  exit 1
fi

echo "🔍 Checking for binary..."
if [ ! -f "./$BINARY_NAME" ]; then
  echo "❌ Error: '$BINARY_NAME' binary not found in the current directory. Please build it first."
  exit 1
fi

echo "📦 Installing $BINARY_NAME to $INSTALL_DIR..."
cp "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Generate the uninstaller script globally
cat << 'EOF' > "$INSTALL_DIR/goshare-uninstall"
#!/bin/bash
if [ "$EUID" -ne 0 ]; then
  echo "❌ Please run uninstallation as root (sudo goshare-uninstall)"
  exit 1
fi
echo "🗑️ Removing GoShare..."
rm -f /usr/local/bin/goshare
rm -f /usr/local/bin/goshare-uninstall
echo "✨ GoShare has been successfully uninstalled from your system!"
EOF
chmod +x "$INSTALL_DIR/goshare-uninstall"

echo ""
echo "🎉 Installation complete!"
echo "👉 You can now open a terminal in any directory and type: goshare"
echo "👉 To uninstall later, simply run: sudo goshare-uninstall"