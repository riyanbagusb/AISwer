package ui

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"aiswer/internal/ai"
	"aiswer/internal/config"

	"github.com/getlantern/systray"
)

var (
	mAPIKey    *systray.MenuItem
	mModel     *systray.MenuItem
	mQuit      *systray.MenuItem
	modelItems []*systray.MenuItem
)

// Setup creates the systray menu items and starts their event listeners.
func Setup() {
	systray.SetTitle("-")
	systray.SetTooltip("AISwer")

	mAPIKey = systray.AddMenuItem("Set API Key...", "Set Gemini API Key")
	mModel = systray.AddMenuItem("Model", "Choose Gemini Model")
	systray.AddSeparator()
	mQuit = systray.AddMenuItem("Quit", "Quit AISwer")

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" || apiKey == "your_gemini_api_key_here" {
		systray.SetTitle("! Key")
	}

	go refreshModelsMenu(apiKey)
	go handleEvents()
}

func refreshModelsMenu(key string) {
	for _, item := range modelItems {
		item.Hide()
	}
	modelItems = nil

	if key == "" || key == "your_gemini_api_key_here" {
		item := mModel.AddSubMenuItem("No API Key", "")
		item.Disable()
		modelItems = append(modelItems, item)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	models, err := ai.FetchModels(ctx, key)
	if err != nil {
		item := mModel.AddSubMenuItem("Error fetching models", "")
		item.Disable()
		modelItems = append(modelItems, item)
		return
	}

	currentModel := os.Getenv("GEMINI_MODEL")
	if currentModel == "" {
		currentModel = "gemini-3.8-flash"
	}

	for _, mName := range models {
		item := mModel.AddSubMenuItem(mName, mName)
		if mName == currentModel {
			item.Check()
		}
		modelItems = append(modelItems, item)

		go func(it *systray.MenuItem, name string) {
			for range it.ClickedCh {
				config.Update("GEMINI_MODEL", name)
				for _, mi := range modelItems {
					mi.Uncheck()
				}
				it.Check()
			}
		}(item, mName)
	}
}

func handleEvents() {
	for {
		select {
		case <-mAPIKey.ClickedCh:
			cmd := exec.Command("osascript", "-e", `display dialog "Enter Gemini API Key:" default answer "" with title "API Key"`, "-e", `text returned of result`)
			out, err := cmd.Output()
			if err == nil {
				newKey := strings.TrimSpace(string(out))
				if newKey != "" {
					config.Update("GEMINI_API_KEY", newKey)
					systray.SetTitle("-")
					go refreshModelsMenu(newKey)
				}
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}
