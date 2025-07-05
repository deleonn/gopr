# prgo

prgo is a Go-based CLI tool that automatically generates pull request descriptions by comparing your current branch with the main branch. It uses a local Ollama instance to generate professional PR descriptions based on your actual code changes.

## Features

- Automatically compares current branch with main branch
- Extracts git diff and commit history
- Generates professional PR descriptions using local Ollama models
- Outputs markdown format for easy integration with GitHub CLI
- No manual input required - everything is calculated from your git repository

## Requirements

- Go 1.20+
- Git repository with a main branch
- Ollama installed and running locally
- An Ollama model (default: llama3.2)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/deleonn/gopr.git
   cd gopr
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the binary:

   ```bash
   go build -o gopr cmd/main.go
   ```

## Usage

### Basic Usage

Generate a PR description for your current branch:

```bash
./gopr
```

This will output the generated PR description to stdout, which you can then pipe to other tools.

### With GitHub CLI

Generate and update your PR description:

```bash
./gopr | gh pr edit --body-file -
```

### Copy to Clipboard

On macOS:
```bash
./gopr | pbcopy
```

On Linux:
```bash
./gopr | xclip -selection clipboard
```

### Command Line Options

```bash
./gopr -h
```

Available options:
- `-ollama-url`: Ollama server URL (default: http://localhost:11434)
- `-model`: Ollama model to use (default: llama3.2)
- `-verbose`: Enable verbose output for debugging

### Examples

Use a different model:
```bash
./gopr -model codellama
```

Enable verbose output:
```bash
./gopr -verbose
```

Use a remote Ollama instance:
```bash
./gopr -ollama-url http://192.168.1.100:11434
```

## How It Works

1. **Branch Detection**: Determines your current branch name
2. **Diff Generation**: Compares your current branch with main using `git diff main...`
3. **Commit History**: Extracts commit messages since the main branch
4. **LLM Processing**: Formats the information and sends it to Ollama
5. **Output**: Returns a professional PR description in markdown format

## Project Structure

- `cmd/main.go`: CLI entry point
- `internal/service/`: Contains the PR generation logic and Ollama integration
- `internal/models/`: Data models (legacy, can be removed)

## Contributing

Contributions are welcome! Please open an issue or submit a PR.
