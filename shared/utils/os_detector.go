package utils

import (
	"regexp"
	"strings"
)

// OSInfo contains detected operating system information
type OSInfo struct {
	OS      string `json:"os"`
	Version string `json:"version,omitempty"`
	Device  string `json:"device,omitempty"`
}

// DetectOSFromUserAgent detects OS from User-Agent header
func DetectOSFromUserAgent(userAgent string) OSInfo {
	if userAgent == "" {
		return OSInfo{OS: "Unknown"}
	}

	userAgent = strings.ToLower(userAgent)

	// Mobile Operating Systems (check first as they're more specific)
	if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") {
		version := extractIOSVersion(userAgent)
		device := "iPhone"
		if strings.Contains(userAgent, "ipad") {
			device = "iPad"
		}
		return OSInfo{OS: "iOS", Version: version, Device: device}
	}

	if strings.Contains(userAgent, "android") {
		version := extractAndroidVersion(userAgent)
		return OSInfo{OS: "Android", Version: version}
	}

	// Desktop Operating Systems
	if strings.Contains(userAgent, "windows nt") {
		version := extractWindowsVersion(userAgent)
		return OSInfo{OS: "Windows", Version: version}
	}

	if strings.Contains(userAgent, "mac os x") || strings.Contains(userAgent, "macos") {
		version := extractMacOSVersion(userAgent)
		return OSInfo{OS: "macOS", Version: version}
	}

	if strings.Contains(userAgent, "linux") {
		// Check for specific Linux distributions
		if strings.Contains(userAgent, "ubuntu") {
			return OSInfo{OS: "Ubuntu"}
		}
		if strings.Contains(userAgent, "fedora") {
			return OSInfo{OS: "Fedora"}
		}
		if strings.Contains(userAgent, "debian") {
			return OSInfo{OS: "Debian"}
		}
		return OSInfo{OS: "Linux"}
	}

	// Other systems
	if strings.Contains(userAgent, "freebsd") {
		return OSInfo{OS: "FreeBSD"}
	}

	if strings.Contains(userAgent, "openbsd") {
		return OSInfo{OS: "OpenBSD"}
	}

	return OSInfo{OS: "Unknown"}
}

// extractIOSVersion extracts iOS version from User-Agent
func extractIOSVersion(userAgent string) string {
	re := regexp.MustCompile(`os (\d+)_(\d+)(?:_(\d+))?`)
	matches := re.FindStringSubmatch(userAgent)
	if len(matches) >= 3 {
		version := matches[1] + "." + matches[2]
		if len(matches) > 3 && matches[3] != "" {
			version += "." + matches[3]
		}
		return version
	}
	return ""
}

// extractAndroidVersion extracts Android version from User-Agent
func extractAndroidVersion(userAgent string) string {
	re := regexp.MustCompile(`android (\d+(?:\.\d+)*)`)
	matches := re.FindStringSubmatch(userAgent)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractWindowsVersion extracts Windows version from User-Agent
func extractWindowsVersion(userAgent string) string {
	versionMap := map[string]string{
		"windows nt 10.0": "10",
		"windows nt 6.3":  "8.1",
		"windows nt 6.2":  "8",
		"windows nt 6.1":  "7",
		"windows nt 6.0":  "Vista",
		"windows nt 5.1":  "XP",
	}

	for pattern, version := range versionMap {
		if strings.Contains(userAgent, pattern) {
			return version
		}
	}
	return ""
}

// extractMacOSVersion extracts macOS version from User-Agent
func extractMacOSVersion(userAgent string) string {
	re := regexp.MustCompile(`mac os x (\d+)_(\d+)(?:_(\d+))?`)
	matches := re.FindStringSubmatch(userAgent)
	if len(matches) >= 3 {
		version := matches[1] + "." + matches[2]
		if len(matches) > 3 && matches[3] != "" {
			version += "." + matches[3]
		}
		return version
	}
	return ""
}

// GetSimpleOS returns just the OS name without version details
func GetSimpleOS(userAgent string) string {
	osInfo := DetectOSFromUserAgent(userAgent)
	return osInfo.OS
}
