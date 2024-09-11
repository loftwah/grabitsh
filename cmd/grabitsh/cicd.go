package grabitsh

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type CICDSystem struct {
	Name  string   `json:"name"`
	File  string   `json:"file"`
	Steps []string `json:"steps"`
}

func analyzeCICDWorkflows() []CICDSystem {
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
			files, _ := filepath.Glob(filePattern)
			for _, file := range files {
				content, err := ioutil.ReadFile(file)
				if err != nil {
					continue
				}
				steps := analyzeCICDSteps(string(content))
				results = append(results, CICDSystem{
					Name:  system.name,
					File:  filepath.Base(file),
					Steps: steps,
				})
			}
		}
	}

	return results
}

func analyzeCICDSteps(content string) []string {
	var steps []string

	if strings.Contains(content, "npm test") || strings.Contains(content, "yarn test") {
		steps = append(steps, "Testing")
	}
	if strings.Contains(content, "npm run build") || strings.Contains(content, "yarn build") {
		steps = append(steps, "Build")
	}
	if strings.Contains(content, "docker build") || strings.Contains(content, "docker-compose") {
		steps = append(steps, "Docker operations")
	}
	if strings.Contains(content, "deploy") || strings.Contains(content, "kubectl") {
		steps = append(steps, "Deployment")
	}
	if strings.Contains(content, "lint") || strings.Contains(content, "eslint") {
		steps = append(steps, "Linting")
	}
	if strings.Contains(content, "security") || strings.Contains(content, "scan") {
		steps = append(steps, "Security scanning")
	}
	if strings.Contains(content, "coverage") || strings.Contains(content, "codecov") {
		steps = append(steps, "Code coverage")
	}

	return steps
}
