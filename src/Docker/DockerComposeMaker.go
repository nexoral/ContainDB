package Docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// ContainerInfo holds details about a Docker container
type ContainerInfo struct {
	Name          string
	Image         string
	Ports         []string
	Volumes       []string
	EnvVars       []string
	Networks      []string
	Dependencies  []string
	RestartPolicy string
	Command       string
}

// MakeDockerComposeWithAllServices creates a Docker Compose file from running containers
func MakeDockerComposeWithAllServices() string {
	fmt.Println("Generating Docker Compose file from running containers...")

	// Get all running containers on ContainDB network
	containers, err := ListRunningDatabases()
	if err != nil {
		fmt.Printf("Error listing containers: %v\n", err)
		return ""
	}

	if len(containers) == 0 {
		fmt.Println("No containers found running on ContainDB-Network")
		return ""
	}

	// Map to store container info
	containerInfoMap := make(map[string]ContainerInfo)

	// Get details for each container
	for _, containerName := range containers {
		info, err := getContainerInfo(containerName)
		if err != nil {
			fmt.Printf("Error getting info for container %s: %v\n", containerName, err)
			continue
		}
		containerInfoMap[containerName] = info
	}

	// Generate Docker Compose YAML content
	composeContent := generateComposeYAML(containerInfoMap)

	// Get current working directory to save the file
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return ""
	}

	// Create file path
	filePath := filepath.Join(cwd, "docker-compose.yml")

	// Trim trailing whitespace and newlines to ensure no empty line at the end
	composeContent = strings.TrimRight(composeContent, "\n\r\t ")

	// Write the YAML to file
	err = os.WriteFile(filePath, []byte(composeContent), 0644)
	if err != nil {
		fmt.Printf("Error writing Docker Compose file: %v\n", err)
		return ""
	}

	fmt.Printf("Docker Compose file created at: %s\n", filePath)
	return filePath
}

// getContainerInfo extracts all relevant information from a container
func getContainerInfo(containerName string) (ContainerInfo, error) {
	info := ContainerInfo{
		Name: containerName,
	}

	// Get container image
	cmd := exec.Command("docker", "inspect", "--format", "{{.Config.Image}}", containerName)
	output, err := cmd.Output()
	if err != nil {
		return info, fmt.Errorf("error getting container image: %v", err)
	}
	info.Image = strings.TrimSpace(string(output))

	// Get exposed ports
	cmd = exec.Command("docker", "inspect", "--format", "{{json .NetworkSettings.Ports}}", containerName)
	output, err = cmd.Output()
	if err != nil {
		return info, fmt.Errorf("error getting container ports: %v", err)
	}

	// Simple port extraction (this could be more robust with JSON parsing)
	ports := strings.TrimSpace(string(output))
	if ports != "null" && ports != "{}" {
		// Extract port mappings using another docker inspect command that's easier to parse
		cmd = exec.Command("docker", "inspect", "--format", "{{range $p, $conf := .NetworkSettings.Ports}}{{if $conf}}{{(index $conf 0).HostPort}}:{{$p}}{{end}} {{end}}", containerName)
		output, err = cmd.Output()
		if err == nil {
			portMappings := strings.Fields(strings.TrimSpace(string(output)))
			for _, mapping := range portMappings {
				if mapping != "" {
					// Convert "hostPort:containerPort/tcp" to "hostPort:containerPort"
					mapping = strings.Replace(mapping, "/tcp", "", 1)
					mapping = strings.Replace(mapping, "/udp", "", 1)
					info.Ports = append(info.Ports, mapping)
				}
			}
		}
	}

	// Get environment variables - using a different approach to preserve spaces in values
	cmd = exec.Command("docker", "inspect", "--format", "{{json .Config.Env}}", containerName)
	output, err = cmd.Output()
	if err == nil {
		// The output is a JSON array, remove the surrounding brackets and quotes
		envString := strings.TrimSpace(string(output))
		envString = strings.TrimPrefix(envString, "[")
		envString = strings.TrimSuffix(envString, "]")

		// Split by "," but be careful about commas within quotes
		var envVars []string
		inQuote := false
		lastStart := 0

		for i, c := range envString {
			if c == '"' && (i == 0 || envString[i-1] != '\\') {
				inQuote = !inQuote
			} else if c == ',' && !inQuote {
				part := envString[lastStart:i]
				part = strings.TrimSpace(part)
				if len(part) > 0 {
					// Remove surrounding quotes
					part = strings.TrimPrefix(part, "\"")
					part = strings.TrimSuffix(part, "\"")
					// Unescape escaped quotes
					part = strings.ReplaceAll(part, "\\\"", "\"")
					envVars = append(envVars, part)
				}
				lastStart = i + 1
			}
		}

		// Add the last part
		if lastStart < len(envString) {
			part := envString[lastStart:]
			part = strings.TrimSpace(part)
			part = strings.TrimPrefix(part, "\"")
			part = strings.TrimSuffix(part, "\"")
			part = strings.ReplaceAll(part, "\\\"", "\"")
			if len(part) > 0 {
				envVars = append(envVars, part)
			}
		}

		// Filter environment variables - keep only meaningful ones for container functionality
		var filteredEnvVars []string
		isPhpMyAdmin := strings.Contains(strings.ToLower(containerName), "phpmyadmin")

		// List of system/internal env vars to exclude
		excludedPrefixes := []string{
			"PATH=", "PHPIZE_DEPS=", "PHP_INI_DIR=",
			"APACHE_CONFDIR=", "APACHE_ENVVARS=",
			"PHP_CFLAGS=", "PHP_CPPFLAGS=", "PHP_LDFLAGS=",
			"GPG_KEYS=", "PHP_VERSION=", "PHP_URL=", "PHP_ASC_URL=", "PHP_SHA256=",
			"GOSU_VERSION=", "JSYAML_VERSION=", "JSYAML_CHECKSUM=",
			"MONGO_PACKAGE=", "MONGO_REPO=", "MONGO_MAJOR=", "MONGO_VERSION=",
			"GLIBC_TUNABLES=", "HOME=",
			"VERSION=", "SHA256=", "URL=",
		}

		// For phpMyAdmin, only include specific variables
		if isPhpMyAdmin {
			for _, env := range envVars {
				if strings.HasPrefix(env, "PMA_HOST=") ||
					strings.HasPrefix(env, "PMA_PORT=") ||
					strings.HasPrefix(env, "PMA_USER=") ||
					strings.HasPrefix(env, "PMA_PASSWORD=") ||
					strings.HasPrefix(env, "PMA_DATABASE=") ||
					strings.HasPrefix(env, "PMA_ARBITRARY=") ||
					strings.HasPrefix(env, "PMA_SSL=") ||
					strings.HasPrefix(env, "PMA_SSL_VERIFY=") {
					filteredEnvVars = append(filteredEnvVars, env)
				}
			}
		} else {
			// For other containers, exclude system variables
			for _, env := range envVars {
				exclude := false
				for _, prefix := range excludedPrefixes {
					if strings.HasPrefix(env, prefix) {
						exclude = true
						break
					}
				}
				if !exclude {
					filteredEnvVars = append(filteredEnvVars, env)
				}
			}
		}

		info.EnvVars = filteredEnvVars
	}

	// Get volumes
	cmd = exec.Command("docker", "inspect", "--format", "{{range .Mounts}}{{.Source}}:{{.Destination}} {{end}}", containerName)
	output, err = cmd.Output()
	if err == nil {
		volumes := strings.Fields(strings.TrimSpace(string(output)))
		info.Volumes = volumes
	}

	// Get networks
	cmd = exec.Command("docker", "inspect", "--format", "{{range $k, $v := .NetworkSettings.Networks}}{{$k}} {{end}}", containerName)
	output, err = cmd.Output()
	if err == nil {
		networks := strings.Fields(strings.TrimSpace(string(output)))
		info.Networks = networks
	}

	// Get restart policy
	cmd = exec.Command("docker", "inspect", "--format", "{{.HostConfig.RestartPolicy.Name}}", containerName)
	output, err = cmd.Output()
	if err == nil {
		info.RestartPolicy = strings.TrimSpace(string(output))
	}

	// Get command if any
	cmd = exec.Command("docker", "inspect", "--format", "{{if .Config.Cmd}}{{join .Config.Cmd \" \"}}{{end}}", containerName)
	output, err = cmd.Output()
	if err == nil {
		info.Command = strings.TrimSpace(string(output))
	}

	return info, nil
}

