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

	// Repository structure
	outputBuffer.WriteString("### Repository Structure ###\n")
	outputBuffer.WriteString(runCommand("ls", "-lah"))
	if _, err := exec.LookPath("tree"); err == nil {
		outputBuffer.WriteString(runCommand("tree", "-L", "2", "-a"))
	} else {
		outputBuffer.WriteString("(Tree command not available)\n")
	}

	// Git information
	outputBuffer.WriteString("\n### Git Information ###\n")
	outputBuffer.WriteString("Recent Commits:\n")
	outputBuffer.WriteString(runCommand("git", "log", "--oneline", "-n", "5"))
	outputBuffer.WriteString("\nBranches:\n")
	outputBuffer.WriteString(runCommand("git", "branch", "-a"))
	outputBuffer.WriteString("\nRemote Repositories:\n")
	outputBuffer.WriteString(runCommand("git", "remote", "-v"))
	outputBuffer.WriteString("\nGit Status:\n")
	outputBuffer.WriteString(runCommand("git", "status", "--short"))

	// Configuration and important files
	outputBuffer.WriteString("\n### Configuration and Important Files ###\n")
	importantFiles := []string{
		".gitignore", "README*", "LICENSE*", "Dockerfile*", ".env*",
		"Makefile", "package.json", "go.mod", "requirements.txt",
		"Gemfile", "composer.json", "build.gradle", "pom.xml",
	}
	for _, file := range importantFiles {
		matches, err := filepath.Glob(filepath.Join(".", file))
		if err == nil && len(matches) > 0 {
			for _, match := range matches {
				outputBuffer.WriteString(match + "\n")
			}
		}
	}

	// Large files
	outputBuffer.WriteString("\n### Large Files (top 5) ###\n")
	outputBuffer.WriteString(findLargeFiles())

	// File types summary
	outputBuffer.WriteString("\n### File Types Summary ###\n")
	outputBuffer.WriteString(summarizeFileTypes())

	// Recent changes
	outputBuffer.WriteString("\n### Recently Modified Files ###\n")
	outputBuffer.WriteString(findRecentlyModifiedFiles())

	// Project type detection
	outputBuffer.WriteString("\n### Project Type Detection ###\n")
	projectTypes := DetectProjectTypes()
	for _, projectType := range projectTypes {
		outputBuffer.WriteString(fmt.Sprintf("- %s\n", projectType))
	}

	// Output results
	finalizeOutput(outputBuffer.String())
}

func runCommand(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running command %s %s: %v\n", name, strings.Join(arg, " "), err)
	}
	return string(out)
}

func findLargeFiles() string {
	cmd := exec.Command("bash", "-c", "find . -type f -exec du -h {} + | sort -rh | head -n 5")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error finding large files: %v\n", err)
	}
	return string(out)
}

func summarizeFileTypes() string {
	cmd := exec.Command("bash", "-c", "find . -type f | sed -e 's/.*\\.//' | sort | uniq -c | sort -rn | head -n 10")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error summarizing file types: %v\n", err)
	}
	return string(out)
}

func findRecentlyModifiedFiles() string {
	cmd := exec.Command("find", ".", "-type", "f", "-mtime", "-7", "-not", "-path", "./.git/*")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error finding recently modified files: %v\n", err)
	}
	return string(out)
}

func finalizeOutput(content string) {
	switch outputMethod {
	case "stdout":
		color.Green(content)
	case "clipboard":
		err := clipboard.WriteAll(content)
		if err != nil {
			color.Red("Failed to copy to clipboard: %v", err)
		} else {
			color.Green("Copied to clipboard successfully.")
		}
	case "file":
		if outputFile == "" {
			color.Red("Output file path must be specified when using file output method.")
			return
		}
		err := os.WriteFile(outputFile, []byte(content), 0644)
		if err != nil {
			color.Red("Failed to write to file: %v", err)
		} else {
			color.Green("Output written to file: %s", outputFile)
		}
	default:
		color.Red("Invalid output method. Choose stdout, clipboard, or file.")
	}
}
