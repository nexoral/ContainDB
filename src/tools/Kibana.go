package tools

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func StartKibana() {
	if Docker.IsContainerRunning("kibana-container", true) {
		fmt.Println("Kibana container is already running.")
		if Docker.AskYesNo("Remove existing Kibana container and recreate?") {
			fmt.Println("Removing existing Kibana container...")
			cmd := exec.Command("docker", "rm", "-f", "kibana-container")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error removing Kibana:", err)
				return
			}
		} else {
			fmt.Println("Keeping existing container. Aborting setup.")
			return
		}
	}

	esContainers := Docker.ListOfContainers([]string{"elasticsearch"})
	if len(esContainers) == 0 {
		fmt.Println("No running Elasticsearch containers found. Start Elasticsearch first.")
		return
	}

	var filtered []string
	for _, name := range esContainers {
		if name != "kibana-container" {
			filtered = append(filtered, name)
		}
	}

	selected := filtered[0]
	if len(filtered) > 1 {
		fmt.Println("Select an Elasticsearch container to link with Kibana:")
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

	port := AskForInput("Enter host port for Kibana", "5601")

	fmt.Println("Pulling Kibana Docker image...")
	cmd := exec.Command("docker", "pull", "kibana:latest")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	_ = cmd.Run()

	fmt.Println("Creating Kibana container...")
	esHosts := fmt.Sprintf("http://%s:9200", selected)
	cmd = exec.Command("docker", "run",
		"-d",
		"--restart", "unless-stopped",
		"--network", "ContainDB-Network",
		"--name", "kibana-container",
		"-e", fmt.Sprintf("ELASTICSEARCH_HOSTS=%s", esHosts),
		"-p", fmt.Sprintf("%s:5601", port),
		"kibana:latest",
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting Kibana:", err)
	} else {
		fmt.Printf("✅ Kibana started! Access it at http://localhost:%s\n", port)
		fmt.Printf("   Connected to Elasticsearch container: %s\n", selected)
	}
}
