# Grabit.sh

[DEMO VIDEO](https://www.youtube.com/watch?v=M6FjvHV1qak)

Grabit.sh is a powerful command-line tool designed to quickly gather and summarize useful information from Git repositories. It provides developers with a comprehensive overview of a project's structure, Git information, important files, and more.

![# Grabit sh Grabit sh is a powerful command-line tool designed to quickly gather and summarize useful information from Git repositories  It provides developers with a comprehensive overview of a p](https://github.com/user-attachments/assets/f1fb1602-88bc-4898-a3c4-dbbcf92a584e)

## Features

- Repository structure visualization
- Git information summary (recent commits, branches, remotes, status)
- Identification of important configuration files
- Large file detection
- File type summary
- Recently modified files list
- Project type detection
- LLM-friendly output chunks for easy integration with AI models

## Installation

### Pre-built Binaries

You can download pre-built binaries for your platform from the [Releases](https://github.com/loftwah/grabitsh/releases) page or use `curl` for direct download.

#### Linux

1. Download the `grabitsh-linux-amd64` file from the latest release:

   ```bash
   curl -L -o grabitsh-linux-amd64 https://github.com/loftwah/grabitsh/releases/latest/download/grabitsh-linux-amd64
   ```

2. Make it executable:

   ```bash
   chmod +x grabitsh-linux-amd64
   ```

3. Optionally, move it to a directory in your PATH:

   ```bash
   sudo mv grabitsh-linux-amd64 /usr/local/bin/grabitsh
   ```

#### macOS

1. Download the `grabitsh-darwin-amd64` file from the latest release:

   ```bash
   curl -L -o grabitsh-darwin-amd64 https://github.com/loftwah/grabitsh/releases/latest/download/grabitsh-darwin-amd64
   ```

2. Make it executable:

   ```bash
   chmod +x grabitsh-darwin-amd64
   ```

3. Optionally, move it to a directory in your PATH:

   ```bash
   sudo mv grabitsh-darwin-amd64 /usr/local/bin/grabitsh
   ```

#### Windows

1. Download the `grabitsh-windows-amd64.exe` file from the latest release using a browser or:

   ```bash
   curl -L -o grabitsh-windows-amd64.exe https://github.com/loftwah/grabitsh/releases/latest/download/grabitsh-windows-amd64.exe
   ```

2. Optionally, rename it to `grabitsh.exe` for convenience.

3. Add the directory containing the executable to your PATH environment variable.

### Building from Source

If you prefer to build from source or want the latest development version:

1. Ensure you have Go installed on your system.

2. Clone the repository:

   ```bash
   git clone https://github.com/loftwah/grabitsh.git
   ```

3. Navigate to the project directory:

   ```bash
   cd grabit.sh
   ```

4. Build the project:

   ```bash
   go build -o grabitsh
   ```

## Usage

To use Grabit.sh, navigate to a Git repository directory and run the following command:

```bash
grabitsh --output <output_method>
```

Replace `<output_method>` with one of the following options:

- `stdout`: Display the output in the terminal (default)
- `clipboard`: Copy the output to your clipboard
- `file`: Save the output to a file (use the `-f` flag to specify the file path)
- `llm-chunks`: Generate LLM-friendly chunks of the output (new feature)

### Examples

1. Display output in the terminal:

   ```bash
   grabitsh --output stdout
   ```

2. Copy output to clipboard:

   ```bash
   grabitsh --output clipboard
   ```

3. Save output to a file:

   ```bash
   grabitsh --output file -f output.txt
   ```

4. Generate LLM-friendly chunks:

   ```bash
   grabitsh --output llm-chunks
   ```

   This will create multiple text files, each containing a portion of the output with a preamble suitable for use with Large Language Models.

5. Customize chunk size for LLM output:

   ```bash
   grabitsh --output llm-chunks --chunk-size 50000
   ```

   This sets the chunk size to 50,000 tokens. The default is 100,000 tokens.

### LLM-Chunks Feature

The LLM-chunks output method is designed to create AI-friendly chunks of the Grabit.sh output. Each chunk includes a preamble that provides context about the tool, its purpose, and instructions for the AI model. This feature is particularly useful when you want to analyze the output using a Large Language Model or other AI tools.

Key points about LLM-chunks:

- Each chunk is saved as a separate text file (`grabitsh_chunk_1.txt`, `grabitsh_chunk_2.txt`, etc.).
- The default chunk size is 100,000 tokens, which can be customized using the `--chunk-size` flag.
- The preamble in each chunk helps the AI understand the context and purpose of the information.
- This feature makes it easy to feed the Grabit.sh output into AI models for further analysis or to generate insights about the repository.

## Web Server

Grabit.sh also includes a web server feature. To start the web server, use the following command:

```bash
grabitsh serve
```

This will start a web server on port 42069. You can access it by navigating to `http://localhost:42069` in your web browser.

## Future Development

### Expansion Plans

Grabit.sh is evolving beyond its current CLI functionality to include web-based capabilities. The following features and versions are planned for future releases:

1. **Free Hosted Version (Public Repositories)**

   - A web-based platform where users can input public Git repository URLs (GitHub, GitLab, Gitea, etc.).
   - The system will scan the repository online and generate summaries and documentation similar to the CLI version.
   - This version will be completely free for all users working with public repositories.

2. **Paid Hosted Version (Private Repositories)**

   - A subscription-based model designed for users who need to analyze private repositories.

   - Offers advanced features such as:

     - Priority processing
     - Enhanced security protocols
     - Team collaboration tools
     - Multi-repo management

   - The pricing will be tiered based on the number of repositories and users, with options for monthly or annual subscriptions.

### Key Benefits

- **Time Efficiency**: Quickly identifies essential files and generates necessary documentation, saving significant time for developers and teams.
- **User-Friendly**: Offers both CLI and web-based interfaces to cater to different user preferences.
- **Integration and Flexibility**: Compatible with any Git platform and supports multiple repositories.
- **Comprehensive Free Options**: Open-source CLI and free hosted options encourage widespread adoption, while the paid model is ideal for monetizing private, enterprise usage.

### Requirements for Future Versions

- Users will need to bring their own OpenAI API key to utilize the LLM features, which adds flexibility and reduces operating costs for the service.

## Project Structure

```go
grabit.sh/
├── cmd/
│   └── grabitsh/
│       ├── root.go
│       ├── serve.go
│       └── project_detection.go
├── main.go
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── .gitignore
```

## Contributing

Contributions to Grabit.sh are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms of the license included in the [LICENSE](LICENSE) file.

## Contact

For any queries or suggestions, please open an issue on the GitHub repository.
