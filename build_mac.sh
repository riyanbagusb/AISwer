#!/bin/bash

# Exit on error
set -e

APP_NAME="AISwer"
EXECUTABLE_NAME="aiswer"
DIST_DIR="dist"
APP_DIR="$DIST_DIR/$APP_NAME.app"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"

echo "Building Go binary..."
go build -o "$DIST_DIR/$EXECUTABLE_NAME" .

echo "Creating App Bundle structure..."
mkdir -p "$MACOS_DIR"

echo "Copying binary..."
cp "$DIST_DIR/$EXECUTABLE_NAME" "$MACOS_DIR/"

echo "Creating Info.plist..."
cat <<EOF > "$CONTENTS_DIR/Info.plist"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>$EXECUTABLE_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>com.bagus.aiswer</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
EOF

echo "Done! The app is located at: $APP_DIR"
