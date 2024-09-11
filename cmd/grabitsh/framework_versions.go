package grabitsh

import (
	"io/ioutil"
	"regexp"
	"strings"
)

func extractFrameworkVersions() map[string]string {
	versions := make(map[string]string)

	checkFrameworkVersion := func(file, framework, regex string) {
		if fileExists(file) {
			content, _ := ioutil.ReadFile(file)
			re := regexp.MustCompile(regex)
			matches := re.FindStringSubmatch(string(content))
			if len(matches) > 1 {
				versions[framework] = matches[1]
			}
		}
	}

	checkFrameworkVersion("Gemfile.lock", "Rails", `rails \((\d+\.\d+\.\d+)\)`)
	checkFrameworkVersion("package.json", "React", `"react": "(\^|~)?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("package.json", "Vue.js", `"vue": "(\^|~)?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("package.json", "Angular", `"@angular/core": "(\^|~)?(\d+\.\d+\.\d+)"`)
	checkFrameworkVersion("requirements.txt", "Django", `Django==(\d+\.\d+\.\d+)`)
	checkFrameworkVersion("requirements.txt", "Flask", `Flask==(\d+\.\d+\.\d+)`)
	checkFrameworkVersion("pom.xml", "Spring Boot", `<spring-boot.version>(\d+\.\d+\.\d+)`)
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

	return versions
}
