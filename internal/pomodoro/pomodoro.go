package pomodoro

import (
	"fmt"
	"os"
	"time"

	"github.com/aureliomalheiros/aragomodoro/internal/sound"
)

func ValidateDurations(focusDuration, breakDuration int) error {
	if focusDuration <= 0 || breakDuration <= 0 {
		return fmt.Errorf("❌ Invalid duration. Please provide positive integers for focus and break durations.")
	}
	if focusDuration > 60 || breakDuration > 60 {
		return fmt.Errorf("⚠️ Focus and break durations should not exceed 60 minutes.")
	}
	return nil
}

func PomodoroTimer(focusDuration int, breakDuration int) {
	
	if err := ValidateDurations(focusDuration, breakDuration); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}

	fmt.Printf("🧭 Aragomodoro begins! Focus for %d minutes.\n", focusDuration)
	startTimer(time.Duration(focusDuration) * time.Minute)
	sound.ThemeAragorn()

	fmt.Printf("🌿 Time for a break! Rest for %d minutes.\n", breakDuration)
	startTimer(time.Duration(breakDuration) * time.Minute)
	sound.ThemeMountDoom()

	fmt.Println("🍅 Session complete. Ready for the next adventure?")
}

func startTimer(duration time.Duration) {
	for remaining := duration; remaining > 0; remaining -= time.Second {
		fmt.Printf("\r⏳ %v remaining", remaining.Truncate(time.Second))
		time.Sleep(time.Second)
	}
	fmt.Println("\r✅ Done!                        ")
}

