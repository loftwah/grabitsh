package grabitsh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const maxContentLength = 1000

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	return err == nil && info.IsDir()
}

func runCommand(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running command %s %s: %v\n", name, strings.Join(arg, " "), err)
	}
	return string(out)
}

func appendUnique(slice []string, item string) []string {
	for _, element := range slice {
		if element == item {
			return slice
		}
	}
	return append(slice, item)
}

func parseBasicTextFile(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}
	buffer.WriteString(fmt.Sprintf("%s contents:\n", filename))
	buffer.Write(fileContent)
}

type JSONData struct {
	Data map[string]interface{} `json:"data"`
}

func parseJSONFile(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}

	var jsonData JSONData
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		buffer.WriteString(fmt.Sprintf("Error parsing %s: %v\n", filename, err))
		return
	}

	buffer.WriteString(fmt.Sprintf("\nParsed %s JSON:\n", filename))
	for key, value := range jsonData.Data {
		buffer.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
	}
}

func parseYAMLFile(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}

	var parsed map[string]interface{}
	if err := yaml.Unmarshal(fileContent, &parsed); err != nil {
		buffer.WriteString(fmt.Sprintf("Error parsing %s: %v\n", filename, err))
		return
	}

	buffer.WriteString(fmt.Sprintf("\nParsed %s YAML:\n", filename))
	for key, value := range parsed {
		buffer.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
	}
}

func truncateContent(content string) string {
	if len(content) > maxContentLength {
		return content[:maxContentLength] + "...\n(content truncated)"
	}
	return content
}

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

func fileExistsWithExtensions(baseName string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(baseName+ext, ext) {
			if fileExists(baseName + ext) {
				return true
			}
		}
	}
	return false
}

func parseGitConfig(filename string, buffer *bytes.Buffer) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return
	}
	buffer.WriteString(fmt.Sprintf("Git config contents:\n%s", truncateContent(string(fileContent))))
}

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

func parseHelmFiles(directory string, buffer *bytes.Buffer) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			buffer.WriteString(fmt.Sprintf("\nHelm chart file found: %s\n", path))
			parseYAMLFile(path, buffer)
		}
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error walking through Helm directory: %v\n", err))
	}
}

func parseDirectoryContents(directory string, buffer *bytes.Buffer) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			buffer.WriteString(fmt.Sprintf("Error walking directory %s: %v\n", directory, err))
			return nil
		}
		buffer.WriteString(fmt.Sprintf("\nDirectory: %s\n", path))
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error walking through directory %s: %v\n", directory, err))
	}
}

func parseGemfile(filename string, buffer *bytes.Buffer) error {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return err
	}
	buffer.WriteString("Gemfile contents:\n")
	buffer.Write(fileContent)
	return nil
}

func parsePackageJSON(filename string, buffer *bytes.Buffer) error {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error reading %s: %v\n", filename, err))
		return err
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(fileContent, &packageJSON); err != nil {
		buffer.WriteString(fmt.Sprintf("Error parsing %s: %v\n", filename, err))
		return err
	}

	buffer.WriteString(fmt.Sprintf("\nParsed %s JSON:\n", filename))
	for key, value := range packageJSON {
		buffer.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
	}
	return nil
}
