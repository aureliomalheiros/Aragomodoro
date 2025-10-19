package pomodoro

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/aureliomalheiros/aragomodoro/internal/ascii_text"
	"github.com/aureliomalheiros/aragomodoro/internal/sound"
)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	ascii_text.PrintAsciiTextAragomodoro()
}

func ValidateDurations(focusDuration, breakDuration, repeatCount int) error {

	if focusDuration <= 0 || breakDuration <= 0 {
		return fmt.Errorf("❌ Invalid duration. Please provide positive integers for focus and break durations.")
	}
	if focusDuration > 60 || breakDuration > 60 {
		return fmt.Errorf("⚠️ Focus and break durations should not exceed 60 minutes.")
	}
	if repeatCount <= 0 {
		return fmt.Errorf("❌ Repeat count must be a positive integer.")
	}

	return nil
}

func PomodoroTimer(focusDuration int, breakDuration int, repeatCount int, continueOnBreak bool) {

	if err := ValidateDurations(focusDuration, breakDuration, repeatCount); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}

	if repeatCount > 1 {
		repeatPomodoro(focusDuration, breakDuration, repeatCount)
		return
	}

	fmt.Printf("🧭 Aragomodoro begins! Focus for %d minutes.\n", focusDuration)
	startTimer(time.Duration(focusDuration) * time.Minute)
	sound.ThemeAragorn()

	fmt.Printf("🌿 Time for a break! Rest for %d minutes.\n", breakDuration)
	startTimer(time.Duration(breakDuration) * time.Minute)
	sound.ThemeMountDoom()

	clearScreen()

	if continueOnBreak {
		continueOnBreak = false
		continuePomodoro(focusDuration, breakDuration)
	}

}

func startTimer(duration time.Duration) {
	for remaining := duration; remaining > 0; remaining -= time.Second {
		fmt.Printf("\r⏳ %v remaining", remaining.Truncate(time.Second))
		time.Sleep(time.Second)
	}
	fmt.Println("\r✅ Done!                        ")
}

func continuePomodoro(focusDuration int, breakDuration int) {
	for {
		PomodoroTimer(focusDuration, breakDuration, 1, true)
	}
}

func repeatPomodoro(focusDuration int, breakDuration int, repeatCount int) {
	for i := 1; i <= repeatCount; i++ {
		fmt.Printf("🔁 Starting Pomodoro session %d/%d...\n", i, repeatCount)
		PomodoroTimer(focusDuration, breakDuration, 1, false)
		if i < repeatCount {
			fmt.Println("🌟 Get ready for the next Pomodoro!")
		} else if i == repeatCount {
			fmt.Println("🎉 All Pomodoros completed! Great job!")
			fmt.Println("🍅 Time for a well-deserved long break!")
		}
	}
}

// func webModePomodoro(focusDuration int, breakDuration int, repeatCount int, continueOnBreak bool, webPort int) {
	
// }
