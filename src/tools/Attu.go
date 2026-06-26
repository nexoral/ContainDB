package tools

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func StartAttu() {
	if Docker.IsContainerRunning("attu-container", true) {
		fmt.Println("Attu container is already running.")
		if Docker.AskYesNo("Remove existing Attu container and recreate?") {
			fmt.Println("Removing existing Attu container...")
			cmd := exec.Command("docker", "rm", "-f", "attu-container")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error removing Attu:", err)
				return
			}
		} else {
			fmt.Println("Keeping existing container. Aborting setup.")
			return
		}
	}

	milvusContainers := Docker.ListOfContainers([]string{"milvusdb/milvus"})
	if len(milvusContainers) == 0 {
		fmt.Println("No running Milvus containers found. Start Milvus first.")
		return
	}

	var filtered []string
	for _, name := range milvusContainers {
		if name != "attu-container" {
			filtered = append(filtered, name)
		}
	}

	selected := filtered[0]
	if len(filtered) > 1 {
		fmt.Println("Select a Milvus container to link with Attu:")
		for i, name := range filtered {
			fmt.Printf("  [%d] %s\n", i+1, name)
		}
		choice := AskForInput("Enter number", "1")
		idx := 0
		for i := range filtered {
			if fmt.Sprintf("%d", i+1) == strings.TrimSpace(choice) {
				idx = i
				break
			}
		}
		selected = filtered[idx]
	}

	port := AskForInput("Enter host port for Attu", "3000")

	fmt.Println("Pulling Attu Docker image...")
	cmd := exec.Command("docker", "pull", "zilliz/attu:latest")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	_ = cmd.Run()

	fmt.Println("Creating Attu container...")
	milvusURL := fmt.Sprintf("http://%s:19530", selected)
	cmd = exec.Command("docker", "run",
		"-d",
		"--restart", "unless-stopped",
		"--network", "ContainDB-Network",
		"--name", "attu-container",
		"-e", fmt.Sprintf("MILVUS_URL=%s", milvusURL),
		"-p", fmt.Sprintf("%s:3000", port),
		"zilliz/attu:latest",
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting Attu:", err)
	} else {
		fmt.Printf("✅ Attu started! Access it at http://localhost:%s\n", port)
		fmt.Printf("   Connected to Milvus container: %s\n", selected)
	}
}
