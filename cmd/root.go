package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/aureliomalheiros/aragomodoro/internal/pomodoro"
	"github.com/aureliomalheiros/aragomodoro/internal/ascii_text"

)

var (
	focusDuration int
	breakDuration int
	repeatCount   int
	continueOnBreak bool
)

var rootCmd = &cobra.Command{
	Use:   "aragomodoro",
	Short: "Aragomodoro: A playful Pomodoro timer inspired by Aragorn",
	Long:  "Aragomodoro is a playful take on the Pomodoro technique, inspired by the spirit of Aragorn from The Lord of the Rings.",
	Run: func(cmd *cobra.Command, args []string) {
		ascii_text.PrintAsciiText()
		pomodoro.PomodoroTimer(focusDuration, breakDuration, repeatCount, continueOnBreak)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&focusDuration, "focus", "f", 25, "Focus duration in minutes")
	rootCmd.Flags().IntVarP(&breakDuration, "break", "b", 5, "Break duration in minutes")
	rootCmd.Flags().IntVarP(&repeatCount, "repeat", "r", 1, "Number of Pomodoros before a long break")
	rootCmd.Flags().BoolVarP(&continueOnBreak, "continue", "c", false, "Continue the timer during breaks")
}




