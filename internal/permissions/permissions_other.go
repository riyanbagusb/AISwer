//go:build !darwin

package permissions

func CheckAndPromptAccessibility() bool {
	return true
}

func CheckAccessibility() bool {
	return true
}
