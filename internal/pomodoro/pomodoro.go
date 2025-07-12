package pomodoro

import (
	"fmt"
	"os"
	"time"
)

func PomodoroTimer(focusDuration int, breakDuration int) {
	if focusDuration <= 0 || breakDuration <= 0 {
		fmt.Println("❌ Invalid duration. Please provide positive integers for focus and break durations.")
		os.Exit(1)
	}
	if focusDuration > 60 || breakDuration > 60 {
		fmt.Println("⚠️ Focus and break durations should not exceed 60 minutes.")
		os.Exit(1)
	}

	fmt.Printf("🧭 Aragomodoro begins! Focus for %d minutes.\n", focusDuration)
	startTimer(time.Duration(focusDuration) * time.Minute)
	sound.PlaySound()

	fmt.Printf("🌿 Time for a break! Rest for %d minutes.\n", breakDuration)
	startTimer(time.Duration(breakDuration) * time.Minute)
	sound.PlaySound()

	fmt.Println("🍅 Session complete. Ready for the next adventure?")
}

func startTimer(duration time.Duration) {
	for remaining := duration; remaining > 0; remaining -= time.Second {
		fmt.Printf("\r⏳ %v remaining", remaining.Truncate(time.Second))
		time.Sleep(time.Second)
	}
	fmt.Println("\r✅ Done!                        ")
}

