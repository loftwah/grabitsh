package grabitsh

import (
	"os"
	"path/filepath"
	"regexp"
)

type APIInfo struct {
	Files       []string `json:"files"`
	Swagger     bool     `json:"swagger"`
	GraphQL     bool     `json:"graphql"`
	Endpoints   []string `json:"endpoints"`
	HTTPMethods []string `json:"http_methods"`
}

func analyzeAPIStructure() APIInfo {
	var apiInfo APIInfo

	apiPatterns := []string{"*api*.go", "*controller*.rb", "*views*.py", "routes/*.js", "controllers/*.js"}
	for _, pattern := range apiPatterns {
		files, _ := filepath.Glob(pattern)
		apiInfo.Files = append(apiInfo.Files, files...)
	}

	apiInfo.Swagger = fileExists("swagger.json") || fileExists("swagger.yaml")
	apiInfo.GraphQL = fileExists("schema.graphql") || fileExists("schema.gql")

	// Analyze API endpoints and HTTP methods
	for _, file := range apiInfo.Files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Extract endpoints (this is a simple example, you might need more sophisticated regex for your specific use case)
		endpointRegex := regexp.MustCompile(`(GET|POST|PUT|DELETE|PATCH)\s+["']([^"']+)["']`)
		matches := endpointRegex.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			apiInfo.HTTPMethods = appendUnique(apiInfo.HTTPMethods, match[1])
			apiInfo.Endpoints = appendUnique(apiInfo.Endpoints, match[2])
		}
	}

	return apiInfo
}
