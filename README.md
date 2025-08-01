# PR Description Generator

PR Description Generator is a Go-based CLI tool that automatically generates pull request descriptions by comparing your current branch with any branch. It supports multiple LLM providers including local Ollama models, OpenAI, and Anthropic to generate professional PR descriptions based on your actual code changes.

## Features

- Automatically compares current branch with any branch
- Extracts git diff and commit history
- Generates professional PR descriptions using multiple LLM providers
- Supports Ollama (local), OpenAI, and Anthropic
- Outputs markdown format for easy integration with GitHub CLI
- No manual input required - everything is calculated from your git repository
- Config file support for persistent settings
- Enhanced accuracy with file type analysis and response validation
- Retry logic for better reliability
- Temperature control for more focused responses

## Requirements

- Go 1.20+
- Git repository
- One of the following LLM providers:
  - **Ollama**: Local models (default: qwen2.5-coder:14b-instruct-q8_0)
  - **OpenAI**: API key and model (e.g., gpt-4, gpt-3.5-turbo)
  - **Anthropic**: API key and model (e.g., claude-3-sonnet-20240229)

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

Create a `.goprrc` file in your project root or home directory. The configuration depends on your chosen provider:

**For Ollama (local models):**
```ini
provider=ollama
base_url=http://localhost:11434
model=devstral:latest
temperature=0.1
```

**For OpenAI:**
```ini
provider=openai
api_key=your_openai_api_key_here
model=gpt-4
temperature=0.1
```

**For Anthropic:**
```ini
provider=anthropic
api_key=your_anthropic_api_key_here
model=claude-3-sonnet-20240229
temperature=0.1
```

The tool will look for config files in this order:

1. `.goprrc` in current directory
2. `~/.goprrc` in home directory

### Environment Variables

You can also use environment variables:

- `GOPR_PROVIDER`
- `GOPR_MODEL`
- `GOPR_API_KEY`
- `GOPR_BASE_URL`
- `GOPR_TEMPERATURE`

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

- `-provider`: LLM provider (ollama, openai, anthropic)
- `-model`: Model to use (varies by provider)
- `-api-key`: API key for the provider (required for OpenAI/Anthropic)
- `-base-url`: Base URL for the provider (optional, defaults vary by provider)
- `-temperature`: Temperature for generation (default: 0.1)
- `-branch`: Branch to compare current changes against (default: `main`)
- `-verbose`: Enable verbose output for debugging

### Examples

**Using Ollama (local models):**
```bash
./gopr -provider ollama -model codellama
```

**Using OpenAI:**
```bash
./gopr -provider openai -model gpt-4 -api-key your_api_key_here
```

**Using Anthropic:**
```bash
./gopr -provider anthropic -model claude-3-sonnet-20240229 -api-key your_api_key_here
```

**Enable verbose output:**
```bash
./gopr -verbose
```

**Use a remote Ollama instance:**
```bash
./gopr -provider ollama -base-url http://192.168.1.100:11434
```

**Full command with all parameters:**
```bash
./gopr -provider openai -model gpt-4 -api-key your_key -temperature 0.1 -branch main -verbose
```

## Recommended Models

Based on testing, these models perform best for PR description generation:

### Ollama Models (Local)
1. **devstral:latest** (23.6B) - Best accuracy and understanding of code changes
2. **phi4:latest** (14.7B) - Good balance of size and accuracy
3. **deepseek-coder-v2:latest** (15.7B) - Code-focused but may be less accurate
4. **qwen2.5-coder:latest** (32B) - Large model, good accuracy but slower
5. **qwen2.5-coder:14b-instruct-q8_0** (14B) - Good accuracy but slower

### OpenAI Models
1. **gpt-4** - Excellent code understanding and PR description quality
2. **gpt-3.5-turbo** - Good performance with faster response times
3. **gpt-4-turbo** - Best balance of quality and speed

### Anthropic Models
1. **claude-3-sonnet-20240229** - Excellent code analysis and PR descriptions
2. **claude-3-haiku-20240307** - Fast and efficient for smaller changes
3. **claude-3-opus-20240229** - Highest quality but slower responses

## How It Works

1. **Branch Detection**: Determines your current branch name and the one you want to compare it with using `main` or the provided one by the `branch` flag
2. **Diff Generation**: Compares your current branch with main or `branch` using `git diff <branch>...`
3. **Commit History**: Extracts commit messages since the desired branch
4. **File Analysis**: Analyzes what types of files were changed
5. **LLM Processing**: Formats the information and sends it to the configured LLM provider with low temperature (0.1)
6. **Response Validation**: Checks for generic responses and retries if needed
7. **Output**: Returns a professional PR description in markdown format

## Project Structure

- `cmd/main.go`: CLI entry point with config file support
- `internal/models/`: Defines the LLM provider interface and configuration structures
- `internal/service/`: Contains the PR generation logic and LLM provider implementations

## Contributing

Contributions are welcome! Please open an issue or submit a PR.
