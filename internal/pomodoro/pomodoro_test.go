package pomodoro

import (
	"testing"
)

func TestValidateDurations(t *testing.T) {
	tests := []struct {
		focus       int
		breakTime   int
		expectError bool
	}{
		{25, 5, false},
		{0, 5, true},
		{-1, 5, true},
		{25, 0, true},
		{25, -5, true},
		{61, 5, true},
		{25, 61, true},
	}

	for _, tt := range tests {
		err := ValidateDurations(tt.focus, tt.breakTime)
		if tt.expectError && err == nil {
			t.Errorf("Expected error for input focus=%d, break=%d", tt.focus, tt.breakTime)
		}
		if !tt.expectError && err != nil {
			t.Errorf("Unexpected error for input focus=%d, break=%d: %v", tt.focus, tt.breakTime, err)
		}
	}
}

