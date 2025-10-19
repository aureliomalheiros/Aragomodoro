package sound

import (
	"testing"
	"time"
)

func TestSoftFocusComplete(t *testing.T) {
	originalMute := Mute
	Mute = true
	defer func() { Mute = originalMute }()

	SoftFocusComplete()
}

func TestSoftBreakComplete(t *testing.T) {
	originalMute := Mute
	Mute = true
	defer func() { Mute = originalMute }()

	SoftBreakComplete()
}

func TestSoftSoundsExecution(t *testing.T) {
	originalMute := Mute
	Mute = true
	defer func() { Mute = originalMute }()

	tests := []struct {
		name string
		fn   func()
	}{
		{"SoftFocusComplete", SoftFocusComplete},
		{"SoftBreakComplete", SoftBreakComplete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			tt.fn()
			duration := time.Since(start)

			if duration > 100*time.Millisecond {
				t.Errorf("Function %s took too long when muted: %v", tt.name, duration)
			}
		})
	}
}

func BenchmarkSoftFocusComplete(b *testing.B) {
	originalMute := Mute
	Mute = true
	defer func() { Mute = originalMute }()

	for i := 0; i < b.N; i++ {
		SoftFocusComplete()
	}
}

func BenchmarkSoftBreakComplete(b *testing.B) {
	originalMute := Mute
	Mute = true
	defer func() { Mute = originalMute }()

	for i := 0; i < b.N; i++ {
		SoftBreakComplete()
	}
}
