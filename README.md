# presidium-rst-to-markdown (rst2md)

A cli tool: rst2md, convert reStructuredText to a Presidium Docsite.

## Prerequesite: Pandoc

As rst2md shells out to pandoc, before you install rst2md please install
https://pandoc.org/installing.html [Pandoc](https://pandoc.org/) for your platform.
## Installation methods

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

```bash
go install github.com/spandigital/presidium-rst-to-markdown/cmd/rst2md
```

### Source code and pre-built binaries

Source code and pre-built binaries can be found on the [releases](/releases)

