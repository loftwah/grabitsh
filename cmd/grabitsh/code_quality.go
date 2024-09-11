package grabitsh

func analyzeCodeQuality() []string {
	var tools []string

	lintConfigs := map[string]string{
		".eslintrc":      "ESLint",
		".rubocop.yml":   "RuboCop",
		".golangci.yml":  "golangci-lint",
		"pylintrc":       "Pylint",
		".checkstyle":    "Checkstyle (Java)",
		"tslint.json":    "TSLint",
		".stylelintrc":   "Stylelint",
		".prettierrc":    "Prettier",
		".scalafmt.conf": "Scalafmt",
	}

	for config, tool := range lintConfigs {
		if fileExists(config) {
			tools = append(tools, tool)
		}
	}

	if fileExists("sonar-project.properties") {
		tools = append(tools, "SonarQube")
	}

	if fileExists(".codeclimate.yml") {
		tools = append(tools, "CodeClimate")
	}

	return tools
}
