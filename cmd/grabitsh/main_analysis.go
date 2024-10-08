package grabitsh

import (
	"bytes"
	"encoding/json"
	"sync"
)

type AnalysisResult struct {
	Architecture         string            `json:"architecture"`
	FrameworkVersions    map[string]string `json:"framework_versions"`
	CICDSystems          []CICDSystem      `json:"cicd_systems"`
	APIStructure         APIInfo           `json:"api_structure"`
	DatabaseUsage        DatabaseInfo      `json:"database_usage"`
	TestingFrameworks    []string          `json:"testing_frameworks"`
	CodeQualityTools     []string          `json:"code_quality_tools"`
	DependencyManagement []string          `json:"dependency_management"`
}

func PerformAdvancedAnalysis(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Advanced Analysis ###\n")

	var result AnalysisResult
	var wg sync.WaitGroup
	wg.Add(7)

	go func() {
		defer wg.Done()
		result.Architecture = detectArchitecture()
	}()

	go func() {
		defer wg.Done()
		result.FrameworkVersions = extractFrameworkVersions()
	}()

	// Fix for capturing two return values
	go func() {
		defer wg.Done()
		cicdSystems, err := analyzeCICDWorkflows() // Capture both values
		if err != nil {                            // Handle the error
			buffer.WriteString("Error analyzing CI/CD workflows: " + err.Error() + "\n")
			return
		}
		result.CICDSystems = cicdSystems // Assign result if no error
	}()

	go func() {
		defer wg.Done()
		result.APIStructure = analyzeAPIStructure()
	}()

	go func() {
		defer wg.Done()
		result.DatabaseUsage = analyzeDatabaseUsage()
	}()

	go func() {
		defer wg.Done()
		result.TestingFrameworks = analyzeTestingFrameworks()
	}()

	go func() {
		defer wg.Done()
		result.CodeQualityTools = analyzeCodeQuality()
		result.DependencyManagement = analyzeDependencyManagement()
	}()

	wg.Wait()

	// Marshal the result to JSON and write to buffer
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	buffer.WriteString(string(jsonResult))
}
