package grabitsh

import (
	"bytes"
	"fmt"
	"os"
)

// Detect and collect important configuration files
func DetectImportantFiles(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Important Configuration Files ###\n")

	// Helper function to check if file exists and then parse
	checkAndParseIfExists := func(filename string, parser func(string, *bytes.Buffer), buffer *bytes.Buffer) {
		if _, err := os.Stat(filename); err == nil {
			parser(filename, buffer)
		}
	}

	// Helper function for parsers that return an error
	checkAndParseIfExistsWithError := func(filename string, parser func(string, *bytes.Buffer) error, buffer *bytes.Buffer) {
		if _, err := os.Stat(filename); err == nil {
			if err := parser(filename, buffer); err != nil {
				buffer.WriteString(fmt.Sprintf("Error parsing %s: %v\n", filename, err))
			}
		}
	}

	// 1. Version Control
	checkAndParseIfExists(".git/config", parseGitConfig, buffer)
	checkAndParseIfExists(".gitattributes", parseBasicTextFile, buffer)
	checkAndParseIfExists(".gitmodules", parseBasicTextFile, buffer)
	checkAndParseIfExists(".gitmessage", parseBasicTextFile, buffer)
	checkAndParseIfExists(".gitflow", parseBasicTextFile, buffer)

	// 2. Project Configuration
	checkAndParseIfExists(".editorconfig", parseBasicTextFile, buffer)
	checkAndParseIfExists(".vscode/settings.json", parseJSONFile, buffer)
	checkAndParseIfExists(".eslintignore", parseBasicTextFile, buffer)
	checkAndParseIfExists(".npmrc", parseBasicTextFile, buffer)
	checkAndParseIfExists(".nvmrc", parseBasicTextFile, buffer)
	checkAndParseIfExists(".yarnrc", parseBasicTextFile, buffer)
	checkAndParseIfExists("lerna.json", parseJSONFile, buffer)
	checkAndParseIfExists("nx.json", parseJSONFile, buffer)

	// 3. CI/CD and DevOps
	checkAndParseIfExists(".travis.yml", parseYAMLFile, buffer)
	checkAndParseIfExists(".circleci/config.yml", parseYAMLFile, buffer)
	checkAndParseIfExists(".gitlab-ci.yml", parseYAMLFile, buffer)
	checkAndParseIfExists(".github/workflows", parseGithubActionsWorkflows, buffer)
	checkAndParseIfExists("Jenkinsfile", parseBasicTextFile, buffer)
	checkAndParseIfExists("azure-pipelines.yml", parseYAMLFile, buffer)
	checkAndParseIfExists("bitbucket-pipelines.yml", parseYAMLFile, buffer)
	checkAndParseIfExists("sonar-project.properties", parseBasicTextFile, buffer)
	checkAndParseIfExists("codecov.yml", parseYAMLFile, buffer)
	checkAndParseIfExists(".snyk", parseBasicTextFile, buffer)

	// 4. Docker and Containerization
	checkAndParseIfExists("docker-compose.yml", parseYAMLFile, buffer)
	checkAndParseIfExists("docker-compose.override.yml", parseYAMLFile, buffer)
	checkAndParseIfExists("Dockerfile", parseDockerfile, buffer)
	checkAndParseIfExists(".dockerignore", parseBasicTextFile, buffer)
	checkAndParseIfExists("docker/", parseDockerDir, buffer)

	// 5. Kubernetes
	checkAndParseIfExists("k8s/", parseK8sFiles, buffer)
	checkAndParseIfExists("helm/", parseHelmFiles, buffer)
	checkAndParseIfExists("values.yaml", parseYAMLFile, buffer)
	checkAndParseIfExists("kustomization.yaml", parseYAMLFile, buffer)

	// 6. Cloud Providers
	checkAndParseIfExists("serverless.yml", parseYAMLFile, buffer)
	checkAndParseIfExists(".aws/", parseDirectoryContents, buffer)
	checkAndParseIfExists("firebase.json", parseJSONFile, buffer)
	checkAndParseIfExists("vercel.json", parseJSONFile, buffer)
	checkAndParseIfExists("netlify.toml", parseYAMLFile, buffer)

	// 7. Infrastructure as Code
	checkAndParseIfExists("main.tf", parseBasicTextFile, buffer)
	checkAndParseIfExists("Vagrantfile", parseBasicTextFile, buffer)
	checkAndParseIfExists("Pulumi.yaml", parseYAMLFile, buffer)

	// 8. Language-Specific
	checkAndParseIfExistsWithError("Gemfile", parseGemfile, buffer)
	checkAndParseIfExists(".ruby-version", parseBasicTextFile, buffer)
	checkAndParseIfExists("config/application.rb", parseBasicTextFile, buffer)

	// Python
	checkAndParseIfExists("requirements.txt", parseBasicTextFile, buffer)
	checkAndParseIfExists("setup.py", parseBasicTextFile, buffer)

	// JavaScript / TypeScript
	checkAndParseIfExistsWithError("package.json", parsePackageJSON, buffer)
	checkAndParseIfExists("tsconfig.json", parseJSONFile, buffer)

	// Go
	checkAndParseIfExists("go.mod", parseBasicTextFile, buffer)
	checkAndParseIfExists("go.sum", parseBasicTextFile, buffer)

	// PHP
	checkAndParseIfExists("composer.json", parseJSONFile, buffer)
	checkAndParseIfExists("phpunit.xml", parseBasicTextFile, buffer)

	// Java / Kotlin
	checkAndParseIfExists("pom.xml", parseBasicTextFile, buffer)
	checkAndParseIfExists("build.gradle", parseBasicTextFile, buffer)

	// 9. Web Frameworks
	checkAndParseIfExists("next.config.js", parseBasicTextFile, buffer)
	checkAndParseIfExists("nuxt.config.js", parseBasicTextFile, buffer)

	// 10. Miscellaneous
	checkAndParseIfExists("CHANGELOG.md", parseBasicTextFile, buffer)
	checkAndParseIfExists("CONTRIBUTING.md", parseBasicTextFile, buffer)
	checkAndParseIfExists("CODE_OF_CONDUCT.md", parseBasicTextFile, buffer)
}
