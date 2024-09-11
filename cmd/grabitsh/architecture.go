package grabitsh

func detectArchitecture() string {
	if dirExists("services") || dirExists("microservices") {
		return "Microservices"
	} else if fileExists("serverless.yml") || fileExists("serverless.yaml") {
		return "Serverless"
	} else if dirExists("app") && dirExists("config") && dirExists("db") {
		return "Monolithic (Rails-like)"
	} else if fileExists("package.json") && fileExists("server.js") {
		return "Monolithic (Node.js)"
	} else if fileExists("pom.xml") || fileExists("build.gradle") {
		return "Monolithic (Java)"
	} else if fileExists("manage.py") && dirExists("apps") {
		return "Monolithic (Django)"
	}

	return "Undetermined"
}
