package pomodoro

import (
	"testing"
)

func TestValidateDurations(t *testing.T) {
	tests := []struct {
		focus       int
		breakTime   int
		repeatTime  int
		expectError bool
	}{
		{25, 5, 2, false},
		{0, 5, 1, true},
		{-1, 5, 0, true},
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

func TestValidateDurations_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		focus       int
		breakTime   int
		repeatTime  int
		expectError bool
		description string
	}{
		{"ValidMinimum", 1, 1, 1, false, "Minimum valid values"},
		{"ValidMaximum", 60, 60, 100, false, "Maximum valid values"},
		{"ZeroFocus", 0, 5, 1, true, "Zero focus duration"},
		{"ZeroBreak", 25, 0, 1, true, "Zero break duration"},
		{"ZeroRepeat", 25, 5, 0, true, "Zero repeat count"},
		{"NegativeFocus", -1, 5, 1, true, "Negative focus duration"},
		{"NegativeBreak", 25, -1, 1, true, "Negative break duration"},
		{"NegativeRepeat", 25, 5, -1, true, "Negative repeat count"},
		{"ExcessiveFocus", 61, 5, 1, true, "Focus duration over 60 minutes"},
		{"ExcessiveBreak", 25, 61, 1, true, "Break duration over 60 minutes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDurations(tt.focus, tt.breakTime, tt.repeatTime)
			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error for focus=%d, break=%d, repeat=%d",
					tt.description, tt.focus, tt.breakTime, tt.repeatTime)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error for focus=%d, break=%d, repeat=%d: %v",
					tt.description, tt.focus, tt.breakTime, tt.repeatTime, err)
			}
		})
	}
}

func BenchmarkValidateDurations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateDurations(25, 5, 1)
	}
}
