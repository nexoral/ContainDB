package base

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

const Version = "6.16.41-stable"

// GetInstallationMethod detects how ContainDB was installed
func GetInstallationMethod() string {
	// Check if installed via dpkg (Debian package)
	cmd := exec.Command("dpkg", "-s", "containdb")
	if err := cmd.Run(); err == nil {
		return "Debian Package (.deb)"
	}

	// Check if installed via npm
	cmd = exec.Command("npm", "list", "-g", "containdb")
	output, err := cmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), "containdb") {
		return "NPM Package"
	}

	// Default to manual installation
	return "Manual Installation"
}

func ShowBanner() {
	// Define styles
	boldGreen := color.New(color.FgGreen, color.Bold).SprintFunc()
	boldCyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	boldWhite := color.New(color.FgWhite, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	// green := color.New(color.FgGreen).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	clear := "\033[2J\033[H" // Clear terminal
	fmt.Print(clear)

	// Banner
	termWidth := 80 // Default width, could be made dynamic with terminal size detection

	welcomeMsg := "An Awesome Project to Manage Your Databases in Containers"
	padding := (termWidth - len(welcomeMsg)) / 2

	banner := "\n" + strings.Repeat(" ", padding) +
		strings.ToUpper(welcomeMsg) +
		strings.Repeat(" ", padding) + "\n"

	// Create a box around the message
	horizontalLine := strings.Repeat("═", termWidth)
	banner = "\n" + horizontalLine + "\n" + banner + horizontalLine + "\n"

	fmt.Println(boldGreen(banner))

	border := boldCyan(strings.Repeat("─", 80))
	fmt.Println(border)

	fmt.Printf("%s\n", boldWhite("🛠️  Welcome to ")+boldGreen("ContainDB")+boldWhite(" - Containerized Database Manager CLI"))
	fmt.Println(border)

	// Get installation method
	installMethod := GetInstallationMethod()

	// Info Block
	fmt.Printf("%s %s\n", boldCyan("📦 Version:"), white(Version))
	fmt.Printf("%s %s\n", boldCyan("📥 Installed via:"), white(installMethod))
	fmt.Printf("%s %s\n", boldCyan("👨‍💻 Author:"), white("Ankan Saha"))
	fmt.Printf("%s %s\n", boldCyan("🔗 GitHub:"), cyan("https://github.com/nexoral/ContainDB"))
	fmt.Printf("%s %s\n", boldCyan("💖 Sponsor:"), cyan("https://github.com/sponsors/AnkanSaha"))
	fmt.Printf("%s %s\n", boldCyan("📄 Docs:"), cyan("https://github.com/nexoral/ContainDB/wiki"))
	fmt.Printf("%s %s\n", boldCyan("🔐 License:"), white("MIT License"))
	fmt.Printf("%s %s\n", boldCyan("💬 Feedback:"), white("Feel free to open issues or discussions on GitHub"))

	fmt.Println(border)
	fmt.Printf("%s\n", boldCyan("⚡ Tip: ")+dim("Run `containDB --help` to see available commands."))
	fmt.Println(border)
	fmt.Println()
}
