package Docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ListRunningDatabases returns names of containers on the ContainDB-Network
func ListRunningDatabases() ([]string, error) {
	cmd := exec.Command("docker", "ps",
		"--filter", "network=ContainDB-Network",
		"--format", "{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	return lines, nil
}

// RemoveDatabase forcibly removes the given container,
// optionally deleting its associated data volumes.
func RemoveDatabase(name string) error {
	// ask user whether to delete attached volumes
	deleteVolumes := AskYesNo("Do you want to delete associated data volumes?")

	// Get container type from name (assuming container name follows pattern like "mongodb-container")
	containerType := ""
	if parts := strings.Split(name, "-"); len(parts) > 0 {
		containerType = parts[0] // Extract database type from container name
	}

	// First remove the container itself
	args := []string{"rm", "-f"}
	if deleteVolumes {
		args = append(args, "-v") // This only removes anonymous volumes
	}
	args = append(args, name)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error removing container: %v", err)
	}

	// If user chose to delete volumes, and we could identify the database type, remove named volume
	if deleteVolumes && containerType != "" {
		volumeName := fmt.Sprintf("%s-data", containerType)
		if VolumeExists(volumeName) {
			fmt.Printf("Removing associated volume: %s\n", volumeName)
			if err := RemoveVolume(volumeName); err != nil {
				return fmt.Errorf("error removing volume %s: %v", volumeName, err)
			}
		}
	}

	return nil
}

// ListDatabaseImages returns a list of database images pulled by ContainDB
func ListDatabaseImages() ([]string, error) {
	// List of common database images used by ContainDB
	dbImages := []string{
		// Core databases
		"mongo", "mysql", "postgres", "redis", "mariadb",
		// Management tools (core)
		"phpmyadmin", "dpage/pgadmin4", "redis/redisinsight",
		// Vector databases
		"qdrant/qdrant", "cr.weaviate.io/semitechnologies/weaviate",
		"milvusdb/milvus", "chromadb/chroma", "pgvector/pgvector",
		"redis/redis-stack", "elasticsearch", "opensearchproject/opensearch",
		"marqoai/marqo", "vespaengine/vespa", "typesense/typesense",
		// Management tools (vector DBs)
		"zilliz/attu", "kibana", "opensearchproject/opensearch-dashboards",
	}

	// Build docker command to list images
	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var images []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this image matches any of our database patterns
		for _, dbImg := range dbImages {
			if strings.HasPrefix(line, dbImg) {
				images = append(images, line)
				break
			}
		}
	}

	return images, nil
}

// IsImageInUse checks if the given image is currently used by any running container
func IsImageInUse(image string) (bool, string, error) {
	cmd := exec.Command("docker", "ps", "--format", "{{.Image}} {{.Names}}", "--filter", fmt.Sprintf("ancestor=%s", image))
	output, err := cmd.Output()
	if err != nil {
		return false, "", fmt.Errorf("failed to check if image is in use: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 && lines[0] != "" {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return true, parts[1], nil // Return the container name
		}
		return true, "unknown", nil
	}

	return false, "", nil
}

// RemoveImage removes a Docker image
func RemoveImage(image string) error {
	cmd := exec.Command("docker", "rmi", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove image: %v", err)
	}
	return nil
}

// ListContainDBVolumes returns a list of volumes created by ContainDB
func ListContainDBVolumes() ([]string, error) {
	// Common volume prefixes used by ContainDB
	prefixes := []string{
		// Core databases
		"mongodb-data", "mysql-data", "postgresql-data", "redis-data", "mariadb-data",
		// Vector databases
		"qdrant-data", "weaviate-data", "milvus-data", "chroma-data", "pgvector-data",
		"redis-stack-data", "elasticsearch-data", "opensearch-data",
		"vespa-data", "typesense-data",
	}

	cmd := exec.Command("docker", "volume", "ls", "--format", "{{.Name}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var volumes []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this volume matches any of our database patterns
		for _, prefix := range prefixes {
			if strings.HasPrefix(line, prefix) {
				volumes = append(volumes, line)
				break
			}
		}
	}

	return volumes, nil
}

// IsVolumeInUse checks if the given volume is currently used by any container
func IsVolumeInUse(volume string) (bool, string, error) {
	cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("volume=%s", volume), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false, "", fmt.Errorf("failed to check if volume is in use: %v", err)
	}

	containerName := strings.TrimSpace(string(output))
	if containerName != "" {
		return true, containerName, nil
	}

	return false, "", nil
}
