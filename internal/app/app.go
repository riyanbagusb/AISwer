package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"aiswer/internal/ai"
	"aiswer/internal/capture"
	"aiswer/internal/permissions"

	"github.com/getlantern/systray"
	"golang.design/x/hotkey"
)

// RunBackground starts the background hotkey listener and permission polling
func RunBackground() {
	go func() {
		if !permissions.CheckAccessibility() {
			permissions.CheckAndPromptAccessibility() // Ask once
			for !permissions.CheckAccessibility() {
				systray.SetTitle("! Perms")
				time.Sleep(2 * time.Second)
			}
		}
		
		// Reset title if API key is valid
		currentKey := os.Getenv("GEMINI_API_KEY")
		if currentKey != "" && currentKey != "your_gemini_api_key_here" {
			systray.SetTitle("-")
		}

		// Trigger Screen Recording permission prompt (if not already granted)
		_, _ = capture.Take()

		hk := setupHotkey()
		if err := hk.Register(); err != nil {
			log.Printf("failed to register hotkey: %v", err)
			systray.SetTitle("! Hotkey")
			return
		}
		defer hk.Unregister()

		fmt.Printf("AISwer is running in menu bar...\n")
		
		listenHotkey(hk)
	}()
}

func setupHotkey() *hotkey.Hotkey {
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

	fmt.Printf("Listening for hotkey: %s + %s\n", strings.ToUpper(modifiersStr), strings.ToUpper(keyStr))
	return hotkey.New(mods, key)
}

func listenHotkey(hk *hotkey.Hotkey) {
	ctx := context.Background()
	var resetTimer *time.Timer

	for {
		<-hk.Keydown()
		fmt.Println("Hotkey triggered! Processing...")

		if resetTimer != nil {
			resetTimer.Stop()
		}

		systray.SetTitle("...")

		imgBytes, err := capture.Take()
		if err != nil {
			log.Printf("Capture error: %v", err)
			systray.SetTitle("! Capture")
			resetTimer = time.AfterFunc(2*time.Second, func() {
				systray.SetTitle("-")
			})
			continue
		}

		currentKey := os.Getenv("GEMINI_API_KEY")
		if currentKey != "" && currentKey != "your_gemini_api_key_here" {
			answer, err := ai.SolveMCQ(ctx, currentKey, imgBytes)
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
}
