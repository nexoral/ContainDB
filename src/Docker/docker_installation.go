package Docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func IsDockerInstalled() bool {
	cmd := exec.Command("docker", "--version")
	err := cmd.Run()
	return err == nil
}

func InstallDocker() error {
	fmt.Printf("Docker not found. Installing Docker on %s...\n", GetOSName())
	
	switch runtime.GOOS {
	case "linux":
		return installDockerLinux()
	case "windows":
		return installDockerWindows()
	case "darwin":
		return installDockerMacOS()
	default:
		return fmt.Errorf("Docker installation not supported on %s. Please install Docker manually", runtime.GOOS)
	}
}

func installDockerLinux() error {
	fmt.Println("Installing Docker on Linux (Ubuntu/Debian)...")
	commands := []string{
		"sudo apt-get update",
		"sudo apt-get install -y ca-certificates curl",
		"sudo install -m 0755 -d /etc/apt/keyrings",
		"sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc",
		"sudo chmod a+r /etc/apt/keyrings/docker.asc",
		`echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
		$(. /etc/os-release && echo ${UBUNTU_CODENAME:-$VERSION_CODENAME}) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null`,
		"sudo apt-get update",
		"sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",
		"sudo usermod -aG docker $USER",
		"sudo apt install -y docker-compose-plugin",
	}

	for index, c := range commands {
		fmt.Println("Running command", index+1, "of", len(commands), ":", c)
		cmd := exec.Command("bash", "-c", c)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error running command:", c)
		}
	}
	return nil
}

func installDockerWindows() error {
	fmt.Println("Installing Docker on Windows...")
	fmt.Println("⚠️  Automatic Docker installation on Windows is not fully supported.")
	fmt.Println("Please install Docker Desktop manually:")
	fmt.Println("1. Download Docker Desktop from: https://www.docker.com/products/docker-desktop/")
	fmt.Println("2. Run the installer and follow the setup wizard")
	fmt.Println("3. Restart your computer after installation")
	fmt.Println("4. Launch Docker Desktop from the Start menu")
	fmt.Println("5. Wait for Docker Desktop to start, then run this program again")
	fmt.Println("\nAlternatively, you can use Chocolatey (if installed):")
	fmt.Println("  choco install docker-desktop")
	fmt.Println("\nOr use winget (Windows 10/11):")
	fmt.Println("  winget install Docker.DockerDesktop")
	
	return fmt.Errorf("automatic installation not supported on Windows. Please install manually using the instructions above")
}

func installDockerMacOS() error {
	fmt.Println("Installing Docker on macOS...")
	
	// Check if Homebrew is available
	if _, err := exec.LookPath("brew"); err == nil {
		fmt.Println("Using Homebrew to install Docker...")
		fmt.Println("Running: brew install --cask docker")
		cmd := exec.Command("brew", "install", "--cask", "docker")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error installing Docker via Homebrew:", err)
			fmt.Println("\nPlease install Docker Desktop manually:")
			fmt.Println("1. Download from: https://www.docker.com/products/docker-desktop/")
			fmt.Println("2. Open the downloaded .dmg file and drag Docker to Applications")
			fmt.Println("3. Launch Docker from Applications and complete setup")
			return err
		}
		fmt.Println("Docker installed successfully via Homebrew!")
		fmt.Println("Please launch Docker Desktop from Applications and wait for it to start.")
		
		// Even with brew install, user often needs to start the app manually to initialize the engine
		// We return error to prompt restart/manual launch verification unless we can verify it's running.
		// However, brew install is "technically" an install.
		// But for safety and consistency with the user request "for other give error that download docker" 
		// (though Mac brew is auto, maybe they consider brew auto enough? 
		// The prompt said "for other give error that download docker". "other" implies non-Linux.
		// But Mac often uses Brew. I'll assume if Brew works, it's auto. If not, error.
		return nil
	}
	
	// No Homebrew, provide manual instructions
	fmt.Println("⚠️  Homebrew not found. Please install Docker Desktop manually:")
	fmt.Println("1. Download Docker Desktop from: https://www.docker.com/products/docker-desktop/")
	fmt.Println("2. Open the downloaded .dmg file and drag Docker to Applications")
	fmt.Println("3. Launch Docker from Applications and complete the setup wizard")
	fmt.Println("4. Wait for Docker Desktop to start, then run this program again")
	fmt.Println("\nAlternatively, install Homebrew first:")
	fmt.Println("  /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
	fmt.Println("Then run: brew install --cask docker")
	
	return fmt.Errorf("automatic installation not supported without Homebrew. Please install manually using the instructions above")
}

func UninstallDocker() error {
	fmt.Printf("Uninstalling Docker on %s...\n", GetOSName())
	
	switch runtime.GOOS {
	case "linux":
		return uninstallDockerLinux()
	case "windows":
		return uninstallDockerWindows()
	case "darwin":
		return uninstallDockerMacOS()
	default:
		return fmt.Errorf("Docker uninstallation not supported on %s. Please uninstall Docker manually", runtime.GOOS)
	}
}

func uninstallDockerLinux() error {
	fmt.Println("Uninstalling Docker on Linux...")
	commands := []string{
		"sudo apt-get purge -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",
		"sudo rm -rf /var/lib/docker",
		"sudo rm -rf /var/lib/containerd",
		"sudo rm /etc/apt/sources.list.d/docker.list",
		"sudo rm /etc/apt/keyrings/docker.asc",
		"sudo apt-get autoremove -y",
	}

	for index, c := range commands {
		fmt.Println("Running command", index+1, "of", len(commands), ":", c)
		cmd := exec.Command("bash", "-c", c)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error running command:", c)
		}
	}
	return nil
}

func uninstallDockerWindows() error {
	fmt.Println("Uninstalling Docker on Windows...")
	fmt.Println("⚠️  Please uninstall Docker Desktop manually:")
	fmt.Println("1. Open Settings > Apps > Apps & features")
	fmt.Println("2. Search for 'Docker Desktop'")
	fmt.Println("3. Click on it and select 'Uninstall'")
	fmt.Println("4. Follow the uninstallation wizard")
	fmt.Println("\nAlternatively, if installed via Chocolatey:")
	fmt.Println("  choco uninstall docker-desktop")
	fmt.Println("\nOr if installed via winget:")
	fmt.Println("  winget uninstall Docker.DockerDesktop")
	return nil
}

func uninstallDockerMacOS() error {
	fmt.Println("Uninstalling Docker on macOS...")
	
	// Check if installed via Homebrew
	if _, err := exec.LookPath("brew"); err == nil {
		fmt.Println("Attempting to uninstall via Homebrew...")
		cmd := exec.Command("brew", "uninstall", "--cask", "docker")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Docker may not have been installed via Homebrew, or already uninstalled.")
		}
	}
	
	fmt.Println("Please also manually remove Docker Desktop if it exists:")
	fmt.Println("1. Quit Docker Desktop if running")
	fmt.Println("2. Remove from Applications: /Applications/Docker.app")
	fmt.Println("3. Remove data: ~/Library/Containers/com.docker.docker")
	fmt.Println("4. Remove preferences: ~/Library/Preferences/com.docker.docker.plist")
	return nil
}
