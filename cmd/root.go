package cmd

import (
	"fmt"
	"os"

	"github.com/aureliomalheiros/aragomodoro/internal/ascii_text"
	"github.com/aureliomalheiros/aragomodoro/internal/pomodoro"
	"github.com/aureliomalheiros/aragomodoro/internal/web"
	"github.com/spf13/cobra"
)

var (
	focusDuration   int
	breakDuration   int
	repeatCount     int
	continueOnBreak bool
	webMode         bool
	webPort         int
)

var rootCmd = &cobra.Command{
	Use:   "aragomodoro",
	Short: "Aragomodoro: A playful Pomodoro timer inspired by Aragorn",
	Long:  "Aragomodoro is a playful take on the Pomodoro technique, inspired by the spirit of Aragorn from The Lord of the Rings.",
	Run: func(cmd *cobra.Command, args []string) {
		if webMode {
			var port int = webPort
			if port == 0 {
				port = 8080
			}
			fmt.Printf("Starting Aragomodoro web interface on port %d...\n", port)
			fmt.Printf("Access at: http://localhost:%d\n", port)
			
			webServer := web.NewServer(port)
			if err := webServer.Start(); err != nil {
				panic(err)
			}
		} else {
			ascii_text.PrintAsciiTextAragomodoro()
			pomodoro.PomodoroTimer(focusDuration, breakDuration, repeatCount, continueOnBreak)
		}
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
	rootCmd.Flags().BoolVarP(&webMode, "web", "w", false, "Start the web interface")
	rootCmd.Flags().IntVarP(&webPort, "port", "p", 8080, "Port for the web server")
}
