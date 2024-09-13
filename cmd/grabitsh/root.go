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
	PerformAdvancedAnalysis(&outputBuffer)

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

	// Improved exclusion: Exclude grabitsh_chunk files and root.go itself to avoid recursive results
	todoCommand := `grep -r -n --exclude-dir={.git,node_modules,vendor} --exclude=\*.min.js --exclude=\*.min.css --exclude=\*grabitsh_chunk_*.txt --exclude=root.go --binary-files=without-match "TODO\|FIXME" .`

	// Execute the command
	todos := runCommand("bash", "-c", todoCommand)

	// Handle cases where grep fails to find anything or errors out
	if strings.TrimSpace(todos) == "" {
		buffer.WriteString("No TODOs or FIXMEs found.\n")
	} else {
		buffer.WriteString("Found TODOs and FIXMEs:\n")
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

	// Define the analysis template function
	getAnalysisTemplate := func(part, total, estimatedTokens int, grabitOutput string) string {
		return fmt.Sprintf(`You are an AI assistant tasked with analyzing and extracting valuable information from the output of a tool called grabit.sh. This tool scans a newly cloned git repository and gathers as much information as possible. Your goal is to process this output and provide a comprehensive summary of the repository's characteristics and contents.

This analysis may be split into multiple chunks due to token limitations. The current chunk number and total number of chunks are provided to you. Please behave accordingly:

- If this is chunk 1 of 1, provide a complete analysis.
- If this is chunk 1 of multiple chunks, begin your analysis but indicate that it's incomplete.
- If this is a middle chunk, continue the analysis without repeating information from previous chunks.
- If this is the final chunk of multiple chunks, complete the analysis and provide a comprehensive summary.

Here is the output from grabit.sh for chunk %d of %d:

<grabit_output>
%s
</grabit_output>

Analyze the provided output and extract all relevant and useful information about the repository. Focus on the following aspects:

1. Repository structure and organization
2. Main programming languages used
3. Key files and their purposes
4. Dependencies and external libraries
5. Build system or package management tools
6. Testing framework and test coverage
7. Documentation availability and quality
8. Coding standards and style guidelines
9. Version control information (branches, tags, commit history)
10. Licensing information
11. Any unique or noteworthy features of the project

As you analyze the output, think critically about what information would be most valuable to a developer or project manager trying to understand this repository quickly. Look for patterns, unusual elements, or potential areas of concern.

Organize your findings into a clear and concise summary. Use bullet points or short paragraphs to present the information in an easily digestible format. If you encounter any ambiguities or areas that require further investigation, mention these as well.

Present your analysis and summary within <analysis> tags. If this is not the final chunk, end your analysis with a note indicating that the summary is incomplete and will be continued in the next chunk.

Remember to focus on extracting and presenting the most relevant and useful information from the grabit.sh output, providing valuable insights into the repository's characteristics and contents.

**Chunk %d/%d (Estimated %d tokens):**

`, part, total, grabitOutput, part, total, estimatedTokens)
	}

	// Define the final chunk message
	getFinalChunkMessage := func(part, total int) string {
		if part < total {
			return fmt.Sprintf("\n**Continue to next chunk (Chunk %d/%d)**\n", part+1, total)
		}
		return "\n**Final chunk—no more parts to follow.**\n"
	}

	for i, chunk := range chunks {
		// Estimate the number of tokens
		estimatedTokens := len(strings.Fields(chunk)) + len(chunk)/3

		// Build the full content using the analysis template
		fullContent := getAnalysisTemplate(i+1, totalChunks, estimatedTokens, chunk)
		fullContent += getFinalChunkMessage(i+1, totalChunks)

		// Write the chunk to file
		filename := fmt.Sprintf("grabitsh_chunk_%d.txt", i+1)
		if err := os.WriteFile(filename, []byte(fullContent), 0644); err != nil {
			return fmt.Errorf("failed to write chunk %d: %v", i+1, err)
		}

		// Output success message
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
