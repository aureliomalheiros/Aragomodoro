package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/aureliomalheiros/aragomodoro/internal/pomodoro"
)

var (
	focusDuration int
	breakDuration int
)

var rootCmd = &cobra.Command{
	Use:   "aragomodoro",
	Short: "Aragomodoro: A playful Pomodoro timer inspired by Aragorn",
	Long:  "Aragomodoro is a playful take on the Pomodoro technique, inspired by the spirit of Aragorn from The Lord of the Rings.",
	Run: func(cmd *cobra.Command, args []string) {
		pomodoro.PomodoroTimer(focusDuration, breakDuration)
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
}




