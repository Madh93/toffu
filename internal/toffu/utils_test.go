package toffu

import (
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

func TestSecondsToHumanReadable(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "Zero duration",
			duration: 0 * time.Second,
			expected: "0h 0m 0s",
		},
		{
			name:     "Seconds only",
			duration: 45 * time.Second,
			expected: "0h 0m 45s",
		},
		{
			name:     "Minutes and seconds",
			duration: 150 * time.Second,
			expected: "0h 2m 30s",
		},
		{
			name:     "Hours, minutes, and seconds",
			duration: 3665 * time.Second,
			expected: "1h 1m 5s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := secondsToHumanReadable(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}
