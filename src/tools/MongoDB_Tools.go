package tools

import (
	"ContainDB/src/Docker"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func DownloadMongoDBCompass() {
	// MongoDB Compass installation is only supported on Linux
	if runtime.GOOS != "linux" {
		fmt.Println("⚠️  MongoDB Compass GUI installation is currently only supported on Linux.")
		fmt.Println("Please download MongoDB Compass manually for your platform:")
		
		switch runtime.GOOS {
		case "windows":
			fmt.Println("Windows: Download from https://www.mongodb.com/try/download/compass")
		case "darwin":
			fmt.Println("macOS: Download from https://www.mongodb.com/try/download/compass or use Homebrew: brew install mongodb-compass")
		}
		return
	}

	fmt.Println("Downloading MongoDB Compass...")
	
	downloadURL := "https://downloads.mongodb.com/compass/mongodb-compass_1.46.2_amd64.deb"
	tempDir := Docker.GetTempDir()
	debPath := filepath.Join(tempDir, "mongodb-compass.deb")

	// Download using Go HTTP client (cross-platform)
	fmt.Println("Downloading from:", downloadURL)
	resp, err := http.Get(downloadURL)
	if err != nil {
		fmt.Printf("Error downloading MongoDB Compass: %v\n", err)
		fmt.Println("Please download it manually from: https://www.mongodb.com/try/download/compass")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error downloading MongoDB Compass: HTTP %d\n", resp.StatusCode)
		return
	}

	// Create the file
	out, err := os.Create(debPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Println("Download completed:", debPath)

	// Install the downloaded deb file (Linux only)
	installCmd := exec.Command("sudo", "dpkg", "-i", debPath)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Printf("Error installing MongoDB Compass: %v\n", err)
		fmt.Println("You may need to install dependencies. Try: sudo apt-get install -f")
	} else {
		fmt.Println("MongoDB Compass downloaded and installed successfully.")

		// Clean up the downloaded file
		if err := os.Remove(debPath); err != nil {
			fmt.Printf("Error cleaning up downloaded file: %v\n", err)
		} else {
			fmt.Println("Temporary files cleaned up successfully.")
		}
		fmt.Println("You can now launch MongoDB Compass from your applications menu or by running 'mongodb-compass' in the terminal.")
		fmt.Println("Note: If you encounter any issues, please ensure you have the necessary dependencies installed.")
		fmt.Println("For more information, visit: https://www.mongodb.com/docs/compass/current/install/")
		fmt.Println("Enjoy using MongoDB Compass!")
	}
}
