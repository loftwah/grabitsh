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
	rootCmd      *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "grabitsh",
		Short: "Grabit.sh gathers useful information from a Git repository",
		Long:  `Grabit.sh simplifies working with Git repositories by gathering useful information and outputting it to stdout, a file, or the clipboard.`,
		Run:   runGrabit,
	}

	rootCmd.Flags().StringVarP(&outputMethod, "output", "o", "stdout", "Output method: stdout, clipboard, or file")
	rootCmd.Flags().StringVarP(&outputFile, "file", "f", "", "Output file path (required if output method is file)")

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
	collectProjectAnalysis(&outputBuffer)
	collectLargeFiles(&outputBuffer)
	collectFileTypeSummary(&outputBuffer)
	collectRecentlyModifiedFiles(&outputBuffer)
	collectProjectTypes(&outputBuffer)
	collectTODOs(&outputBuffer)
	collectSecurityAnalysis(&outputBuffer)
	collectPerformanceMetrics(&outputBuffer)

	// Output results
	finalizeOutput(outputBuffer.String())
}

func collectRepoStructure(buffer *bytes.Buffer) {
	buffer.WriteString("### Repository Structure ###\n")
	buffer.WriteString(runCommand("ls", "-lah"))
	if _, err := exec.LookPath("tree"); err == nil {
		buffer.WriteString(runCommand("tree", "-L", "3", "-a"))
	} else {
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
	buffer.WriteString(AnalyzeRepository())
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
	todoCommand := `grep -r -n "TODO\|FIXME" --exclude-dir={.git,node_modules,vendor} .`
	todos := runCommand("bash", "-c", todoCommand)
	if todos != "" {
		buffer.WriteString(todos)
	} else {
		buffer.WriteString("No TODOs or FIXMEs found.\n")
	}
}

func collectSecurityAnalysis(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Security Analysis ###\n")

	// Check for sensitive files
	sensitiveFiles := []string{".env", "id_rsa", "id_dsa", "*.pem", "*.key"}
	for _, pattern := range sensitiveFiles {
		files, _ := filepath.Glob(pattern)
		if len(files) > 0 {
			buffer.WriteString(fmt.Sprintf("Warning: Potentially sensitive files found: %v\n", files))
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
	default:
		color.Red("Invalid output method. Choose stdout, clipboard, or file.")
	}
}
