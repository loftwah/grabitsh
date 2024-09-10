package grabitsh

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	outputMethod string
	outputFile   string
	chunkSize    int
	rootCmd      *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "grabitsh",
		Short: "Grabit.sh gathers useful information from a Git repository",
		Long:  `Grabit.sh simplifies working with Git repositories by gathering useful information and outputting it to stdout, a file, the clipboard, or LLM-friendly chunks.`,
		Run:   runGrabit,
	}

	rootCmd.Flags().StringVarP(&outputMethod, "output", "o", "stdout", "Output method: stdout, clipboard, file, or llm-chunks")
	rootCmd.Flags().StringVarP(&outputFile, "file", "f", "", "Output file path (required if output method is file)")
	rootCmd.Flags().IntVarP(&chunkSize, "chunk-size", "c", 100000, "Token size for LLM chunks (default 100000)")

	rootCmd.AddCommand(serveCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func runGrabit(cmd *cobra.Command, args []string) {
	var outputBuffer bytes.Buffer

	// Collect all sections
	collectRepoStructure(&outputBuffer)
	collectGitInfo(&outputBuffer)
	analyzeGitDir(&outputBuffer)
	collectProjectAnalysis(&outputBuffer)
	collectLargeFiles(&outputBuffer)
	collectFileTypeSummary(&outputBuffer)
	collectRecentlyModifiedFiles(&outputBuffer)
	collectProjectTypes(&outputBuffer)
	collectTODOs(&outputBuffer)
	collectSecurityAnalysis(&outputBuffer)
	collectPerformanceMetrics(&outputBuffer)
	DetectImportantFiles(&outputBuffer)

	// Output results
	finalizeOutput(outputBuffer.String())
}

func collectRepoStructure(buffer *bytes.Buffer) {
	buffer.WriteString("### Repository Structure ###\n")

	excludeDirs := []string{"node_modules", ".git/objects", ".git/logs", ".git/packs"}

	// Use the tree command or ls based on availability and exclude the directories
	if _, err := exec.LookPath("tree"); err == nil {
		buffer.WriteString(runCommand("tree", "-L", "3", "-a", "--prune", "-I", strings.Join(excludeDirs, "|")))
	} else {
		buffer.WriteString(runCommand("ls", "-lah"))
		buffer.WriteString("(Tree command not available)\n")
	}
}

func collectGitInfo(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Git Information ###\n")
	buffer.WriteString("Recent Commits:\n")
	buffer.WriteString(runCommand("git", "log", "--oneline", "-n", "10"))
	buffer.WriteString("\nBranches:\n")
	buffer.WriteString(runCommand("git", "branch", "-a"))
	buffer.WriteString("\nRemote Repositories:\n")
	buffer.WriteString(runCommand("git", "remote", "-v"))
	buffer.WriteString("\nGit Status:\n")
	buffer.WriteString(runCommand("git", "status", "--short"))
}

func collectProjectAnalysis(buffer *bytes.Buffer) {
	buffer.WriteString("\n")
	buffer.WriteString(AnalyzeRepository()) // Calls analysis from project_detection.go
}

func collectLargeFiles(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Large Files (top 5) ###\n")
	buffer.WriteString(runCommand("bash", "-c", "find . -type f -exec du -h {} + | sort -rh | head -n 5"))
}

func collectFileTypeSummary(buffer *bytes.Buffer) {
	buffer.WriteString("\n### File Types Summary ###\n")
	buffer.WriteString(runCommand("bash", "-c", "find . -type f | sed -e 's/.*\\.//' | sort | uniq -c | sort -rn | head -n 10"))
}

func collectRecentlyModifiedFiles(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Recently Modified Files ###\n")
	buffer.WriteString(runCommand("find", ".", "-type", "f", "-mtime", "-7", "-not", "-path", "./.git/*"))
}

func collectProjectTypes(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Project Type Detection ###\n")
	projectTypes := DetectProjectTypes()
	for _, projectType := range projectTypes {
		buffer.WriteString(fmt.Sprintf("- %s\n", projectType))
	}
}

func collectTODOs(buffer *bytes.Buffer) {
	buffer.WriteString("\n### TODOs and FIXMEs ###\n")
	todoCommand := `grep -r -n --binary-files=without-match "TODO\|FIXME" --exclude-dir={.git,node_modules,vendor} .`
	todos := runCommand("bash", "-c", todoCommand)
	if todos == "" {
		buffer.WriteString("No TODOs or FIXMEs found.\n")
	} else {
		buffer.WriteString(todos)
	}
}

func collectSecurityAnalysis(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Security Analysis ###\n")

	// Check for sensitive files
	sensitiveFiles := []string{".env", "id_rsa", "id_dsa", "*.pem", "*.key"}
	for _, pattern := range sensitiveFiles {
		files, _ := filepath.Glob(pattern)
		if len(files) > 0 {
			for _, file := range files {
				if file == ".env" {
					// Output .env file as sanitized example
					content, _ := os.ReadFile(file)
					buffer.WriteString(fmt.Sprintf("Sanitized .env Example:\n%s\n", sanitizeEnvFile(string(content))))
				} else {
					buffer.WriteString(fmt.Sprintf("Warning: Sensitive file detected: %s\n", file))
				}
			}
		}
	}

	// Check for outdated dependencies (example for Node.js projects)
	if fileExists("package.json") {
		buffer.WriteString(runCommand("npm", "audit"))
	}
}

func collectPerformanceMetrics(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Performance Metrics ###\n")

	// Repository size
	buffer.WriteString("Repository size:\n")
	buffer.WriteString(runCommand("du", "-sh", "."))

	// Number of files
	buffer.WriteString("\nTotal number of files:\n")
	buffer.WriteString(runCommand("bash", "-c", "find . -type f | wc -l"))

	// Lines of code (excluding .git directory)
	buffer.WriteString("\nTotal lines of code:\n")
	buffer.WriteString(runCommand("bash", "-c", "find . -name '*.go' -not -path './.git/*' | xargs wc -l"))
}

func runCommand(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running command %s %s: %v\n", name, strings.Join(arg, " "), err)
	}
	return string(out)
}

