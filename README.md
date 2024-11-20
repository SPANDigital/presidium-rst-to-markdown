# presidium-rst-to-markdown (rst2md)

[![Open in Dev Containers](https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/SPANDigital/presidium-rst-to-markdown)
![Go Dev](https://github.com/SPANDigital/presidium-rst-to-markdown/actions/workflows/go.yml/badge.svg?branch=dev)
![Go Main](https://github.com/SPANDigital/presidium-rst-to-markdown/actions/workflows/go.yml/badge.svg?branch=main)

A cli tool: rst2md, convert reStructuredText to a Presidium Docsite.

## Prerequesite: Pandoc

As rst2md shells out to pandoc, before you install rst2md please [install](https://pandoc.org/installing.html)
[Pandoc](https://pandoc.org/) for your chosen platform.

## Installation

Choose one installation method, they are listed in order of preference

### Via homebrew (recommended)

This requires [homebrew](https://brew.sh/) to be installed.

```shell
brew tap SPANDigital/homebrew-tap
brew update
brew install rst2md
```

### Via go install (for go developers)

This requires as least [Go 1.22.x](https://go.dev/doc/install) to be installed for your operating system.

```shell
go install github.com/spandigital/presidium-rst-to-markdown/cmd/rst2md
```
### Source code and pre-built binaries

Source code and pre-built binaries for various architecutrs can be found at [releases](/releases)

## Usage

```
Usage of rst2md:
  -depth int
        Heading depth level to split sections (default 2)
  -force
        Force overwrite of output directory
  -input string
        Input directory
  -output string
        Output directory
  -pandoc-path string
        Path to the Pandoc executable (default "pandoc")
  -parallel int
        Maximum number of parallel processes (default 4)
  -v    Enable verbose logging
```

## Tools Required for Development

#### Golangci-lint

```brew
brew install golangci-lint
brew upgrade golangci-lint
```

### Nodejs

You should have LTS version of NodeJS installed.

```shell
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.0/install.sh | bash
nvm install --lts
```

### commitlint

For linting conventional commits

```shell
npm install -g @commitlint/cli @commitlint/config-conventional
```

### Pre-commit hooks

To install pre-commit

On MacOSX
```shell
brew install pre-commit
```

## Clone and install pre-commit hooks

```shell
git clone git@github.com:SPANDigital/presidium-rst-to-markdown.git
pre-commit install
```

## Commit messages

Commit messages to adhere to the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) standard.
