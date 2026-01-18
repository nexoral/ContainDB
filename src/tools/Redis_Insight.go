package tools

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func StartRedisInsight() {
	// Check if RedisInsight is already running
	if Docker.IsContainerRunning("redisinsight", true) {
		fmt.Println("RedisInsight is already running.")
		if Docker.AskYesNo("Do you want to remove the existing RedisInsight container and create a new one?") {
			fmt.Println("Removing existing RedisInsight container...")
			cmd := exec.Command("docker", "rm", "-f", "redisinsight")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error removing RedisInsight container:", err)
				return
			}
			fmt.Println("Existing RedisInsight container removed successfully.")
		} else {
			fmt.Println("Keeping existing RedisInsight container. Setup aborted.")
			return
		}
	}

	// Look for running Redis containers
	redisContainers := Docker.ListOfContainers([]string{"redis"})
	if len(redisContainers) == 0 {
		fmt.Println("No running Redis containers found.")
		return
	}

	items := append(redisContainers, "Exit")
	prompt := promptui.Select{
		Label: "Select a Redis container to link with RedisInsight",
		Items: items,
	}
	_, selectedContainer, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		Cleanup()
	}
	if selectedContainer == "Exit" {
		fmt.Println("Exiting RedisInsight setup.")
		return
	}

	port := AskForInput("Enter host port to expose RedisInsight", "8001")

	fmt.Printf("Pulling RedisInsight image...\n")
	cmd := exec.Command("docker", "pull", "redis/redisinsight:latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	// Build docker run command as args array
	args := []string{
		"run", "-d",
		"--restart", "unless-stopped",
		"--network", "ContainDB-Network",
		"--name", "redisinsight",
		"-p", fmt.Sprintf("%s:5540", port),
		"redis/redisinsight:latest",
	}

	fmt.Println("Running: docker", strings.Join(args, " "))
	cmd = exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting RedisInsight:", err)
	} else {
		fmt.Printf("\n✅ RedisInsight started. Access it at: http://localhost:%s\n", port)
		fmt.Printf("👉 In the RedisInsight UI, add a Redis database with host: `%s`, port: `6379`\n", selectedContainer)
		fmt.Println("   (RedisInsight will resolve container name using Docker network DNS.)")
	}
}
