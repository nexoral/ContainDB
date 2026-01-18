package Docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func AskYesNo(label string) bool {
	items := []string{"Yes", "No", "Exit"}
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		// Handle cleanup locally or call a function that doesn't create an import cycle
		os.Exit(1)
	}
	if index == len(items)-1 {
		fmt.Println("Exiting...")
		os.Exit(0)
	}
	return index == 0
}

func IsContainerRunning(nameOrImage string, checkByName bool) bool {
	var cmd *exec.Cmd
	if checkByName {
		cmd = exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", nameOrImage), "--format", "{{.Names}}")
	} else {
		cmd = exec.Command("docker", "ps", "--filter", fmt.Sprintf("ancestor=%s", nameOrImage), "--format", "{{.Names}}")
	}
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) != ""
}

func ListOfContainers(images []string) []string {
	if len(images) == 0 {
		return []string{}
	}

	// Get all running containers with their names and images
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}} {{.Image}}")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	// Filter containers by matching image names
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var containers []string
	imageMap := make(map[string]bool)
	for _, img := range images {
		imageMap[img] = true
	}

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			containerName := parts[0]
			containerImage := parts[1]
			// Check if any of the target images match
			for img := range imageMap {
				if strings.Contains(containerImage, img) {
					containers = append(containers, containerName)
					break
				}
			}
		}
	}
	return containers
}

// VolumeExists returns true if Docker volume with given name exists
func VolumeExists(name string) bool {
	cmd := exec.Command("docker", "volume", "inspect", name)
	err := cmd.Run()
	return err == nil
}

// CreateVolume creates a Docker volume with given name
func CreateVolume(name string) error {
	cmd := exec.Command("docker", "volume", "create", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemoveVolume force-removes a Docker volume with given name
func RemoveVolume(name string) error {
	// First check if volume exists
	if !VolumeExists(name) {
		return fmt.Errorf("volume %s does not exist", name)
	}

	fmt.Printf("Removing volume %s...\n", name)
	cmd := exec.Command("docker", "volume", "rm", "-f", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove volume: %v", err)
	}
	return nil
}
