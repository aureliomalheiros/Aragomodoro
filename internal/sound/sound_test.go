package sound

import (
	"log"
	"os"
	"testing"
)

func TestPlaySound_InvalidFile(t *testing.T) {
	null, _ := os.Open(os.DevNull)
	log.SetOutput(null)
	defer log.SetOutput(os.Stderr)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlaySound should not panic. Got: %v", r)
		}
	}()

	PlaySound("nonexistent.wav")
}
