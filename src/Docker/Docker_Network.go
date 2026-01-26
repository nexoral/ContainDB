package Docker

import (
	"os"
	"os/exec"
)

func CreateDockerNetworkIfNotExists() error {
	// Check if network exists
	cmdCheck := exec.Command("docker", "network", "inspect", "ContainDB-Network")
	err := cmdCheck.Run()
	if err == nil {
		// Network exists, no need to create
		return nil
	}
	// Network does not exist, create it
	cmdCreate := exec.Command("docker", "network", "create", "ContainDB-Network")
	cmdCreate.Stdout = os.Stdout
	cmdCreate.Stderr = os.Stderr
	return cmdCreate.Run()
}
