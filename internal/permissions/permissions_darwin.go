//go:build darwin

package permissions

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreFoundation
#include <ApplicationServices/ApplicationServices.h>
#include <CoreFoundation/CoreFoundation.h>

bool checkAndPromptAccessibility() {
    const void *keys[] = { kAXTrustedCheckOptionPrompt };
    const void *values[] = { kCFBooleanTrue };
    CFDictionaryRef options = CFDictionaryCreate(kCFAllocatorDefault,
                                                 keys, values, 1,
                                                 &kCFCopyStringDictionaryKeyCallBacks,
                                                 &kCFTypeDictionaryValueCallBacks);
    bool trusted = AXIsProcessTrustedWithOptions(options);
    CFRelease(options);
    return trusted;
}

bool checkAccessibility() {
    return AXIsProcessTrusted();
}
*/
import "C"

func CheckAndPromptAccessibility() bool {
	return bool(C.checkAndPromptAccessibility())
}

func CheckAccessibility() bool {
	return bool(C.checkAccessibility())
}
