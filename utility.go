package main

import (
	"strings"
	"time"
)

// Determine if a label should be included based on the excluded labels list
func shouldInclude(label string) bool {
	for _, excludedLabel := range excludedLabels {
		if label == excludedLabel {
			return false
		}
	}
	return true
}

func formatTimestamp(ts string) string {
	parsedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts // If parsing fails, return the original timestamp
	}
	return parsedTime.Format(timeFormat)
}

func titleize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// Usage: {{ .StartsAt | formatTimestampWithOffset -10 }}
func formatTimestampWithOffset(t time.Time, offsetMinutes int) string {
	return t.Add(time.Duration(offsetMinutes) * time.Minute).Format(timeFormat)
}

// replaceAll replaces all occurrences of old with new in the input string
func replaceAll(string, before, after string) string {
	return strings.ReplaceAll(string, before, after)
}
