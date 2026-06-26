package tools

import (
	"ContainDB/src/Docker"
	"fmt"
)

// AfterContainerToolInstaller provides post-installation setup for database management tools.
//
// For MySQL/MariaDB, it offers to install or reinstall phpMyAdmin. If phpMyAdmin is
// already running, it asks if the user wants to reinstall it.
//
// For MongoDB, it offers to install MongoDB Compass as a GUI management tool.
//
// For PostgreSQL, it offers to install PgAdmin as a GUI management tool.
//
// Parameters:
//   - database: A string identifying the database type ("mysql", "mariadb", "mongodb", or "postgresql")
//
// The function doesn't return any values but initiates the installation of
// the respective management tool based on user consent.
func AfterContainerToolInstaller(database string) {
	switch database {
	case "mysql", "mariadb":
		// Check if phpMyAdmin is already running
		if Docker.IsContainerRunning("phpmyadmin", true) {
			fmt.Println("phpMyAdmin is already running.")
			consentPhpMyAdmin := Docker.AskYesNo("Do you want to reinstall phpMyAdmin for this database?")
			if consentPhpMyAdmin {
				StartPHPMyAdmin()
			} else {
				fmt.Println("You can reinstall phpMyAdmin later using the 'phpmyadmin' option.")
			}
		} else {
			consentPhpMyAdmin := Docker.AskYesNo("Do you want to install phpMyAdmin for this database?")
			if consentPhpMyAdmin {
				StartPHPMyAdmin()
			} else {
				fmt.Println("You can install phpMyAdmin later using the 'phpmyadmin' option.")
			}
		}
	case "mongodb":
		consentCompass := Docker.AskYesNo("Do you want to install MongoDB Compass?")
		if consentCompass {
			DownloadMongoDBCompass()
		} else {
			fmt.Println("You can install MongoDB Compass later using the 'mongodb compass' option.")
		}
	case "postgresql":
		pgAdminConsent := Docker.AskYesNo("Do you want to install PgAdmin? (yes/no)")
		if pgAdminConsent {
			StartPgAdmin()
		}
	case "redis":
		redisInsightConsent := Docker.AskYesNo("Do you want to install Redis Insight? (yes/no)")
		if redisInsightConsent {
			StartRedisInsight()
		} else {
			fmt.Println("You can install Redis Insight later using the 'redis insight' option.")
		}

	// Vector database tool suggestions
	case "milvus":
		if Docker.AskYesNo("Do you want to install Attu (Milvus Web UI)?") {
			StartAttu()
		} else {
			fmt.Println("You can install Attu later using the 'Attu' option.")
		}
	case "elasticsearch":
		if Docker.AskYesNo("Do you want to install Kibana?") {
			StartKibana()
		} else {
			fmt.Println("You can install Kibana later using the 'Kibana' option.")
		}
	case "opensearch":
		if Docker.AskYesNo("Do you want to install OpenSearch Dashboards?") {
			StartOpenSearchDashboards()
		} else {
			fmt.Println("You can install OpenSearch Dashboards later using the 'OpenSearch Dashboards' option.")
		}
	case "pgvector":
		if Docker.AskYesNo("Do you want to install PgAdmin?") {
			StartPgAdmin()
		} else {
			fmt.Println("You can install PgAdmin later using the 'PgAdmin' option.")
		}
	case "qdrant":
		fmt.Println("Qdrant Web UI is built-in — access it at http://localhost:6333/dashboard")
	case "redis-stack":
		fmt.Println("RedisInsight is built-in in Redis Stack — access it at http://localhost:8001")

	default:
		fmt.Println("No additional tools available for this database type.")
	}
}
