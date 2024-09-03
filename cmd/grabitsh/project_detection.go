package grabitsh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const maxContentLength = 1000 // Maximum number of characters to display for file contents

func DetectProjectTypes() []string {
	var projectTypes []string

	if fileExists("package.json") {
		projectTypes = append(projectTypes, "Node.js project")
		if fileExistsWithExtensions("next.config", []string{".js", ".ts", ".mjs", ".mts"}) {
			projectTypes = append(projectTypes, "Next.js framework")
		}
		if fileExists("react-scripts.config.js") || (dirExists("src") && fileExists("src/App.js")) {
			projectTypes = append(projectTypes, "React project")
		}
		if fileExistsWithExtensions("astro.config", []string{".js", ".ts", ".mjs", ".mts"}) {
			projectTypes = append(projectTypes, "Astro framework")
		}
		if fileExistsWithExtensions("vite.config", []string{".js", ".ts", ".mjs", ".mts"}) {
			projectTypes = append(projectTypes, "Vite project")
		}
	}

	if fileExists("config/application.rb") && dirExists("app") && dirExists("config") {
		projectTypes = append(projectTypes, "Ruby on Rails project")
	}

	if fileExists("artisan") && dirExists("app") && dirExists("public") {
		projectTypes = append(projectTypes, "Laravel (PHP) project")
	}

	if fileExists("manage.py") && dirExists("templates") {
		projectTypes = append(projectTypes, "Django (Python) project")
	}

	if fileExists("app.py") || fileExists("wsgi.py") {
		projectTypes = append(projectTypes, "Flask/FastAPI (Python) project")
	}

	if fileExistsWithExtensions("vue.config", []string{".js", ".ts"}) {
		projectTypes = append(projectTypes, "Vue.js project")
	}

	if fileExists("angular.json") {
		projectTypes = append(projectTypes, "Angular project")
	}

	if fileExists("Program.cs") && dirExists("bin") && dirExists("obj") {
		projectTypes = append(projectTypes, ".NET Core project")
	}

	if fileExists("pom.xml") && dirExists("src/main/java") {
		projectTypes = append(projectTypes, "Java Spring Boot project")
	}

	if fileExists("go.mod") {
		projectTypes = append(projectTypes, "Go project")
	}

	if dirExists("terraform") || fileExists("main.tf") {
		projectTypes = append(projectTypes, "Terraform project")
	}

	if fileExistsWithExtensions("docker-compose", []string{".yml", ".yaml"}) || fileExists("compose.yml") || fileExists("compose.yaml") {
		projectTypes = append(projectTypes, "Docker Compose project")
	}

	if fileExists("Dockerfile") {
		projectTypes = append(projectTypes, "Docker project")
	}

	if fileExists("Vagrantfile") {
		projectTypes = append(projectTypes, "Vagrant project")
	}

	if fileExists("ansible.cfg") || dirExists("roles") {
		projectTypes = append(projectTypes, "Ansible project")
	}

	if fileExists("Jenkinsfile") {
		projectTypes = append(projectTypes, "Jenkins pipeline")
	}

	if fileExists("cloudbuild.yaml") || fileExists("cloudbuild.yml") {
		projectTypes = append(projectTypes, "Google Cloud Build project")
	}

	if fileExists("serverless.yml") || fileExists("serverless.yaml") {
		projectTypes = append(projectTypes, "Serverless Framework project")
	}

	if fileExists("Chart.yaml") {
		projectTypes = append(projectTypes, "Helm Chart")
	}

	return projectTypes
}

func AnalyzeRepository() string {
	var output strings.Builder

	// Analyze root directory
	output.WriteString("### Repository Structure ###\n")
	analyzeDirectory(".", &output, 0)

	// Analyze specific directories and files
	analyzeGitHubDir(&output)
	analyzeImportantDirs(&output)
	analyzeImportantFiles(&output)
	analyzeGoProject(&output)
	analyzeDependencies(&output)
	analyzeConfiguration(&output)
	analyzeDocumentation(&output)
	analyzeContainerization(&output)
	analyzeInfrastructureAsCode(&output)
	analyzeCICDPipelines(&output)

	return output.String()
}

func analyzeDirectory(dir string, output *strings.Builder, depth int) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(output, "Error reading directory %s: %v\n", dir, err)
		return
	}

	for _, file := range files {
		indent := strings.Repeat("  ", depth)
		path := filepath.Join(dir, file.Name())

		if file.IsDir() {
			fmt.Fprintf(output, "%s📁 %s\n", indent, file.Name())
			if depth < 2 { // Limit depth to avoid excessive output
				analyzeDirectory(path, output, depth+1)
			}
		} else {
			fmt.Fprintf(output, "%s📄 %s\n", indent, file.Name())
		}
	}
}

