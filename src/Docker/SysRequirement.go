package Docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func CheckSystemRequirements() {
	// Check if Docker is installed and accessible
	if err := checkDockerInstallation(); err != nil {
		fmt.Printf("WARNING: %v\n", err)
	}

	// Check if system has enough RAM (minimum 2GB recommended)
	if err := checkRAM(2); err != nil {
		fmt.Printf("WARNING: %v\n", err)
		fmt.Println("Exiting program due to RAM requirement failure")
		os.Exit(1)
	}

	// Check if system has enough disk space (minimum 10GB recommended)
	if err := checkDiskSpace(10); err != nil {
		fmt.Printf("WARNING: %v\n", err)
		fmt.Println("Exiting program due to disk space requirement failure")
		os.Exit(1)
	}

	fmt.Println("All system requirements checks passed!")
}

func checkDockerInstallation() error {
	cmd := exec.Command("docker", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker is not installed or not accessible: %v", err)
	}

	fmt.Printf("Docker is available: %s\n", strings.TrimSpace(string(output)))
	return nil
}

func checkRAM(minGB float64) error {
	var totalGB float64

	switch runtime.GOOS {
	case "linux":
		// Read /proc/meminfo
		data, err := os.ReadFile("/proc/meminfo")
		if err != nil {
			return fmt.Errorf("failed to read /proc/meminfo: %v", err)
		}
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					kb, err := strconv.ParseFloat(fields[1], 64)
					if err == nil {
						totalGB = kb / (1024 * 1024)
					}
				}
				break
			}
		}
	case "darwin":
		// Use sysctl on macOS
		cmd := exec.Command("sysctl", "-n", "hw.memsize")
		output, err := cmd.Output()
		if err == nil {
			bytes, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
			if err == nil {
				totalGB = bytes / (1024 * 1024 * 1024)
			}
		}
	case "windows":
		// Use wmic on Windows
		cmd := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && line != "TotalPhysicalMemory" {
					bytes, err := strconv.ParseFloat(line, 64)
					if err == nil {
						totalGB = bytes / (1024 * 1024 * 1024)
						break
					}
				}
			}
		}
	default:
		// Fallback to runtime stats (inaccurate but better than nothing)
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		totalGB = float64(m.Sys) / (1024 * 1024 * 1024)
	}

	if totalGB == 0 {
		fmt.Println("Warning: Could not detect system RAM, skipping check.")
		return nil
	}

	if totalGB < minGB {
		return fmt.Errorf("insufficient RAM. Detected: %.2f GB, Required: %.2f GB", totalGB, minGB)
	}

	fmt.Printf("RAM check passed. Total: %.2f GB\n", totalGB)
	return nil
}

func checkDiskSpace(minGB float64) error {
	var path string
	if runtime.GOOS == "windows" {
		path = "C:\\"
	} else {
		path = "/"
	}

	freeGB := 0.0

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command", "Get-PSDrive C | Select-Object Free")
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to check disk space: %v", err)
		}
		
		lines := strings.Split(string(output), "\n")
		// PowerShell output might be "Free\n-----\n123456789\n" or just the number depending on version/formatting
		// We look for the first valid number
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if val, err := strconv.ParseFloat(line, 64); err == nil && val > 0 {
				freeGB = val / (1024 * 1024 * 1024)
				break
			}
		}
	} else {
		// Use df -kP for portability (outputs in 1K blocks)
		cmd := exec.Command("df", "-kP", path)
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to check disk space: %v", err)
		}

		// Parse Unix/Linux output
		// Filesystem 1024-blocks Used Available Capacity Mounted on
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i == 0 {
				continue // Skip header
			}
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				// Available is usually the 4th field in -kP output
				// Linux: Filesystem 1k-blocks Used Available Use% Mounted on
				// macOS: Filesystem 1024-blocks Used Available Capacity Mounted on
				// Sometimes strict POSIX might vary but usually Available is 4th.
				
				// Check which field matches numbers
				idx := 3 // 0-indexed, so 4th field
				if val, err := strconv.ParseFloat(fields[idx], 64); err == nil {
					freeGB = val / (1024 * 1024) // KB to GB
					break
				}
			}
		}
	}

	if freeGB == 0 {
		fmt.Println("Warning: Could not detect free disk space, skipping check.")
		return nil 
	}

	if freeGB < minGB {
		return fmt.Errorf("insufficient disk space. Available: %.2f GB, Required: %.2f GB", freeGB, minGB)
	}

	fmt.Printf("Disk space check passed. Available: %.2f GB\n", freeGB)
	return nil
}
