package grabitsh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Step struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CICDSystem struct {
	Name  string `json:"name"`
	File  string `json:"file"`
	Steps []Step `json:"steps"`
}

func analyzeCICDWorkflows() ([]CICDSystem, error) {
	cicdSystems := []struct {
		name  string
		files []string
	}{
		{"GitHub Actions", []string{".github/workflows/*.yml", ".github/workflows/*.yaml"}},
		{"GitLab CI", []string{".gitlab-ci.yml"}},
		{"Jenkins", []string{"Jenkinsfile"}},
		{"CircleCI", []string{".circleci/config.yml"}},
		{"Travis CI", []string{".travis.yml"}},
		{"Azure Pipelines", []string{"azure-pipelines.yml"}},
		{"Bitbucket Pipelines", []string{"bitbucket-pipelines.yml"}},
		{"AWS CodeBuild", []string{"buildspec.yml"}},
		{"Drone CI", []string{".drone.yml"}},
		{"Semaphore", []string{".semaphore/semaphore.yml"}},
	}

	var results []CICDSystem

	for _, system := range cicdSystems {
		for _, filePattern := range system.files {
			files, err := filepath.Glob(filePattern)
			if err != nil {
				return nil, fmt.Errorf("error globbing files: %w", err)
			}
			for _, file := range files {
				content, err := os.ReadFile(file)
				if err != nil {
					return nil, fmt.Errorf("error reading file %s: %w", file, err)
				}
				steps, err := analyzeCICDSteps(string(content))
				if err != nil {
					return nil, fmt.Errorf("error analyzing steps in file %s: %w", file, err)
				}
				results = append(results, CICDSystem{
					Name:  system.name,
					File:  filepath.Base(file),
					Steps: steps,
				})
			}
		}
	}

	return results, nil
}

func analyzeCICDSteps(content string) ([]Step, error) {
	var steps []Step

	if strings.Contains(content, "npm test") || strings.Contains(content, "yarn test") {
		steps = append(steps, Step{Name: "Testing", Description: "Runs tests"})
	}
	if strings.Contains(content, "npm run build") || strings.Contains(content, "yarn build") {
		steps = append(steps, Step{Name: "Build", Description: "Builds the project"})
	}
	if strings.Contains(content, "docker build") || strings.Contains(content, "docker-compose") {
		steps = append(steps, Step{Name: "Docker operations", Description: "Performs Docker operations"})
	}
	if strings.Contains(content, "deploy") || strings.Contains(content, "kubectl") {
		steps = append(steps, Step{Name: "Deployment", Description: "Deploys the project"})
	}
	if strings.Contains(content, "lint") || strings.Contains(content, "eslint") {
		steps = append(steps, Step{Name: "Linting", Description: "Runs linter"})
	}
	if strings.Contains(content, "security") || strings.Contains(content, "scan") {
		steps = append(steps, Step{Name: "Security scanning", Description: "Performs security scanning"})
	}
	if strings.Contains(content, "coverage") || strings.Contains(content, "codecov") {
		steps = append(steps, Step{Name: "Code coverage", Description: "Checks code coverage"})
	}

	return steps, nil
}
