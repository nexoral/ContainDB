package tools

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Cleanup stops and removes any created containers and temporary artifacts.
func Cleanup() {
	fmt.Println("🧹 Cleaning up resources...")
	
	// remove exited/dead containers - get list first, then remove
	fmt.Println("- Removing failed containers...")
	
	statuses := []string{"exited", "dead", "created"}
	for _, status := range statuses {
		// Get containers with the specific status
		cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("status=%s", status), "--format", "{{.ID}}")
		output, err := cmd.Output()
		if err == nil {
			containerIDs := strings.Fields(strings.TrimSpace(string(output)))
			for _, id := range containerIDs {
				if id != "" {
					rmCmd := exec.Command("docker", "rm", "-f", id)
					rmCmd.Run()
				}
			}
		}
	}
	
	// remove dangling images
	fmt.Println("- Removing dangling images...")
	exec.Command("docker", "image", "prune", "-f").Run()

	// clean up MongoDB Compass download - use cross-platform temp dir
	tempDir := Docker.GetTempDir()
	debPath := filepath.Join(tempDir, "mongodb-compass.deb")
	os.Remove(debPath)

	fmt.Println("✅ Cleanup completed.")
	// exit immediately to prevent any further interactive prompts
	os.Exit(1)
}
