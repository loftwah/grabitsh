package grabitsh

import (
	"os"
	"regexp"
	"strings"
)

func extractFrameworkVersions() map[string]string {
	versions := make(map[string]string)

	// Check for versions of various frameworks using specific files and regex patterns.
	checkFrameworkVersion := func(file, framework, regex string) {
		if fileExists(file) {
			content, _ := os.ReadFile(file)
			re := regexp.MustCompile(regex)
			matches := re.FindStringSubmatch(string(content))
			if len(matches) > 1 {
				versions[framework] = matches[1]
			}
		}
	}

	// Ruby/Rails framework
	checkFrameworkVersion("Gemfile.lock", "Rails", `rails \((\d+\.\d+\.\d+)\)`)

	// JavaScript/Node.js frameworks
	checkFrameworkVersion("package.json", "React", `"react": "(\^|~)?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("package.json", "Vue.js", `"vue": "(\^|~)?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("package.json", "Angular", `"@angular/core": "(\^|~)?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("package.json", "Express", `"express": "(\^|~)?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("package.json", "Next.js", `"next": "(\^|~)?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("package.json", "Svelte", `"svelte": "(\^|~)?(\d+\.\d+\.\d+)"`)

	// Python frameworks
	checkFrameworkVersion("requirements.txt", "Django", `Django==(\d+\.\d+\.\d+)`)
	checkFrameworkVersion("requirements.txt", "Flask", `Flask==(\d+\.\d+\.\d+)`)
	checkFrameworkVersion("requirements.txt", "FastAPI", `fastapi==(\d+\.\d+\.\d+)`)

	// PHP frameworks
	checkFrameworkVersion("composer.lock", "Laravel", `"name": "laravel/framework",\s*"version": "v?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("composer.lock", "Symfony", `"name": "symfony/symfony",\s*"version": "v?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("composer.lock", "WordPress", `"name": "wordpress/core",\s*"version": "v?(\d+\.\d+\.\d+)"`)

	// Go frameworks
	checkFrameworkVersion("go.mod", "Gin", `github.com/gin-gonic/gin\s*v(\d+\.\d+\.\d+)`)
	checkFrameworkVersion("go.mod", "Echo", `github.com/labstack/echo/v4\s*v(\d+\.\d+\.\d+)`)
	checkFrameworkVersion("go.mod", "Fiber", `github.com/gofiber/fiber/v2\s*v(\d+\.\d+\.\d+)`)

	// Rust frameworks
	checkFrameworkVersion("Cargo.toml", "Rocket", `rocket\s*=\s*"(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("Cargo.toml", "Actix", `actix-web\s*=\s*"(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("Cargo.toml", "Tide", `tide\s*=\s*"(\d+\.\d+\.\d+)"`)

	// Java frameworks
	checkFrameworkVersion("pom.xml", "Spring Boot", `<spring-boot.version>(\d+\.\d+\.\d+)</spring-boot.version>`)
	checkFrameworkVersion("build.gradle", "Spring Boot", `springBootVersion = '(\d+\.\d+\.\d+)'`)

	// Check for Node.js and npm versions
	if fileExists("package.json") {
		versions["Node.js"] = strings.TrimSpace(runCommand("node", "-v"))
		versions["npm"] = strings.TrimSpace(runCommand("npm", "-v"))
	}

	// Check for Python version
	versions["Python"] = strings.TrimSpace(runCommand("python", "--version"))

	// Check for Go version
	if fileExists("go.mod") {
		versions["Go"] = strings.TrimSpace(strings.TrimPrefix(runCommand("go", "version"), "go version "))
	}

	// Check for PHP version
	if fileExists("composer.lock") {
		versions["PHP"] = strings.TrimSpace(runCommand("php", "-v"))
	}

	// Check for Rust version
	if fileExists("Cargo.toml") {
		versions["Rust"] = strings.TrimSpace(runCommand("rustc", "--version"))
	}

	return versions
}
