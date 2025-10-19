package ascii_text

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestPrintAsciiTextAragomodoro(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintAsciiTextAragomodoro()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if output == "" {
		t.Error("PrintAsciiTextAragomodoro should print something")
	}

	// Check for expected content
	if len(output) < 100 {
		t.Error("ASCII art output seems too short")
	}
}

func TestPrintAsciiTextBreak(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintAsciiTextBreak()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if output == "" {
		t.Error("PrintAsciiTextBreak should print something")
	}

	// Check for expected content
	if len(output) < 50 {
		t.Error("ASCII break text output seems too short")
	}
}

func TestPrintFunctions_NoErrors(t *testing.T) {
	// Test that functions don't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ASCII text functions panicked: %v", r)
		}
	}()

	PrintAsciiTextAragomodoro()
	PrintAsciiTextBreak()
}

func BenchmarkPrintAsciiTextAragomodoro(b *testing.B) {
	// Redirect output to avoid cluttering test output
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PrintAsciiTextAragomodoro()
	}
}

func BenchmarkPrintAsciiTextBreak(b *testing.B) {
	// Redirect output to avoid cluttering test output
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PrintAsciiTextBreak()
	}
}
