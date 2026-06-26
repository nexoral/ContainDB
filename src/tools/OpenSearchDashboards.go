package tools

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func StartOpenSearchDashboards() {
	if Docker.IsContainerRunning("opensearch-dashboards-container", true) {
		fmt.Println("OpenSearch Dashboards container is already running.")
		if Docker.AskYesNo("Remove existing OpenSearch Dashboards container and recreate?") {
			fmt.Println("Removing existing OpenSearch Dashboards container...")
			cmd := exec.Command("docker", "rm", "-f", "opensearch-dashboards-container")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error removing OpenSearch Dashboards:", err)
				return
			}
		} else {
			fmt.Println("Keeping existing container. Aborting setup.")
			return
		}
	}

	osContainers := Docker.ListOfContainers([]string{"opensearchproject/opensearch"})
	if len(osContainers) == 0 {
		fmt.Println("No running OpenSearch containers found. Start OpenSearch first.")
		return
	}

	var filtered []string
	for _, name := range osContainers {
		if name != "opensearch-dashboards-container" {
			filtered = append(filtered, name)
		}
	}

	selected := filtered[0]
	if len(filtered) > 1 {
		fmt.Println("Select an OpenSearch container to link with OpenSearch Dashboards:")
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

	port := AskForInput("Enter host port for OpenSearch Dashboards", "5601")

	fmt.Println("Pulling OpenSearch Dashboards Docker image...")
	cmd := exec.Command("docker", "pull", "opensearchproject/opensearch-dashboards:latest")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	_ = cmd.Run()

	fmt.Println("Creating OpenSearch Dashboards container...")
	osHosts := fmt.Sprintf("http://%s:9200", selected)
	cmd = exec.Command("docker", "run",
		"-d",
		"--restart", "unless-stopped",
		"--network", "ContainDB-Network",
		"--name", "opensearch-dashboards-container",
		"-e", fmt.Sprintf("OPENSEARCH_HOSTS=%s", osHosts),
		"-p", fmt.Sprintf("%s:5601", port),
		"opensearchproject/opensearch-dashboards:latest",
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting OpenSearch Dashboards:", err)
	} else {
		fmt.Printf("✅ OpenSearch Dashboards started! Access it at http://localhost:%s\n", port)
		fmt.Printf("   Connected to OpenSearch container: %s\n", selected)
	}
}
