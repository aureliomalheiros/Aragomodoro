package pomodoro

import (
	"testing"
)

func TestValidateDurations(t *testing.T) {
	tests := []struct {
		focus       	int
		breakTime   	int
		repeatTime  	int
		expectError 	bool
	}{
		{25, 5, 2, false},
		{0, 5, 1, true},
		{-1, 5, 0,  true},
		{25, 0, -1, true},
		{25, -5, 10, true},
		{61, 5, 1, true},
		{25, 61, 10, true},
	}

	for _, tt := range tests {
		err := ValidateDurations(tt.focus, tt.breakTime, tt.repeatTime)
		if tt.expectError && err == nil {
			t.Errorf("Expected error for input focus=%d, break=%d, repeat=%d", tt.focus, tt.breakTime, tt.repeatTime)
		}
		if !tt.expectError && err != nil {
			t.Errorf("Unexpected error for input focus=%d, break=%d, repeat=%d: %v", tt.focus, tt.breakTime, tt.repeatTime, err)
		}
	}
	
}

