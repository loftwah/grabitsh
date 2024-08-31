# Grabit.sh

Grabit.sh is a powerful command-line tool designed to quickly gather and summarize useful information from Git repositories. It provides developers with a comprehensive overview of a project's structure, Git information, important files, and more.

## Features

- Repository structure visualization
- Git information summary (recent commits, branches, remotes, status)
- Identification of important configuration files
- Large file detection
- File type summary
- Recently modified files list
- Project type detection

## Installation

To install Grabit.sh, follow these steps:

1. Clone the repository:

   ```bash
   git clone git@github.com:loftwah/grabit.sh.git
   ```

2. Navigate to the project directory:

   ```bash
   cd grabit.sh
   ```

3. Build the project:

   ```bash
   go build -o grabitsh
   ```

## Usage

To use Grabit.sh, run the following command in your terminal:

```bash
./grabitsh --output <output_method>
```

Replace `<output_method>` with one of the following options:

- `stdout`: Display the output in the terminal (default)
- `clipboard`: Copy the output to your clipboard
- `file`: Save the output to a file (use the `-f` flag to specify the file path)

### Examples

1. Display output in the terminal:

   ```bash
   ./grabitsh --output stdout
   ```

2. Copy output to clipboard:

   ```bash
   ./grabitsh --output clipboard
   ```

3. Save output to a file:

   ```bash
   ./grabitsh --output file -f output.txt
   ```

## Web Server

Grabit.sh also includes a web server feature. To start the web server, use the following command:

```bash
./grabitsh serve
```

This will start a web server on port 42069. You can access it by navigating to `http://localhost:42069` in your web browser.

## Project Structure

```bash
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
