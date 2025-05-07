# git-gen

`git-gen` is a command-line tool developed in Go that generates commit messages and code reviews based on code changes in your project by utilizing OpenAI's ChatGPT API, Ollama API, Anthropic's Claude API, and Google's Gemini API.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Commands](#commands)
- [Supported Platforms](#supported-platforms)
- [Configuration](#configuration)
- [Custom Instructions](#custom-instructions)
- [Contributing](#contributing)
- [License](#license)

## Introduction

`git-gen` is designed to assist developers in creating detailed commit messages and/or performing code reviews automatically depending on their codebase changes. By leveraging the power of ChatGPT, Ollama, Gemini, and Claude, `git-gen` analyzes the changes made to the code and generates meaningful output.

## Features

- Generate commit messages based on code changes
- Perform detailed code reviews
- Support for multiple AI platforms (OpenAI, Ollama, Anthropic, Gemini)
- Customizable model selection
- Configurable token limits and timeout settings
- Test scenario and test code generation
- Custom instruction support via instructions.txt
- Native Git diff implementation for accurate change detection

## Installation

To get started with `git-gen`, you need to have Go installed on your machine. You can download and install Go from [here](https://golang.org/dl/).

Once Go is installed, you can clone the `git-gen` repository and build the tool:

```sh
git clone https://github.com/seymahandekli/git-gen
cd git-gen

go build ./cmd/git-gen
./git-gen register
```

You can install the package by putting the /usr/local/go/bin directory in your PATH environment variable:

```sh
go install github.com/seymahandekli/git-gen/cmd/git-gen@latest
```

## Commands

### commit
Generates a commit message based on your code changes.

```sh
# Basic usage
git gen commit

# With specific source and destination
git gen commit --source "commitID" --dest "commitID"

# With platform and model
git gen commit --platform openai --model gpt-4
```

### review
Performs a code review of your changes.

```sh
# Basic usage
git gen review

# With specific source and destination
git gen review --source "commitID" --dest "commitID"

# With platform and model
git gen review --platform anthropic --model claude-3-sonnet
```

### test-scenarios
Creates test scenarios based on your code changes.

```sh
# Basic usage
git gen test-scenarios

# With specific source and destination
git gen test-scenarios --source "commitID" --dest "commitID"

# With platform and model
git gen test-scenarios --platform gemini --model gemini-pro
```

### test
Creates test implementations based on your code changes.

```sh
# Basic usage
git gen test

# With specific source and destination
git gen test --source "commitID" --dest "commitID"

# With platform and model
git gen test --platform ollama --model llama2
```

### register
Registers the tool to your system for global usage.

```sh
git gen register
```

### help
Shows help information for commands.

```sh
# Show all commands
git gen help

# Show help for specific command
git gen help commit
git gen help review
git gen help test-scenarios
git gen help test
```

## Supported Platforms

`git-gen` supports multiple AI platforms, each with its own characteristics and requirements.

### OpenAI
OpenAI's models provide high-quality responses and are suitable for most use cases.

```sh
# Basic usage with OpenAI
git gen commit --platform openai --model gpt-4

# Available models:
# - gpt-4
# - gpt-3.5-turbo
# - gpt-4-turbo-preview
```

**Requirements:**
- OpenAI API key (set via `--apikey` or `PLATFORM_API_KEY` environment variable)
- Internet connection
- Paid API access

### Ollama
Ollama provides local model deployment, making it perfect for offline use and privacy-focused development.

```sh
# Basic usage with Ollama
git gen commit --platform ollama --model llama2

# Available models:
# - llama2
# - mistral
# - codellama
# - neural-chat
```

**Requirements:**
- Ollama installed locally
- No API key required
- Sufficient local computing resources

### Anthropic
Anthropic's Claude models excel at understanding and generating code-related content.

```sh
# Basic usage with Anthropic
git gen commit --platform anthropic --model claude-3-sonnet

# Available models:
# - claude-3-sonnet
# - claude-3-opus
# - claude-3-haiku
```

**Requirements:**
- Anthropic API key (set via `--apikey` or `PLATFORM_API_KEY` environment variable)
- Internet connection
- Paid API access

### Gemini
Google's Gemini models provide fast and efficient responses, particularly good for code-related tasks.

```sh
# Basic usage with Gemini
git gen commit --platform gemini --model gemini-pro

# Available models:
# - gemini-pro
# - gemini-pro-vision
```

**Requirements:**
- Google API key (set via `--apikey` or `PLATFORM_API_KEY` environment variable)
- Internet connection
- Paid API access

## Configuration

`git-gen` supports various configuration options that can be set through command-line flags or environment variables:

### Platform Options
- `--platform`: Choose between "openai", "ollama", "anthropic", or "gemini" (default: "openai")
- `--model`: Specify the AI model to use (e.g., "gpt-4", "llama2", "claude-3-sonnet", "gemini-pro")
- `--apikey`: Your platform API key (can also be set via PLATFORM_API_KEY environment variable)

### Token and Timeout Settings
- `--prompt-max-tokens`: Maximum tokens for the prompt (default: 3500)
- `--prompt-timeout`: Request timeout in seconds (default: 3600)

### Reference Settings
- `--source`: Source reference for diff (default: "HEAD")
- `--dest`: Destination reference for diff

### Example Configuration

```sh
# Using OpenAI with custom token limit
git gen commit --platform openai --model gpt-4 --prompt-max-tokens 4000

# Using Ollama with custom timeout
git gen commit --platform ollama --model llama2 --prompt-timeout 1800

# Using Anthropic with custom source/destination
git gen commit --platform anthropic --model claude-3-sonnet --source HEAD~2 --dest HEAD

# Using Gemini with custom configuration
git gen commit --platform gemini --model gemini-pro --prompt-max-tokens 4000
```

## Custom Instructions

You can provide custom instructions for any generation task by creating an `instructions.txt` file in your project root. This file will be automatically picked up and used to customize the generation process.

### Using git-gen-dotfile-generator

The easiest way to create your `instructions.txt` file is to use the [git-gen-dotfile-generator](https://github.com/seymahandekli/git-gen-dotfile-generator) tool. This web-based tool provides a user-friendly interface to generate AI instructions and rules for your project.

#### Features of git-gen-dotfile-generator:
- Specify project name, role, and expertise areas
- Define coding preferences and standards
- Add sample code to create contextual rules
- Download the generated instructions/rules directly

#### How to use:
1. Visit [git-gen-dotfile-generator](https://github.com/seymahandekli/git-gen-dotfile-generator)
2. Fill out the form with your project details:
   - Project name
   - Role
   - Expertise areas
   - Coding preferences
3. Add any sample code if needed
4. Generate and download your `instructions.txt` file
5. Place the file in your project root directory

The `instructions.txt` file will be automatically used by `git-gen` to:
- Generate more relevant commit messages
- Create more accurate code reviews
- Generate appropriate test scenarios
- Maintain consistency with your project standards

## Contributing

We welcome contributions from the community! If you'd like to contribute to `git-gen`, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Make your changes and commit them with clear messages.
4. Push your changes to your fork.
5. Submit a pull request to the `main` branch of this repository.

For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contributors

- Eser Özvataf (https://github.com/eser)
- Daniel M. Matongo (https://github.com/mmatongo)

## Acknowledgement

I would like to thank people below for their support and contributions:

- Arda Kılıçdağı (http://github.com/Ardakilic)
- Erman İmer (https://github.com/ermanimer)
- Eser Özvataf (https://github.com/eser)

---

We hope you find `git-gen` useful! If you have any questions or feedback, please feel free to open an issue on GitHub.

Happy coding!
