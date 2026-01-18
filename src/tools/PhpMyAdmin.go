package tools

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func StartPHPMyAdmin() {
	// Check if phpMyAdmin is already running
	if Docker.IsContainerRunning("phpmyadmin", true) {
		fmt.Println("phpMyAdmin is already running.")
		if Docker.AskYesNo("Do you want to remove the existing phpMyAdmin container and create a new one?") {
			fmt.Println("Removing existing phpMyAdmin container...")
			cmd := exec.Command("docker", "rm", "-f", "phpmyadmin")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error removing phpMyAdmin container:", err)
				return
			}
			fmt.Println("Existing phpMyAdmin container removed successfully.")
		} else {
			fmt.Println("Keeping existing phpMyAdmin container. Setup aborted.")
			return
		}
	}

	sqlContainers := Docker.ListOfContainers([]string{"mysql", "mariadb"})
	if len(sqlContainers) == 0 {
		fmt.Println("No running MySQL/MariaDB containers found.")
		return
	}

	items := append(sqlContainers, "Exit")
	prompt := promptui.Select{
		Label: "Select a SQL container to link with phpMyAdmin",
		Items: items,
	}
	_, selectedContainer, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		Cleanup()
	}
	if selectedContainer == "Exit" {
		fmt.Println("Exiting phpMyAdmin setup.")
		return
	}

	port := AskForInput("Enter host port to expose phpMyAdmin", "8080")

	fmt.Printf("Pulling phpMyAdmin image...\n")
	cmd := exec.Command("docker", "pull", "phpmyadmin/phpmyadmin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	// Build docker run command as args array
	args := []string{
		"run", "-d",
		"--restart", "unless-stopped",
		"--network", "ContainDB-Network",
		"--name", "phpmyadmin",
		"-e", fmt.Sprintf("PMA_HOST=%s", selectedContainer),
		"-p", fmt.Sprintf("%s:80", port),
		"phpmyadmin/phpmyadmin",
	}

	fmt.Println("Running: docker", strings.Join(args, " "))
	cmd = exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting phpMyAdmin:", err)
	} else {
		fmt.Printf("phpMyAdmin started. Access it at http://localhost:%s\n", port)
	}
}
