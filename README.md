# PR Description Generator

PR Description Generator is a Go-based CLI tool that automatically generates pull request descriptions by comparing your current branch with the main branch. It uses a local Ollama instance to generate professional PR descriptions based on your actual code changes.

## Features

- Automatically compares current branch with main branch
- Extracts git diff and commit history
- Generates professional PR descriptions using local Ollama models
- Outputs markdown format for easy integration with GitHub CLI
- No manual input required - everything is calculated from your git repository
- **NEW**: Config file support for persistent settings
- **NEW**: Enhanced accuracy with file type analysis and response validation
- **NEW**: Retry logic for better reliability
- **NEW**: Temperature control for more focused responses

## Requirements

- Go 1.20+
- Git repository with a main branch
- Ollama installed and running locally
- An Ollama model (default: devstral:latest)

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

## Configuration

### Config File

Create a `.goprrc` file in your project root or home directory with:

```ini
ollama_url=http://localhost:11434
model=devstral:latest
```

The tool will look for config files in this order:
1. `.goprrc` in current directory
2. `~/.goprrc` in home directory

### Environment Variables

You can also use environment variables:
- `GOPR_OLLAMA_URL`
- `GOPR_MODEL`

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
- `-ollama-url`: Ollama server URL (default: from config or http://localhost:11434)
- `-model`: Ollama model to use (default: from config or devstral:latest)
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

## Recommended Models

Based on testing, these models perform best for PR description generation:

1. **devstral:latest** (23.6B) - Best accuracy and understanding of code changes
2. **phi4:latest** (14.7B) - Good balance of size and accuracy
3. **deepseek-coder-v2:latest** (15.7B) - Code-focused but may be less accurate
4. **qwen2.5-coder:latest** (32B) - Large model, good accuracy but slower

## How It Works

1. **Branch Detection**: Determines your current branch name
2. **Diff Generation**: Compares your current branch with main using `git diff main...`
3. **Commit History**: Extracts commit messages since the main branch
4. **File Analysis**: Analyzes what types of files were changed
5. **LLM Processing**: Formats the information and sends it to Ollama with low temperature (0.1)
6. **Response Validation**: Checks for generic responses and retries if needed
7. **Output**: Returns a professional PR description in markdown format

## Project Structure

- `cmd/main.go`: CLI entry point with config file support
- `internal/service/`: Contains the PR generation logic and Ollama integration
- `internal/models/`: Data models (legacy, can be removed)

## Contributing

Contributions are welcome! Please open an issue or submit a PR.
