package grabitsh

import (
	"path/filepath"
)

func analyzeTestingFrameworks() []string {
	var frameworks []string

	testingFrameworks := map[string]string{
		"test":        "Go testing",
		"spec":        "RSpec (Ruby)",
		"test.js":     "JavaScript testing",
		"test.py":     "Python testing",
		"__tests__":   "Jest (JavaScript)",
		"pytest":      "pytest (Python)",
		"phpunit.xml": "PHPUnit",
		"junit":       "JUnit (Java)",
		"test_*.py":   "unittest (Python)",
		"*_test.go":   "Go testing",
		"*.spec.ts":   "Jasmine/Mocha (TypeScript)",
		"test_*.rb":   "Minitest (Ruby)",
	}

	for pattern, framework := range testingFrameworks {
		if files, _ := filepath.Glob("**/" + pattern); len(files) > 0 {
			frameworks = append(frameworks, framework)
		}
	}

	return frameworks
}
