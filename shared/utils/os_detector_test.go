package utils

import (
	"testing"
)

func TestDetectOSFromUserAgent(t *testing.T) {
	testCases := []struct {
		name      string
		userAgent string
		expected  OSInfo
	}{
		{
			name:      "iPhone Safari",
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
			expected:  OSInfo{OS: "iOS", Version: "17.1.1", Device: "iPhone"},
		},
		{
			name:      "iPad Safari",
			userAgent: "Mozilla/5.0 (iPad; CPU OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
			expected:  OSInfo{OS: "iOS", Version: "17.1", Device: "iPad"},
		},
		{
			name:      "Android Chrome",
			userAgent: "Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Mobile Safari/537.36",
			expected:  OSInfo{OS: "Android", Version: "14"},
		},
		{
			name:      "Windows Chrome",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			expected:  OSInfo{OS: "Windows", Version: "10"},
		},
		{
			name:      "macOS Safari",
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			expected:  OSInfo{OS: "macOS", Version: "10.15.7"},
		},
		{
			name:      "Linux Firefox",
			userAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0",
			expected:  OSInfo{OS: "Linux"},
		},
		{
			name:      "Ubuntu Chrome",
			userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Ubuntu",
			expected:  OSInfo{OS: "Ubuntu"},
		},
		{
			name:      "Empty User-Agent",
			userAgent: "",
			expected:  OSInfo{OS: "Unknown"},
		},
		{
			name:      "Unknown User-Agent",
			userAgent: "SomeCustomBot/1.0",
			expected:  OSInfo{OS: "Unknown"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DetectOSFromUserAgent(tc.userAgent)
			
			if result.OS != tc.expected.OS {
				t.Errorf("Expected OS %s, got %s", tc.expected.OS, result.OS)
			}
			
			if tc.expected.Version != "" && result.Version != tc.expected.Version {
				t.Errorf("Expected version %s, got %s", tc.expected.Version, result.Version)
			}
			
			if tc.expected.Device != "" && result.Device != tc.expected.Device {
				t.Errorf("Expected device %s, got %s", tc.expected.Device, result.Device)
			}
		})
	}
}

func TestGetSimpleOS(t *testing.T) {
	testCases := []struct {
		userAgent string
		expected  string
	}{
		{"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X)", "iOS"},
		{"Mozilla/5.0 (Linux; Android 14; SM-G998B)", "Android"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)", "Windows"},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)", "macOS"},
		{"Mozilla/5.0 (X11; Linux x86_64)", "Linux"},
		{"", "Unknown"},
	}

	for _, tc := range testCases {
		result := GetSimpleOS(tc.userAgent)
		if result != tc.expected {
			t.Errorf("For user agent %s, expected %s, got %s", tc.userAgent, tc.expected, result)
		}
	}
}
