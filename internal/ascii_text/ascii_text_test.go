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

	if buf.String() == "" {
		t.Error("PrintAsciiTextAragomodoro should print something")
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

	if buf.String() == "" {
		t.Error("PrintAsciiTextBreak should print something")
	}
}
