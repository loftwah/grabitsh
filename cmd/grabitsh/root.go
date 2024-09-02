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
	collectConfigFiles(&outputBuffer)
	collectLargeFiles(&outputBuffer)
	collectFileTypeSummary(&outputBuffer)
	collectRecentlyModifiedFiles(&outputBuffer)
	collectProjectTypes(&outputBuffer)

	// Output results
	finalizeOutput(outputBuffer.String())
}

// Helper functions to collect information
func collectRepoStructure(buffer *bytes.Buffer) {
	buffer.WriteString("### Repository Structure ###\n")
	buffer.WriteString(runCommand("ls", "-lah"))
	if _, err := exec.LookPath("tree"); err == nil {
		buffer.WriteString(runCommand("tree", "-L", "2", "-a"))
	} else {
		buffer.WriteString("(Tree command not available)\n")
	}
}

func collectGitInfo(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Git Information ###\n")
	buffer.WriteString("Recent Commits:\n")
	buffer.WriteString(runCommand("git", "log", "--oneline", "-n", "5"))
	buffer.WriteString("\nBranches:\n")
	buffer.WriteString(runCommand("git", "branch", "-a"))
	buffer.WriteString("\nRemote Repositories:\n")
	buffer.WriteString(runCommand("git", "remote", "-v"))
	buffer.WriteString("\nGit Status:\n")
	buffer.WriteString(runCommand("git", "status", "--short"))
}

func collectConfigFiles(buffer *bytes.Buffer) {
	buffer.WriteString("\n### Configuration and Important Files ###\n")
	importantFiles := []string{
		".gitignore", "README*", "LICENSE*", "Dockerfile*", ".env*",
		"Makefile", "package.json", "go.mod", "requirements.txt",
		"Gemfile", "composer.json", "build.gradle", "pom.xml",
	}
	for _, file := range importantFiles {
		matches, err := filepath.Glob(filepath.Join(".", file))
		if err == nil && len(matches) > 0 {
			for _, match := range matches {
				buffer.WriteString(match + "\n")
			}
		}
	}
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
