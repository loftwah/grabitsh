package grabitsh

import (
	"os"
	"path/filepath"
	"strings"
)

type DatabaseInfo struct {
	MigrationsPresent bool     `json:"migrations_present"`
	ConfigFiles       []string `json:"config_files"`
	ORMUsed           bool     `json:"orm_used"`
	DatabaseTypes     []string `json:"database_types"`
}

func analyzeDatabaseUsage() DatabaseInfo {
	var dbInfo DatabaseInfo

	dbInfo.MigrationsPresent = dirExists("migrations") || dirExists("db/migrate")

	dbConfigFiles := []string{"config/database.yml", "knexfile.js", "ormconfig.json", "sequelize.config.js"}
	for _, file := range dbConfigFiles {
		if fileExists(file) {
			dbInfo.ConfigFiles = append(dbInfo.ConfigFiles, file)
		}
	}

	ormFiles := []string{"models.py", "*.model.ts", "*.rb", "entity/*.go", "*.entity.ts"}
	for _, pattern := range ormFiles {
		files, _ := filepath.Glob(pattern)
		if len(files) > 0 {
			dbInfo.ORMUsed = true
			break
		}
	}

	// Detect database types
	dbTypes := map[string][]string{
		"PostgreSQL": {"postgres", "postgresql"},
		"MySQL":      {"mysql"},
		"SQLite":     {"sqlite"},
		"MongoDB":    {"mongodb", "mongo"},
		"Redis":      {"redis"},
		"Cassandra":  {"cassandra"},
	}

	for dbType, keywords := range dbTypes {
		for _, file := range dbInfo.ConfigFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			for _, keyword := range keywords {
				if strings.Contains(string(content), keyword) {
					dbInfo.DatabaseTypes = appendUnique(dbInfo.DatabaseTypes, dbType)
					break
				}
			}
		}
	}

	return dbInfo
}
