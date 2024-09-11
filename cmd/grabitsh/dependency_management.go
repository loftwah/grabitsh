package grabitsh

func analyzeDependencyManagement() []string {
	var tools []string

	depManagement := map[string]string{
		"package-lock.json": "npm",
		"yarn.lock":         "Yarn",
		"Gemfile.lock":      "Bundler (Ruby)",
		"poetry.lock":       "Poetry (Python)",
		"go.sum":            "Go Modules",
		"composer.lock":     "Composer (PHP)",
		"Pipfile.lock":      "Pipenv (Python)",
		"pom.xml":           "Maven (Java)",
		"build.gradle":      "Gradle (Java)",
		"requirements.txt":  "pip (Python)",
		"Cargo.lock":        "Cargo (Rust)",
	}

	for file, tool := range depManagement {
		if fileExists(file) {
			tools = append(tools, tool)
		}
	}

	return tools
}
