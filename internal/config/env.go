package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Init loads environment variables from embedded config (if any), executable dir, and current dir.
func Init(embeddedEnv []byte) {
	if len(embeddedEnv) > 0 {
		if envMap, err := godotenv.UnmarshalBytes(embeddedEnv); err == nil {
			for k, v := range envMap {
				if os.Getenv(k) == "" {
					os.Setenv(k, v)
				}
			}
		}
	}

	if exePath, err := os.Executable(); err == nil {
		_ = godotenv.Overload(filepath.Join(filepath.Dir(exePath), ".env"))
	}
	_ = godotenv.Overload()
}

// Update writes the given key-value pair to the .env file and sets it in the current environment.
func Update(key, val string) {
	os.Setenv(key, val)

	envPath := ".env"
	if exePath, err := os.Executable(); err == nil {
		exeDirEnv := filepath.Join(filepath.Dir(exePath), ".env")
		if _, err := os.Stat(exeDirEnv); err == nil {
			envPath = exeDirEnv
		} else {
			envPath = exeDirEnv
		}
	}

	var content []byte
	var err error
	if content, err = os.ReadFile(envPath); err != nil {
		content = []byte{}
	}

	lines := strings.Split(string(content), "\n")
	if len(content) == 0 {
		lines = []string{}
	}

	found := false
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), key+"=") {
			lines[i] = fmt.Sprintf(`%s="%s"`, key, val)
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, fmt.Sprintf(`%s="%s"`, key, val))
	}

	_ = os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0644)
}
