package base

import (
	"ContainDB/src/tools"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

var categoryItems = map[string][]string{
	"SQL Database": {
		"mysql", "postgresql", "mariadb", "pgvector",
		"phpmyadmin", "PgAdmin",
		"Back",
	},
	"NoSQL Database": {
		"mongodb", "axiodb", "redis",
		"MongoDB Compass", "Redis Insight",
		"Back",
	},
	"Vector Database": {
		"qdrant", "weaviate", "milvus", "chroma", "redis-stack",
		"elasticsearch", "opensearch", "marqo", "vespa", "typesense",
		"Attu", "Kibana", "OpenSearch Dashboards",
		"Back",
	},
}

func SelectDatabase() string {
	categories := []string{"SQL Database", "NoSQL Database", "Vector Database", "Exit"}

	for {
		categoryPrompt := promptui.Select{
			Label: "Select database category",
			Items: categories,
		}
		_, category, err := categoryPrompt.Run()
		if err != nil {
			fmt.Println("\n⚠️ Interrupt received, rolling back...")
			tools.Cleanup()
		}
		if category == "Exit" {
			fmt.Println("Goodbye!")
			os.Exit(0)
		}

		subItems := categoryItems[category]
		subPrompt := promptui.Select{
			Label: fmt.Sprintf("Select service to start (%s)", category),
			Items: subItems,
		}
		_, result, err := subPrompt.Run()
		if err != nil {
			fmt.Println("\n⚠️ Interrupt received, rolling back...")
			tools.Cleanup()
		}
		if result == "Back" {
			continue
		}

		return result
	}
}