func analyzeGitHubDir(output *strings.Builder) {
	if dirExists(".github") {
		output.WriteString("\n### .github Directory Analysis ###\n")
		if dirExists(".github/workflows") {
			output.WriteString("GitHub Actions workflows found:\n")
			workflows, _ := filepath.Glob(".github/workflows/*.yml")
			for _, workflow := range workflows {
				content, err := ioutil.ReadFile(workflow)
				if err == nil {
					output.WriteString(fmt.Sprintf("Workflow: %s\n", filepath.Base(workflow)))
					output.WriteString(truncateContent(string(content)))
					output.WriteString("\n\n")
				}
			}
		}
		if fileExists(".github/PULL_REQUEST_TEMPLATE.md") {
			output.WriteString("Pull Request template found\n")
		}
		if fileExists(".github/FUNDING.yml") {
			output.WriteString("Funding configuration found\n")
		}
		if fileExists(".github/CODEOWNERS") {
			output.WriteString("CODEOWNERS file found\n")
		}
	}
}

func analyzeImportantDirs(output *strings.Builder) {
	importantDirs := []string{"app", "src", "config", "lib", "spec", "test", "public", "cmd"}
	for _, dir := range importantDirs {
		if dirExists(dir) {
			output.WriteString(fmt.Sprintf("\n### %s Directory Contents ###\n", dir))
			analyzeDirectory(dir, output, 0)
		}
	}

	if dirExists("terraform") {
		output.WriteString("\n### Terraform Files ###\n")
		tfFiles, _ := filepath.Glob("terraform/*.tf")
		for _, file := range tfFiles {
			content, err := ioutil.ReadFile(file)
			if err == nil {
				output.WriteString(fmt.Sprintf("File: %s\n", filepath.Base(file)))
				output.WriteString(truncateContent(string(content)))
				output.WriteString("\n\n")
			}
		}
	}
}

func analyzeImportantFiles(output *strings.Builder) {
	importantFiles := []string{
		".dockerignore", ".gitignore", "Dockerfile",
		"Procfile", "Rakefile", "Makefile", ".env", "package.json",
		"Gemfile", "requirements.txt",
		"go.mod", "go.sum", "main.go", "README.md", "LICENSE",
		"Vagrantfile", "ansible.cfg", "Jenkinsfile", "cloudbuild.yaml",
		"serverless.yml", "Chart.yaml",
	}

	output.WriteString("\n### Important Files ###\n")
	for _, file := range importantFiles {
		if fileExists(file) {
			content, err := ioutil.ReadFile(file)
			if err == nil {
				output.WriteString(fmt.Sprintf("File: %s\n", file))
				if file == ".env" {
					output.WriteString(sanitizeEnvFile(string(content)))
				} else {
					output.WriteString(truncateContent(string(content)))
				}
				output.WriteString("\n\n")
			}
		}
	}

	// Check for files with multiple possible extensions
	multiExtensionFiles := map[string][]string{
		"docker-compose": {".yml", ".yaml"},
		"compose":        {".yml", ".yaml"},
		"vite.config":    {".js", ".ts", ".mjs", ".mts"},
		"astro.config":   {".js", ".ts", ".mjs", ".mts"},
		"next.config":    {".js", ".ts", ".mjs", ".mts"},
	}

	for baseName, extensions := range multiExtensionFiles {
		for _, ext := range extensions {
			fileName := baseName + ext
			if fileExists(fileName) {
				content, err := ioutil.ReadFile(fileName)
				if err == nil {
					output.WriteString(fmt.Sprintf("File: %s\n", fileName))
					output.WriteString(truncateContent(string(content)))
					output.WriteString("\n\n")
				}
				break // We only need to find one matching file
			}
		}
	}
}

func analyzeGoProject(output *strings.Builder) {
	if fileExists("go.mod") {
		output.WriteString("\n### Go Project Analysis ###\n")

		// Analyze go.mod
		modContent, _ := ioutil.ReadFile("go.mod")
		output.WriteString("go.mod contents:\n")
		output.WriteString(truncateContent(string(modContent)))
		output.WriteString("\n\n")

		// Analyze main.go if it exists
		if fileExists("main.go") {
			mainContent, _ := ioutil.ReadFile("main.go")
			output.WriteString("main.go contents:\n")
			output.WriteString(truncateContent(string(mainContent)))
			output.WriteString("\n\n")
		}

		// List all Go files
		output.WriteString("Go files in the project:\n")
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
				output.WriteString(fmt.Sprintf("- %s\n", path))
			}
			return nil
		})
		if err != nil {
			output.WriteString(fmt.Sprintf("Error walking the path: %v\n", err))
		}
	}
}

func analyzeDependencies(output *strings.Builder) {
	output.WriteString("\n### Dependencies Analysis ###\n")

	// Analyze package.json for Node.js projects
	if fileExists("package.json") {
		content, _ := ioutil.ReadFile("package.json")
		var packageJSON map[string]interface{}
		if err := json.Unmarshal(content, &packageJSON); err == nil {
			if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
				output.WriteString("Node.js Dependencies:\n")
				for dep, version := range deps {
					output.WriteString(fmt.Sprintf("- %s: %v\n", dep, version))
				}
			}
		}
	}

	// Analyze go.mod for Go projects
	if fileExists("go.mod") {
		content, _ := ioutil.ReadFile("go.mod")
		output.WriteString("Go Dependencies:\n")
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "\t") && !strings.Contains(line, "=>") {
				output.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(line)))
			}
		}
	}
}

