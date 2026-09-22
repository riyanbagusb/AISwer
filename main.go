package main

import (
	_ "embed"
	"fmt"

	"aiswer/internal/app"
	"aiswer/internal/config"
	"aiswer/internal/ui"

	"github.com/getlantern/systray"
)

//go:embed .env
var embeddedEnv []byte

func main() {
	// Initialize systray, which blocks and runs the macOS event loop on the main thread.
	systray.Run(onReady, onExit)
}

func onReady() {
	// 1. Initialize environment configuration
	config.Init(embeddedEnv)

	// 2. Setup System Tray Menu UI
	ui.Setup()

	// 3. Start background logic (Permissions, Hotkeys, Screen Capture)
	app.RunBackground()
}

func onExit() {
	fmt.Println("Exiting AISwer Assistant...")
}
