package ratelimit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"aiswer/internal/config"
)

type Usage struct {
	LastMinute        time.Time `json:"last_minute"`
	RequestsPerMinute int       `json:"requests_per_minute"`
	TokensPerMinute   int       `json:"tokens_per_minute"`

	LastDay        time.Time `json:"last_day"`
	RequestsPerDay int       `json:"requests_per_day"`
}

var (
	usageFile string
	usage     Usage
	mu        sync.Mutex
)

func init() {
	usageFile = filepath.Join(config.GetConfigDir(), "ratelimit.json")
	loadUsage()
}

func loadUsage() {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(usageFile)
	if err == nil {
		_ = json.Unmarshal(data, &usage)
	}

	now := time.Now()
	// Reset minute counters if it's a different minute
	if now.Truncate(time.Minute) != usage.LastMinute.Truncate(time.Minute) {
		usage.RequestsPerMinute = 0
		usage.TokensPerMinute = 0
		usage.LastMinute = now
	}

	// Reset daily counters if it's a different day
	y1, m1, d1 := now.Date()
	y2, m2, d2 := usage.LastDay.Date()
	if y1 != y2 || m1 != m2 || d1 != d2 {
		usage.RequestsPerDay = 0
		usage.LastDay = now
	}
}

func saveUsage() {
	data, err := json.MarshalIndent(usage, "", "  ")
	if err == nil {
		_ = os.WriteFile(usageFile, data, 0644)
	}
}

// CheckLimit returns an error if the model limits are exceeded.
func CheckLimit(modelName string) error {
	tier := os.Getenv("GEMINI_TIER")
	if tier != "free" {
		return nil // No limit for paid tier
	}

	mu.Lock()
	defer mu.Unlock()

	// Update time windows
	now := time.Now()
	if now.Truncate(time.Minute) != usage.LastMinute.Truncate(time.Minute) {
		usage.RequestsPerMinute = 0
		usage.TokensPerMinute = 0
		usage.LastMinute = now
	}
	y1, m1, d1 := now.Date()
	y2, m2, d2 := usage.LastDay.Date()
	if y1 != y2 || m1 != m2 || d1 != d2 {
		usage.RequestsPerDay = 0
		usage.LastDay = now
	}

	lowerName := strings.ToLower(modelName)
	if strings.Contains(lowerName, "-pro") {
		return fmt.Errorf("Pro models are not available on the Free tier")
	}

	// Default to flash limits
	maxRPM, maxTPM, maxRPD := 5, 250000, 20
	if strings.Contains(lowerName, "flash-lite") {
		maxRPM, maxTPM, maxRPD = 15, 250000, 500
	}

	if usage.RequestsPerDay >= maxRPD {
		return fmt.Errorf("Daily request limit reached for Free tier (%d/%d RPD)", usage.RequestsPerDay, maxRPD)
	}
	if usage.RequestsPerMinute >= maxRPM {
		return fmt.Errorf("Minute request limit reached for Free tier (%d/%d RPM)", usage.RequestsPerMinute, maxRPM)
	}
	if usage.TokensPerMinute >= maxTPM {
		return fmt.Errorf("Minute token limit reached for Free tier (%d/%d TPM)", usage.TokensPerMinute, maxTPM)
	}

	return nil
}

// AddUsage records a successful API call's usage.
func AddUsage(modelName string, tokens int) {
	tier := os.Getenv("GEMINI_TIER")
	if tier != "free" {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	usage.RequestsPerMinute++
	usage.TokensPerMinute += tokens
	usage.RequestsPerDay++

	saveUsage()
}
