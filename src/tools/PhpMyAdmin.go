package tools

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

// CloudDBConfig holds cloud database connection parameters
type CloudDBConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	Database  string // Optional
	EnableSSL bool
}

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

	// Detect local MySQL/MariaDB containers
	sqlContainers := Docker.ListOfContainers([]string{"mysql", "mariadb"})

	// Present connection type selection
	connectionType := selectConnectionType(len(sqlContainers) > 0)
	if connectionType == "exit" {
		fmt.Println("Exiting phpMyAdmin setup.")
		return
	}

	// Branch to appropriate handler
	if connectionType == "local" {
		startPHPMyAdminLocal(sqlContainers)
	} else {
		startPHPMyAdminCloud()
	}
}

// selectConnectionType prompts user to choose between local container or cloud database
func selectConnectionType(hasLocalContainers bool) string {
	items := []string{}

	if hasLocalContainers {
		items = append(items, "Local Container")
	}
	items = append(items, "Cloud Database", "Exit")

	prompt := promptui.Select{
		Label: "Select phpMyAdmin connection type",
		Items: items,
	}
	_, selected, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		Cleanup()
		return "exit"
	}

	if selected == "Exit" {
		return "exit"
	} else if selected == "Local Container" {
		return "local"
	} else {
		return "cloud"
	}
}

// startPHPMyAdminLocal handles local container connection (existing logic)
func startPHPMyAdminLocal(sqlContainers []string) {
	items := append(sqlContainers, "Exit")
	prompt := promptui.Select{
		Label: "Select a SQL container to link with phpMyAdmin",
		Items: items,
	}
	_, selectedContainer, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		Cleanup()
		return
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

// startPHPMyAdminCloud handles cloud database connection (new logic)
func startPHPMyAdminCloud() {
	config := getCloudConnectionConfig()
	port := AskForInput("Enter host port to expose phpMyAdmin", "8080")

	// Validate inputs
	if config.Host == "" {
		fmt.Println("Error: Database host cannot be empty")
		return
	}
	if config.Username == "" {
		fmt.Println("Error: Username cannot be empty")
		return
	}
	if config.Password == "" {
		fmt.Println("Error: Password cannot be empty")
		return
	}

	// Pull image
	fmt.Printf("Pulling phpMyAdmin image...\n")
	cmd := exec.Command("docker", "pull", "phpmyadmin/phpmyadmin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	// Build docker run args
	args := []string{
		"run", "-d",
		"--restart", "unless-stopped",
		"--network", "bridge",
		"--name", "phpmyadmin",
		"-e", "PMA_ARBITRARY=1",
		"-e", fmt.Sprintf("PMA_HOST=%s", config.Host),
		"-e", fmt.Sprintf("PMA_PORT=%s", config.Port),
		"-e", fmt.Sprintf("PMA_USER=%s", config.Username),
		"-e", fmt.Sprintf("PMA_PASSWORD=%s", config.Password),
	}

	if config.Database != "" {
		args = append(args, "-e", fmt.Sprintf("PMA_DATABASE=%s", config.Database))
	}

	if config.EnableSSL {
		args = append(args, "-e", "PMA_SSL=1")
	} else {
		args = append(args, "-e", "PMA_SSL=0")
	}

	// Always disable SSL verification by default
	args = append(args, "-e", "PMA_SSL_VERIFY=0")

	args = append(args, "-p", fmt.Sprintf("%s:80", port))
	args = append(args, "phpmyadmin/phpmyadmin")

	fmt.Println("Running: docker", strings.Join(args, " "))
	cmd = exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting phpMyAdmin:", err)
	} else {
		fmt.Printf("\n✅ phpMyAdmin started! Access it at http://localhost:%s\n", port)
		fmt.Printf("📋 Cloud database connection:\n")
		fmt.Printf("   Host: %s:%s\n", config.Host, config.Port)
		fmt.Printf("   User: %s\n", config.Username)
		fmt.Printf("   SSL: %v\n", config.EnableSSL)
	}
}

// getCloudConnectionConfig collects cloud database connection parameters from user
func getCloudConnectionConfig() CloudDBConfig {
	fmt.Println("\n🌐 Cloud Database Configuration")

	host := AskForInput("Enter database host (e.g., learn-mysql-db.mysql.database.azure.com)", "")
	port := AskForInput("Enter database port", "3306")
	username := AskForInput("Enter database username", "")
	password := AskForInput("Enter database password", "")
	database := AskForInput("Enter database name (optional, leave empty for none)", "")

	enableSSL := Docker.AskYesNo("Enable SSL connection?")

	return CloudDBConfig{
		Host:      strings.TrimSpace(host),
		Port:      strings.TrimSpace(port),
		Username:  strings.TrimSpace(username),
		Password:  strings.TrimSpace(password),
		Database:  strings.TrimSpace(database),
		EnableSSL: enableSSL,
	}
}
