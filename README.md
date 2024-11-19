# presidium-rst-to-markdown (rst2md)

A cli tool: rst2md, convert reStructuredText to a Presidium Docsite.

## Prerequesite: Pandoc

As rst2md shells out to pandoc, before you install rst2md please [install]
(https://pandoc.org/installing.html) [Pandoc](https://pandoc.org/) for your platform.

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

## Development

### Nodejs

You should have LTS version of NodeJS installed

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

```shell
brew install pre-commit
```

Run to install git hook scrips

```shell
pre-commit install
```