// generateComposeYAML creates a Docker Compose YAML file from container info
func generateComposeYAML(containers map[string]ContainerInfo) string {
	// Template for Docker Compose file
	const composeTemplate = `version: '3'

services:
{{- range $name, $container := . }}
  {{ sanitizeName $name }}:
    image: {{ $container.Image }}
    container_name: {{ $name }}
{{- if $container.RestartPolicy }}
    restart: {{ $container.RestartPolicy }}
{{- end }}
{{- if $container.Command }}
    command: {{ $container.Command }}
{{- end }}
{{- if $container.Ports }}
    ports:
{{- range $container.Ports }}
      - "{{ . }}"
{{- end }}
{{- end }}
{{- if $container.EnvVars }}
    environment:
{{ formatEnvironment $container.EnvVars | indent 6 }}
{{- end }}
{{- if $container.Volumes }}
    volumes:
{{- range $container.Volumes }}
      - {{ . }}
{{- end }}
{{- end }}
{{- if $container.Networks }}
    networks:
{{- range $container.Networks }}
      - {{ . }}
{{- end }}
{{- end }}
{{- end }}

networks:
  ContainDB-Network:
    external: true
`

	// Create template with custom functions
	tmpl, err := template.New("compose").Funcs(template.FuncMap{
		"sanitizeName": func(name string) string {
			// Replace problematic characters in service names
			return strings.ReplaceAll(name, "-", "_")
		},
		"formatEnvironment": func(envVars []string) string {
			var result strings.Builder
			processedVars := make(map[string]bool) // Track processed variables to avoid duplicates

			// Special handling for GPG_KEYS and similar variables with spaces
			for _, envVar := range envVars {
				parts := strings.SplitN(envVar, "=", 2)
				if len(parts) != 2 {
					// If it doesn't have an equals sign, just add it as is
					if !processedVars[envVar] {
						processedVars[envVar] = true
						result.WriteString(fmt.Sprintf("%s: ''\n", envVar))
					}
					continue
				}

				key := parts[0]
				value := parts[1]

				// Skip if we've already processed this key
				if processedVars[key] {
					continue
				}

				processedVars[key] = true

				// Always quote the value to ensure YAML compatibility
				// Use single quotes and escape any single quotes within the value
				escapedValue := strings.Replace(value, "'", "''", -1)
				result.WriteString(fmt.Sprintf("%s: '%s'\n", key, escapedValue))
			}

			return result.String()
		},
		"indent": func(spaces int, text string) string {
			pad := strings.Repeat(" ", spaces)
			lines := strings.Split(text, "\n")
			for i, line := range lines {
				if line != "" {
					lines[i] = pad + line
				}
			}
			return strings.Join(lines, "\n")
		},
	}).Parse(composeTemplate)
	if err != nil {
		fmt.Printf("Error creating template: %v\n", err)
		return ""
	}

	// Execute template
	var result strings.Builder
	err = tmpl.Execute(&result, containers)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return ""
	}

	return result.String()
}