func finalizeOutput(content string) {
	switch outputMethod {
	case "stdout":
		color.Green(content)
	case "clipboard":
		if err := clipboard.WriteAll(content); err != nil {
			color.Red("Failed to copy to clipboard: %v", err)
		} else {
			color.Green("Copied to clipboard successfully.")
		}
	case "file":
		if outputFile == "" {
			color.Red("Output file path must be specified when using file output method.")
			return
		}
		if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
			color.Red("Failed to write to file: %v", err)
		} else {
			color.Green("Output written to file: %s", outputFile)
		}
	case "llm-chunks":
		if err := writeChunks(content); err != nil {
			color.Red("Failed to write LLM chunks: %v", err)
		} else {
			color.Green("LLM chunks written successfully.")
		}
	default:
		color.Red("Invalid output method. Choose stdout, clipboard, file, or llm-chunks.")
	}
}

func writeChunks(content string) error {
	chunks := splitIntoChunks(content, chunkSize)
	totalChunks := len(chunks)

	preamble := `This is part %d of %d of the output from Grabit.sh, a tool that analyzes Git repositories.

Purpose: This output provides a comprehensive analysis of a Git repository, including its structure, configuration, dependencies, and potential issues. Use this information to understand the project, identify areas for improvement, and make informed decisions about the codebase.

Instructions:
1. Read through the information provided in this chunk.
2. If this is not the final chunk, wait for the next one before drawing conclusions.
3. Use the information to answer questions about the repository, suggest improvements, or identify potential issues.
4. Pay attention to sections like security analysis, performance metrics, and detected project types.

Content of Chunk %d/%d (Estimated %d tokens):

`

	for i, chunk := range chunks {
		filename := fmt.Sprintf("grabitsh_chunk_%d.txt", i+1)
		estimatedTokens := len(strings.Fields(chunk)) + len(chunk)/3 // Same estimation as in splitIntoChunks
		fullContent := fmt.Sprintf(preamble, i+1, totalChunks, i+1, totalChunks, estimatedTokens) + chunk
		if err := os.WriteFile(filename, []byte(fullContent), 0644); err != nil {
			return fmt.Errorf("failed to write chunk %d: %v", i+1, err)
		}
		color.Green("Chunk %d/%d written to %s (Estimated %d tokens)", i+1, totalChunks, filename, estimatedTokens)
	}
	return nil
}

func splitIntoChunks(content string, chunkSize int) []string {
	var chunks []string
	lines := strings.Split(content, "\n")
	currentChunk := ""
	tokenCount := 0
	preambleSize := 250 // Approximate size of the preamble in tokens

	estimateTokens := func(s string) int {
		// This is a rough estimation. Actual tokenization varies by model.
		return len(strings.Fields(s)) + len(s)/3
	}

	for _, line := range lines {
		lineTokens := estimateTokens(line)
		if tokenCount+lineTokens > chunkSize-preambleSize {
			if currentChunk != "" {
				chunks = append(chunks, strings.TrimSpace(currentChunk))
				currentChunk = ""
				tokenCount = 0
			}
		}
		currentChunk += line + "\n"
		tokenCount += lineTokens
	}

	if currentChunk != "" {
		chunks = append(chunks, strings.TrimSpace(currentChunk))
	}

	return chunks
}
