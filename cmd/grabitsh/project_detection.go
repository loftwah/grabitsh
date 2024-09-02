package grabitsh

import (
	"os"
)

func DetectProjectTypes() []string {
	var projectTypes []string

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

	if fileExists("vue.config.js") {
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

	return projectTypes
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	return err == nil && info.IsDir()
}