func analyzeConfiguration(output *strings.Builder) {
	output.WriteString("\n### Configuration Analysis ###\n")

	// Analyze .env file
	if fileExists(".env") {
		content, _ := ioutil.ReadFile(".env")
		output.WriteString("Environment variables (sanitized):\n")
		output.WriteString(sanitizeEnvFile(string(content)))
		output.WriteString("\n")
	}

	// Analyze YAML configuration files
	yamlFiles, err := filepath.Glob("*.yaml")
	if err != nil {
		output.WriteString(fmt.Sprintf("Error searching for YAML files: %v\n", err))
		return
	}
	ymlFiles, err := filepath.Glob("*.yml")
	if err != nil {
		output.WriteString(fmt.Sprintf("Error searching for YML files: %v\n", err))
		return
	}
	yamlFiles = append(yamlFiles, ymlFiles...)

	for _, file := range yamlFiles {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			output.WriteString(fmt.Sprintf("Error reading file %s: %v\n", file, err))
			continue
		}
		var yamlConfig map[string]interface{}
		if err := yaml.Unmarshal(content, &yamlConfig); err == nil {
			output.WriteString(fmt.Sprintf("YAML Configuration (%s):\n", file))
			output.WriteString(truncateContent(fmt.Sprintf("%v", yamlConfig)))
			output.WriteString("\n\n")
		} else {
			output.WriteString(fmt.Sprintf("Error parsing YAML file %s: %v\n", file, err))
		}
	}
}

func analyzeDocumentation(output *strings.Builder) {
	output.WriteString("\n### Documentation Analysis ###\n")

	if fileExists("README.md") {
		content, _ := ioutil.ReadFile("README.md")
		output.WriteString("README.md contents:\n")
		output.WriteString(truncateContent(string(content)))
		output.WriteString("\n\n")
	}

	if fileExists("LICENSE") {
		content, _ := ioutil.ReadFile("LICENSE")
		output.WriteString("LICENSE contents:\n")
		output.WriteString(truncateContent(string(content)))
		output.WriteString("\n\n")
	}

	// Check for other documentation files
	docFiles, _ := filepath.Glob("docs/*.md")
	for _, file := range docFiles {
		content, _ := ioutil.ReadFile(file)
		output.WriteString(fmt.Sprintf("Documentation file: %s\n", file))
		output.WriteString(truncateContent(string(content)))
		output.WriteString("\n\n")
	}
}

func analyzeContainerization(output *strings.Builder) {
	output.WriteString("\n### Containerization Analysis ###\n")

	if fileExists("Dockerfile") {
		content, _ := ioutil.ReadFile("Dockerfile")
		output.WriteString("Dockerfile found:\n")
		output.WriteString(truncateContent(string(content)))
		output.WriteString("\n\n")
	}

	composeFiles := []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"}
	for _, file := range composeFiles {
		if fileExists(file) {
			content, _ := ioutil.ReadFile(file)
			output.WriteString(fmt.Sprintf("%s found:\n", file))
			output.WriteString(truncateContent(string(content)))
			output.WriteString("\n\n")
			break // We only need to find one compose file
		}
	}
}

func analyzeInfrastructureAsCode(output *strings.Builder) {
	output.WriteString("\n### Infrastructure as Code Analysis ###\n")

	if dirExists("terraform") {
		output.WriteString("Terraform configuration found.\n")
		tfFiles, _ := filepath.Glob("terraform/*.tf")
		for _, file := range tfFiles {
			content, _ := ioutil.ReadFile(file)
			output.WriteString(fmt.Sprintf("File: %s\n", filepath.Base(file)))
			output.WriteString(truncateContent(string(content)))
			output.WriteString("\n\n")
		}
	}

	if fileExists("serverless.yml") || fileExists("serverless.yaml") {
		output.WriteString("Serverless Framework configuration found.\n")
		// Add analysis of serverless config here
	}

	if fileExists("Chart.yaml") {
		output.WriteString("Helm Chart found.\n")
		// Add analysis of Helm Chart here
	}
}

func analyzeCICDPipelines(output *strings.Builder) {
	output.WriteString("\n### CI/CD Pipeline Analysis ###\n")

	if fileExists("Jenkinsfile") {
		content, _ := ioutil.ReadFile("Jenkinsfile")
		output.WriteString("Jenkinsfile found:\n")
		output.WriteString(truncateContent(string(content)))
		output.WriteString("\n\n")
	}

	if fileExists("cloudbuild.yaml") || fileExists("cloudbuild.yml") {
		output.WriteString("Google Cloud Build configuration found.\n")
		// Add analysis of Cloud Build config here
	}

	// Add checks for other CI/CD configurations (GitLab CI, CircleCI, etc.)
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
		if fileExists(baseName + ext) {
			return true
		}
	}
	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	return err == nil && info.IsDir()
}
