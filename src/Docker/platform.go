package Docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// IsAdmin checks if the program is running with administrator/root privileges
// Cross-platform implementation
func IsAdmin() bool {
	if runtime.GOOS == "windows" {
		// On Windows, check if running with admin privileges
		cmd := exec.Command("net", "session")
		err := cmd.Run()
		return err == nil
	} else {
		// On Unix-like systems, check if UID is 0
		return os.Geteuid() == 0
	}
}

// GetTempDir returns the appropriate temporary directory for the current platform
func GetTempDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("TEMP")
	}
	return "/tmp"
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsMacOS returns true if running on macOS
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux returns true if running on Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// GetOSName returns a friendly OS name string
func GetOSName() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	default:
		return runtime.GOOS
	}
}

// CheckOSSupport checks if the current OS is supported
func CheckOSSupport() error {
	osName := GetOSName()
	fmt.Printf("Detected OS: %s\n", osName)
	
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		return nil
	default:
		return fmt.Errorf("unsupported operating system: %s", osName)
	}
}

// GetOSRelease reads OS release information (Linux-specific, returns empty on other OS)
func GetOSRelease() string {
	if runtime.GOOS != "linux" {
		return ""
	}
	
	releaseBytes, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}
	return string(releaseBytes)
}

// CheckDockerCommand checks if a Docker command exists and is executable
func CheckDockerCommand(cmd string) error {
	dockerCmd := exec.Command(cmd, "--version")
	_, err := dockerCmd.Output()
	if err != nil {
		return fmt.Errorf("%s is not installed or not accessible", cmd)
	}
	return nil
}

// GetShell returns the appropriate shell for the platform
func GetShell() string {
	if runtime.GOOS == "windows" {
		// Check for PowerShell first, then cmd
		if _, err := exec.LookPath("powershell"); err == nil {
			return "powershell"
		}
		return "cmd"
	}
	// For Unix-like systems, prefer bash but fall back to sh
	if _, err := exec.LookPath("bash"); err == nil {
		return "bash"
	}
	return "sh"
}

// ExecuteCommand runs a command in a platform-appropriate way
func ExecuteCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// BuildDockerRunCommand constructs docker run arguments as a slice for exec.Command
func BuildDockerRunCommand(args []string) []string {
	dockerArgs := []string{"run", "-d"}
	dockerArgs = append(dockerArgs, args...)
	return dockerArgs
}
