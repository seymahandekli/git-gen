# git-gen

`git-gen` is a command-line tool developed in Go that generates commit messages and code reviews based on code changes in your project by utilizing OpenAI's ChatGPT API and Ollama API.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Introduction

`git-gen` is designed to assist developers in creating detailed commit messages and/or performing code reviews automatically depending on their codebase changes. By leveraging the power of ChatGPT and Ollama, `git-gen` analyzes the changes made to the code and generates meaningful output.

## Features

- Generate commit messages based on code changes.
- Performs detailed code reviews.

## Installation

To get started with `git-gen`, you need to have Go installed on your machine. You can download and install Go from [here](https://golang.org/dl/).

Once Go is installed, you can clone the `git-gen` repository and build the tool:

```sh
git clone https://github.com/seymahandekli/git-gen
cd git-gen

go build ./cmd/git-gen
./git-gen register
```

You can install the package should put the /usr/local/go/bin directory in your PATH environment variable.

```sh
go install github.com/seymahandekli/git-gen/cmd/git-gen@latest
```

## Usage

After building `git-gen`, you can use it from the command line. Below are some example commands to help you get started:

### Using OpenAI API

#### Alternative 1:

```sh
# Generate commit message based on your git diff command choices (default platform: openai)
git gen commit --source "commitID" --dest "commitID" --apikey "PLATFORM_API_KEY" --platform openai --model YOUR_MODEL

# default `git diff HEAD´ command
git gen commit --apikey "PLATFORM_API_KEY"

# Generate code review
git gen review --apikey "PLATFORM_API_KEY"
```

#### Alternative 2:

You don't have to specify OPENAI API KEY explicitly, you may store it to `PLATFORM_API_KEY` environment variable.

```sh
export PLATFORM_API_KEY="PLATFORM_API_KEY"

# Generate commit message based on your git diff command choices (default platform: openai and default model: gpt4-o)
git gen commit --source "commitID" --dest "commitID"

# default `git diff HEAD´ command
git gen commit

# Generate code review
git gen review
```

### Using Ollama API (No API Key Required)

You can use the Ollama API, which does not require an API key, for free

```sh
# Generate commit message based on your git diff command choices (default model: llama3)
git gen commit --source "commitID" --dest "commitID" --platform ollama --model YOUR_MODEL

# default `git diff HEAD´ command
git gen commit --platform ollama --model YOUR_MODEL

# Generate code review
git gen review --platform ollama --model YOUR_MODEL
```


For more detailed usage instructions, refer to the `--help` option:

```sh
# General usage
git gen --help

# Help for commit messages feature
git gen commit --help

# Help for code review feature
git gen review --help
```

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
