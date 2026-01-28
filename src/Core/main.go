package main

import (
	"ContainDB/src/Docker"
	"ContainDB/src/base"
	"ContainDB/src/tools"
	"fmt"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	VERSION := "5.14.40-stable"

	// handle version flag without requiring sudo
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("ContainDB CLI Version:", VERSION)
		return
	} else if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println("ContainDB CLI - A tool for managing Docker databases")
		if runtime.GOOS != "windows" {
			fmt.Println("Usage: sudo containdb")
		} else {
			fmt.Println("Usage: containdb (run as Administrator)")
		}
		fmt.Println("Options:")
		fmt.Println("  --version   Show version information")
		fmt.Println("  --help             Show this help message")
		fmt.Println("  --install-docker   Install Docker if not installed")
		fmt.Println("  --uninstall-docker Uninstall Docker if installed")
		fmt.Println("  --export   Export Docker Compose file with all running services")
		fmt.Println("  --import ./docker-compose.yml      Import and run services from a Docker Compose file")
		os.Exit(0) // Exit after handling flags
	} else if len(os.Args) > 1 && os.Args[1] == "--install-docker" {
		if !Docker.IsDockerInstalled() {
			fmt.Println("Docker is not installed. Installing Docker...")
			err := Docker.InstallDocker()
			if err != nil {
				fmt.Println("Failed to install Docker:", err)
			}
			fmt.Println("Docker installed successfully! Please restart the terminal or log out & log in again.")
		} else {
			fmt.Println("Docker is already installed.")
		}
		os.Exit(0) // Exit after handling flags
	}

	// Replace Ctrl+C handler to avoid triggering on normal exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		tools.Cleanup()
		os.Exit(1)
	}()

	// Check OS support
	if err := Docker.CheckOSSupport(); err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	// Check for admin/root privileges (only on Linux, optional on Windows/macOS for Docker)
	if runtime.GOOS == "linux" {
		if !Docker.IsAdmin() {
			fmt.Println("❌ Please run this program with sudo (Linux requires root for Docker)")
			os.Exit(1)
		}
	} else if runtime.GOOS == "windows" {
		if !Docker.IsAdmin() {
			fmt.Println("⚠️  Warning: Not running as Administrator. Docker may require admin privileges.")
			fmt.Println("   Continuing anyway...")
		}
	}

	// Check if Docker is installed and if not, prompt to install it
	base.DockerStarter()

	// Handle command line flags with FlagHandler function
	base.FlagHandler()

	errs := Docker.CreateDockerNetworkIfNotExists()
	if errs != nil {
		fmt.Println("Failed to create Docker network:", errs)
		return
	}

	// Show welcome banner
	base.ShowBanner()

	// Start the base case handler which contains the main menu
	base.BaseCaseHandler()
}
