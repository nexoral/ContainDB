package base

import (
	"ContainDB/src/Docker"
	"ContainDB/src/tools"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func StartContainer(database string) {
	imageMap := map[string]string{
		// Core databases
		"mongodb":    "mongo",
		"redis":      "redis",
		"mysql":      "mysql",
		"postgresql": "postgres",
		"mariadb":    "mariadb",
		"axiodb":     "theankansaha/axiodb",
		// Vector databases
		"qdrant":        "qdrant/qdrant",
		"weaviate":      "cr.weaviate.io/semitechnologies/weaviate",
		"milvus":        "milvusdb/milvus",
		"chroma":        "chromadb/chroma",
		"pgvector":      "pgvector/pgvector:pg17",
		"redis-stack":   "redis/redis-stack",
		"elasticsearch": "elasticsearch",
		"opensearch":    "opensearchproject/opensearch",
		"marqo":         "marqoai/marqo",
		"vespa":         "vespaengine/vespa",
		"typesense":     "typesense/typesense",
	}

	defaultPorts := map[string]string{
		// Core databases
		"mongodb":    "27017",
		"redis":      "6379",
		"mysql":      "3306",
		"postgresql": "5432",
		"mariadb":    "3306",
		"axiodb":     "27018",
		// Vector databases
		"qdrant":        "6333",
		"weaviate":      "8080",
		"milvus":        "19530",
		"chroma":        "8000",
		"pgvector":      "5432",
		"redis-stack":   "6379",
		"elasticsearch": "9200",
		"opensearch":    "9200",
		"marqo":         "8882",
		"vespa":         "8080",
		"typesense":     "8108",
	}

	// Secondary ports exposed automatically alongside the primary port.
	secondaryPorts := map[string]string{
		"qdrant":        "6334",  // gRPC
		"weaviate":      "50051", // gRPC
		"milvus":        "9091",  // metrics / management
		"redis-stack":   "8001",  // built-in RedisInsight UI
		"elasticsearch": "9300",  // cluster transport
		"opensearch":    "9600",  // performance analyzer
		"vespa":         "19071", // config server
		"axiodb":        "27019", // internal port
	}

	image := imageMap[database]
	port := defaultPorts[database]

	if Docker.IsContainerRunning(image, false) {
		fmt.Printf("Database %s is already running on port %s\n", database, port)
		return
	}

	// Pull image
	fmt.Printf("Pulling image %s...\n", image)
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	// Ask for port mapping
	portMapping := ""
	if Docker.AskYesNo("Do you want to map container port with host?") {
		customPort := Docker.AskYesNo("Do you want to use custom host port?")
		if customPort {
			hostPort := tools.AskForInput("Enter custom host port", port)
			// For axiodb, reject 27019 as it's needed internally
			if database == "axiodb" && hostPort == "27019" {
				fmt.Println("Port 27019 is reserved for internal use. Please choose a different port.")
				hostPort = tools.AskForInput("Enter custom host port", port)
			}
			portMapping = fmt.Sprintf("-p %s:%s", hostPort, port)
		} else {
			portMapping = fmt.Sprintf("-p %s:%s", port, port)
		}
	}

	// Automatically add secondary ports for databases that need them
	if secPort, ok := secondaryPorts[database]; ok {
		if portMapping != "" {
			portMapping += fmt.Sprintf(" -p %s:%s", secPort, secPort)
		} else {
			portMapping = fmt.Sprintf("-p %s:%s", secPort, secPort)
		}
	}

	restartFlag := ""
	if Docker.AskYesNo("Do you want the container to auto-restart on system startup?") {
		restartFlag = "--restart unless-stopped"
	}

	// Ask for data persistence
	volumeMapping := ""
	if Docker.AskYesNo("Do you want to persist data?") {
		containerDirs := map[string]string{
			// Core databases
			"mongodb":    "/data/db",
			"redis":      "/data",
			"mysql":      "/var/lib/mysql",
			"postgresql": "/var/lib/postgresql/data",
			"mariadb":    "/var/lib/mysql",
			"axiodb":     "/app/AxioDB",
			// Vector databases
			"qdrant":        "/qdrant/storage",
			"weaviate":      "/var/lib/weaviate",
			"milvus":        "/var/lib/milvus",
			"chroma":        "/chroma/chroma",
			"pgvector":      "/var/lib/postgresql/data",
			"redis-stack":   "/data",
			"elasticsearch": "/usr/share/elasticsearch/data",
			"opensearch":    "/usr/share/opensearch/data",
			"vespa":         "/opt/vespa/var",
			"typesense":     "/data",
			// marqo excluded: no reliable standalone volume path
		}
		if _, hasDir := containerDirs[database]; hasDir {
			volName := fmt.Sprintf("%s-data", database)
			if Docker.VolumeExists(volName) {
				items := []string{"Use existing", "Create fresh", "Exit"}
				prompt := promptui.Select{
					Label: fmt.Sprintf("Volume '%s' exists. Use or recreate?", volName),
					Items: items,
				}
				_, choice, _ := prompt.Run()
				if choice == "Create fresh" {
					fmt.Println("Removing and recreating volume:", volName)
					_ = Docker.RemoveVolume(volName)
					_ = Docker.CreateVolume(volName)
				}
				if choice == "Exit" {
					fmt.Println("Exiting setup.")
					return
				}
			} else {
				_ = Docker.CreateVolume(volName)
			}
			volumeMapping = fmt.Sprintf("-v %s:%s", volName, containerDirs[database])
		} else {
			fmt.Printf("⚠️  Data persistence is not supported for %s in this mode.\n", database)
		}
	}

	env := ""
	switch database {
	case "mysql":
		fmt.Println("You need to set environment variables for MySQL.")
		pass := tools.AskForInput("Enter root password", "")
		if pass == "" {
			fmt.Println("Error: Password cannot be empty.")
			os.Exit(1)
		}
		env = fmt.Sprintf("-e MYSQL_ROOT_PASSWORD=%s", pass)

	case "postgresql":
		fmt.Println("You need to set environment variables for PostgreSQL.")
		user := tools.AskForInput("Enter username", "postgres")
		pass := tools.AskForInput("Enter password", "")
		if user == "" {
			user = "postgres"
		}
		if pass == "" {
			fmt.Println("Error: Password cannot be empty.")
			os.Exit(1)
		}
		env = fmt.Sprintf("-e POSTGRES_PASSWORD=%s -e POSTGRES_USER=%s", pass, user)

	case "mariadb":
		fmt.Println("You need to set environment variables for MariaDB.")
		pass := tools.AskForInput("Enter root password", "")
		if pass == "" {
			fmt.Println("Error: Password cannot be empty.")
			os.Exit(1)
		}
		env = fmt.Sprintf("-e MARIADB_ROOT_PASSWORD=%s", pass)

	case "pgvector":
		fmt.Println("You need to set environment variables for pgvector (PostgreSQL-compatible).")
		user := tools.AskForInput("Enter username", "postgres")
		pass := tools.AskForInput("Enter password", "")
		if user == "" {
			user = "postgres"
		}
		if pass == "" {
			fmt.Println("Error: Password cannot be empty.")
			os.Exit(1)
		}
		env = fmt.Sprintf("-e POSTGRES_PASSWORD=%s -e POSTGRES_USER=%s", pass, user)

	case "elasticsearch":
		fmt.Println("Configuring Elasticsearch (single-node mode).")
		if Docker.AskYesNo("Enable security (password-protected)?") {
			pass := tools.AskForInput("Enter ELASTIC_PASSWORD", "")
			if pass == "" {
				fmt.Println("Error: Password cannot be empty.")
				os.Exit(1)
			}
			env = fmt.Sprintf("-e discovery.type=single-node -e ELASTIC_PASSWORD=%s", pass)
		} else {
			env = "-e discovery.type=single-node -e xpack.security.enabled=false"
			fmt.Println("⚠️  Security disabled — dev mode only, do not use in production.")
		}

	case "opensearch":
		fmt.Println("Configuring OpenSearch (single-node mode).")
		fmt.Println("Admin password must be ≥8 chars with at least 1 uppercase, 1 number, and 1 special character.")
		pass := tools.AskForInput("Enter OPENSEARCH_INITIAL_ADMIN_PASSWORD", "")
		if pass == "" {
			fmt.Println("Error: Password cannot be empty.")
			os.Exit(1)
		}
		env = fmt.Sprintf("-e discovery.type=single-node -e OPENSEARCH_INITIAL_ADMIN_PASSWORD=%s", pass)

	case "weaviate":
		env = "-e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true -e PERSISTENCE_DATA_PATH=/var/lib/weaviate"

	case "typesense":
		fmt.Println("Typesense requires an API key for all requests.")
		apiKey := tools.AskForInput("Enter Typesense API key", "")
		if apiKey == "" {
			fmt.Println("Error: API key cannot be empty.")
			os.Exit(1)
		}
		env = fmt.Sprintf("-e TYPESENSE_DATA_DIR=/data -e TYPESENSE_API_KEY=%s", apiKey)
	}

	containerName := fmt.Sprintf("%s-container", database)

	// Build docker run command as args array
	args := []string{"run", "-d", "--network", "ContainDB-Network"}

	if portMapping != "" {
		args = append(args, strings.Fields(portMapping)...)
	}

	if restartFlag != "" {
		args = append(args, strings.Fields(restartFlag)...)
	}

	if volumeMapping != "" {
		args = append(args, strings.Fields(volumeMapping)...)
	}

	if env != "" {
		args = append(args, strings.Fields(env)...)
	}

	args = append(args, "--name", containerName, image)

	// Some databases need a startup command appended after the image name
	containerCommands := map[string]string{
		"milvus": "milvus run standalone",
	}
	if cmdStr, ok := containerCommands[database]; ok {
		args = append(args, strings.Fields(cmdStr)...)
	}

	fmt.Println("Running: docker", strings.Join(args, " "))
	cmd = exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting container:", err)
	} else {
		fmt.Println("Container started successfully.")
		tools.AfterContainerToolInstaller(database)
	}
}
