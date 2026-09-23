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
	usageMap  map[string]*Usage
	mu        sync.Mutex
)

func init() {
	usageFile = filepath.Join(config.GetConfigDir(), "ratelimit.json")
	usageMap = make(map[string]*Usage)
	loadUsage()
}

func getUsageForModel(modelName string) *Usage {
	if u, exists := usageMap[modelName]; exists {
		return u
	}
	u := &Usage{}
	usageMap[modelName] = u
	return u
}

func updateTimeWindows(u *Usage) {
	now := time.Now()
	// Reset minute counters if it's a different minute
	if now.Truncate(time.Minute) != u.LastMinute.Truncate(time.Minute) {
		u.RequestsPerMinute = 0
		u.TokensPerMinute = 0
		u.LastMinute = now
	}

	// Reset daily counters if it's a different day
	y1, m1, d1 := now.Date()
	y2, m2, d2 := u.LastDay.Date()
	if y1 != y2 || m1 != m2 || d1 != d2 {
		u.RequestsPerDay = 0
		u.LastDay = now
	}
}

func loadUsage() {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(usageFile)
	if err == nil {
		_ = json.Unmarshal(data, &usageMap)
	}

	for _, u := range usageMap {
		updateTimeWindows(u)
	}
}

func saveUsage() {
	data, err := json.MarshalIndent(usageMap, "", "  ")
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

	u := getUsageForModel(modelName)
	updateTimeWindows(u)

	lowerName := strings.ToLower(modelName)
	if strings.Contains(lowerName, "-pro") {
		return fmt.Errorf("Pro models are not available on the Free tier")
	}

	// Default to flash limits
	maxRPM, maxTPM, maxRPD := 5, 250000, 20
	if strings.Contains(lowerName, "flash-lite") {
		maxRPM, maxTPM, maxRPD = 15, 250000, 500
	}

	if u.RequestsPerDay >= maxRPD {
		return fmt.Errorf("Daily request limit reached for %s on Free tier (%d/%d RPD)", modelName, u.RequestsPerDay, maxRPD)
	}
	if u.RequestsPerMinute >= maxRPM {
		return fmt.Errorf("Minute request limit reached for %s on Free tier (%d/%d RPM)", modelName, u.RequestsPerMinute, maxRPM)
	}
	if u.TokensPerMinute >= maxTPM {
		return fmt.Errorf("Minute token limit reached for %s on Free tier (%d/%d TPM)", modelName, u.TokensPerMinute, maxTPM)
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

	u := getUsageForModel(modelName)
	updateTimeWindows(u)

	u.RequestsPerMinute++
	u.TokensPerMinute += tokens
	u.RequestsPerDay++

	saveUsage()
}
