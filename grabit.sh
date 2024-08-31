#!/bin/bash

# Script to analyze project repository and provide insights.

# Check for required dependencies
check_dependencies() {
  local dependencies=("git" "xclip" "jq" "yq")

  for dep in "${dependencies[@]}"; do
    if ! command -v $dep &> /dev/null; then
      echo "Dependency $dep is missing. Please install it:"
      case $dep in
        git)
          echo "sudo apt-get install git -y"
          ;;
        xclip)
          echo "sudo apt-get install xclip -y"
          ;;
        jq)
          echo "sudo apt-get install jq -y"
          ;;
        yq)
          echo "sudo snap install yq"
          ;;
      esac
    fi
  done
}

# Analyze project structure
analyze_structure() {
  echo "### Project Structure ###"
  tree -a -I '.git' . 
}

# Display contents of configuration files
display_config_files() {
  echo "### Configuration Files ###"
  config_files=$(find . -type f \( -name "*.yml" -o -name "*.json" -o -name "*.rb" -o -name "*.js" -o -name "Dockerfile" \))
  for file in $config_files; do
    echo "### Contents of $file ###"
    cat "$file"
    echo
  done
}

# Provide insights into large files
large_files_report() {
  echo "### Large Files ###"
  find . -type f -exec du -h {} + | sort -rh | head -n 5
}

# Summarize recent Git commits, branches, and status
git_summary() {
  echo "### Git Information ###"
  echo "Recent Commits:"
  git log --oneline -5
  echo
  echo "Branches:"
  git branch -a
  echo
  echo "Git Status:"
  git status -s
}

# Output the result
output_result() {
  local output_format=$1
  local output_file=$2

  if [[ "$output_format" == "clipboard" ]]; then
    cat "$output_file" | xclip -selection clipboard
    echo "Output copied to clipboard."
  elif [[ "$output_format" == "file" ]]; then
    echo "Output written to $output_file."
  else
    cat "$output_file"
  fi
}

# Main function to run all checks and output results
main() {
  local output_format=${1:-"terminal"}
  local output_file="/tmp/repo_analysis_output.txt"

  check_dependencies
  {
    analyze_structure
    display_config_files
    large_files_report
    git_summary
  } > "$output_file"

  output_result "$output_format" "$output_file"
}

# Usage message
usage() {
  echo "Usage: $0 [output_format]"
  echo "output_format: terminal (default), clipboard, or file"
  echo "Example: $0 clipboard"
}

# Check if help is requested
if [[ "$1" == "--help" ]]; then
  usage
  exit 0
fi

# Run the main function
main "$@"

