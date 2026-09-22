package notify

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

// Show displays a native OS notification with the provided answer
func Show(answer string) error {
	title := "Exam OCR Assistant"
	message := fmt.Sprintf("%s", answer)

	// beeep.Notify takes (title, message, appIcon)
	err := beeep.Notify(title, message, "")
	if err != nil {
		return fmt.Errorf("failed to show notification: %w", err)
	}
	return nil
}
