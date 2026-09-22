package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"examocr/internal/ai"
	"examocr/internal/capture"

	"github.com/getlantern/systray"
	"github.com/joho/godotenv"
	"golang.design/x/hotkey"
)

//go:embed .env
var embeddedEnv []byte

func initEnv() {
	// 1. Load embedded .env file into environment
	if len(embeddedEnv) > 0 {
		if envMap, err := godotenv.UnmarshalBytes(embeddedEnv); err == nil {
			for k, v := range envMap {
				if os.Getenv(k) == "" {
					os.Setenv(k, v)
				}
			}
		}
	}

	// 2. Overload with external .env from executable directory if exists
	if exePath, err := os.Executable(); err == nil {
		_ = godotenv.Overload(filepath.Join(filepath.Dir(exePath), ".env"))
	}

	// 3. Overload with .env in current working directory if exists
	_ = godotenv.Overload()
}

func main() {
	// Initialize systray, which blocks and runs the macOS event loop on the main thread.
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("-")
	systray.SetTooltip("Exam OCR Assistant")

	mQuit := systray.AddMenuItem("Quit", "Quit Exam OCR")

	initEnv()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" || apiKey == "your_gemini_api_key_here" {
		systray.SetTitle("! Key")
		log.Println("Error: GEMINI_API_KEY is not set in .env file.")
	}

	// Setup Hotkey
	modifiersStr := os.Getenv("HOTKEY_MODIFIERS")
	keyStr := os.Getenv("HOTKEY_KEY")
	if modifiersStr == "" {
		modifiersStr = "cmd+shift"
	}
	if keyStr == "" {
		keyStr = "x"
	}

	var mods []hotkey.Modifier
	for _, m := range strings.Split(modifiersStr, "+") {
		switch strings.ToLower(m) {
		case "cmd":
			mods = append(mods, hotkey.ModCmd)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "ctrl":
			mods = append(mods, hotkey.ModCtrl)
		case "option", "alt":
			mods = append(mods, hotkey.ModOption)
		}
	}

	var key hotkey.Key
	switch strings.ToLower(keyStr) {
	case "x":
		key = hotkey.KeyX
	case "a":
		key = hotkey.KeyA
	case "s":
		key = hotkey.KeyS
	case "1":
		key = hotkey.Key1
	default:
		key = hotkey.KeyX
	}

	hk := hotkey.New(mods, key)

	// We run hotkey registration and listening in a goroutine
	// because onReady is called on a different goroutine (or main) and we don't want to block it forever.
	go func() {
		if err := hk.Register(); err != nil {
			log.Printf("failed to register hotkey: %v", err)
			systray.SetTitle("! Hotkey")
			return
		}
		defer hk.Unregister()

		fmt.Printf("Exam OCR Assistant is running in menu bar...\n")
		fmt.Printf("Listening for hotkey: %s + %s\n", strings.ToUpper(modifiersStr), strings.ToUpper(keyStr))

		ctx := context.Background()

		var resetTimer *time.Timer

		// Infinite loop waiting for hotkey events
		for {
			<-hk.Keydown()
			fmt.Println("Hotkey triggered! Processing...")

			if resetTimer != nil {
				resetTimer.Stop()
			}

			systray.SetTitle("...")

			// 1. Capture screen
			imgBytes, err := capture.Take()
			if err != nil {
				log.Printf("Capture error: %v", err)
				systray.SetTitle("! Capture")
				resetTimer = time.AfterFunc(2*time.Second, func() {
					systray.SetTitle("-")
				})
				continue
			}

			// 2. Ask AI
			if apiKey != "" && apiKey != "your_gemini_api_key_here" {
				answer, err := ai.SolveMCQ(ctx, apiKey, imgBytes)
				if err != nil {
					log.Printf("AI error: %v", err)
					systray.SetTitle("! AI")
					resetTimer = time.AfterFunc(2*time.Second, func() {
						systray.SetTitle("-")
					})
					continue
				}
				fmt.Printf("AI answered: %s\n", answer)
				systray.SetTitle(fmt.Sprintf("%s", answer))
			} else {
				systray.SetTitle("! Key")
			}

			resetTimer = time.AfterFunc(2*time.Second, func() {
				systray.SetTitle("-")
			})
		}
	}()

	// Handle Quit
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// Clean up here if needed
	fmt.Println("Exiting Exam OCR Assistant...")
}
