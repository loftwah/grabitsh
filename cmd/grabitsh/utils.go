package grabitsh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const maxContentLength = 1000

// Utility function to check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

// Utility function to check if a directory exists
func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	return err == nil && info.IsDir()
}

// Utility function to parse basic text files
func parseBasicTextFile(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}
	buffer.WriteString(fmt.Sprintf("%s contents:\n", filename))
	buffer.Write(fileContent)
}

// Utility function to parse JSON files
func parseJSONFile(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(fileContent, &parsed); err != nil {
		buffer.WriteString(fmt.Sprintf("Error parsing %s: %v\n", filename, err))
		return
	}

	buffer.WriteString(fmt.Sprintf("\nParsed %s JSON:\n", filename))
	for key, value := range parsed {
		buffer.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
	}
}

// Utility function to parse YAML files
func parseYAMLFile(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}

	var parsed map[interface{}]interface{}
	if err := yaml.Unmarshal(fileContent, &parsed); err != nil {
		buffer.WriteString(fmt.Sprintf("Error parsing %s: %v\n", filename, err))
		return
	}

	buffer.WriteString(fmt.Sprintf("\nParsed %s YAML:\n", filename))
	for key := range parsed {
		buffer.WriteString(fmt.Sprintf("  %v\n", key))
	}
}

// Utility function to truncate file content if it exceeds a certain length
func truncateContent(content string) string {
	if len(content) > maxContentLength {
		return content[:maxContentLength] + "...\n(content truncated)"
	}
	return content
}

// Utility function to sanitize .env file contents
func sanitizeEnvFile(content string) string {
	lines := strings.Split(content, "\n")
	var sanitized []string
	for _, line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			sanitized = append(sanitized, parts[0]+"=<value>")
		} else {
			sanitized = append(sanitized, line)
		}
	}
	return strings.Join(sanitized, "\n")
}

// Utility function to check for files with specific extensions
func fileExistsWithExtensions(baseName string, extensions []string) bool {
	for _, ext := range extensions {
		if fileExists(baseName + ext) {
			return true
		}
	}
	return false
}

// Utility function to check and parse a file
func checkAndParseFile(filename string, parser func(string, *bytes.Buffer), buffer *bytes.Buffer) {
	if fileExists(filename) {
		buffer.WriteString(fmt.Sprintf("\nFound: %s\n", filename))
		parser(filename, buffer)
	} else {
		buffer.WriteString(fmt.Sprintf("No %s found.\n", filename))
	}
}

// Utility function to parse Git config
func parseGitConfig(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}
	buffer.WriteString(fmt.Sprintf("Git config contents:\n%s", truncateContent(string(fileContent))))
}

// Utility function to parse GitHub Actions workflows
func parseGithubActionsWorkflows(directory string, buffer *bytes.Buffer) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") {
			buffer.WriteString(fmt.Sprintf("\nParsing GitHub Actions workflow: %s\n", path))
			parseYAMLFile(path, buffer)
		}
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error walking through directory %s: %v\n", directory, err))
	}
}

// Utility function to parse Dockerfile
func parseDockerfile(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}
	lines := strings.Split(string(fileContent), "\n")
	buffer.WriteString("Dockerfile (first 10 lines):\n")
	for i := 0; i < 10 && i < len(lines); i++ {
		buffer.WriteString(fmt.Sprintf("%s\n", lines[i]))
	}
}

// Utility function to parse Docker directories
func parseDockerDir(directory string, buffer *bytes.Buffer) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, "Dockerfile") {
			buffer.WriteString(fmt.Sprintf("\nDocker-related file found: %s\n", path))
			parseYAMLFile(path, buffer)
		}
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error walking through Docker directory: %v\n", err))
	}
}

// Utility function to parse Kubernetes files
func parseK8sFiles(directory string, buffer *bytes.Buffer) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			buffer.WriteString(fmt.Sprintf("\nKubernetes file found: %s\n", path))
			parseYAMLFile(path, buffer)
		}
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error walking through Kubernetes directory: %v\n", err))
	}
}

// Utility function to parse Helm chart files
func parseHelmFiles(directory string, buffer *bytes.Buffer) {
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			buffer.WriteString(fmt.Sprintf("\nHelm chart file found: %s\n", path))
			parseYAMLFile(path, buffer)
		}
		return nil
	})
}

// Utility function to parse directories (for cloud providers)
func parseDirectoryContents(directory string, buffer *bytes.Buffer) {
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			buffer.WriteString(fmt.Sprintf("Error walking directory %s: %v\n", directory, err))
			return nil
		}
		buffer.WriteString(fmt.Sprintf("\nDirectory: %s\n", path))
		return nil
	})
}

// Utility function to parse Gemfile
func parseGemfile(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}
	buffer.WriteString("Gemfile contents:\n")
	buffer.Write(fileContent)
}

// Utility function to parse package.json
func parsePackageJSON(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(fileContent, &packageJSON); err != nil {
		buffer.WriteString(fmt.Sprintf("Error parsing %s: %v\n", filename, err))
		return
	}

	buffer.WriteString(fmt.Sprintf("\nParsed %s JSON:\n", filename))
	for key, value := range packageJSON {
		buffer.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
	}
}
