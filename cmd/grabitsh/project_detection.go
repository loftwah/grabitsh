package grabitsh

import (
	"os"
)

func DetectProjectTypes() []string {
	var projectTypes []string

	// Node.js / JavaScript
	if fileExists("package.json") {
		projectTypes = append(projectTypes, "Node.js project")
		if fileExists("next.config.js") {
			projectTypes = append(projectTypes, "Next.js framework")
		}
		if fileExists("react-scripts.config.js") || (dirExists("src") && fileExists("src/App.js")) {
			projectTypes = append(projectTypes, "React project")
		}
		if fileExists("astro.config.mjs") {
			projectTypes = append(projectTypes, "Astro framework")
		}
	}

	// Ruby on Rails
	if fileExists("config/application.rb") && dirExists("app") && dirExists("config") {
		projectTypes = append(projectTypes, "Ruby on Rails project")
	}

	// Laravel (PHP)
	if fileExists("artisan") && dirExists("app") && dirExists("public") {
		projectTypes = append(projectTypes, "Laravel (PHP) project")
	}

	// Django (Python)
	if fileExists("manage.py") && dirExists("templates") {
		projectTypes = append(projectTypes, "Django (Python) project")
	}

	// Flask (Python)
	if fileExists("app.py") || fileExists("wsgi.py") {
		projectTypes = append(projectTypes, "Flask (Python) project")
	}

	// Vue.js
	if fileExists("vue.config.js") {
		projectTypes = append(projectTypes, "Vue.js project")
	}

	// Angular
	if fileExists("angular.json") {
		projectTypes = append(projectTypes, "Angular project")
	}

	// .NET Core
	if fileExists("Program.cs") && dirExists("bin") && dirExists("obj") {
		projectTypes = append(projectTypes, ".NET Core project")
	}

	// Java Spring Boot
	if fileExists("pom.xml") && dirExists("src/main/java") {
		projectTypes = append(projectTypes, "Java Spring Boot project")
	}

	// Go
	if fileExists("go.mod") {
		projectTypes = append(projectTypes, "Go project")
	}

	return projectTypes
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
