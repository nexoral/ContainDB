package base

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func DockerStarter() {
	if !Docker.IsDockerInstalled() {
		fmt.Println("❌ Docker is not installed. Without Docker the tool cannot run.")
		installPrompt := promptui.Select{
			Label: "Would you like to install Docker now?",
			Items: []string{"Yes", "No", "Exit"},
		}
		_, choice, err := installPrompt.Run()
		if err != nil || choice != "Yes" {
			fmt.Println("Exiting. Please install Docker manually and rerun.")
			os.Exit(1)
		}
		err = Docker.InstallDocker()
		if err != nil {
			fmt.Println("Failed to install Docker:", err)
			os.Exit(1)
		}
		fmt.Println("Docker installed successfully! Please restart the terminal or log out & log in again.")
		os.Exit(0) // Exit to ensure user restarts terminal/logs out (or at least stops here)

	}
}
