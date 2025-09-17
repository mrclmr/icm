[![build](https://github.com/mrclmr/icm/actions/workflows/build.yml/badge.svg)](https://github.com/mrclmr/icm/actions/workflows/build.yml)  [![Go Report Card](https://goreportcard.com/badge/github.com/mrclmr/icm)](https://goreportcard.com/report/github.com/mrclmr/icm)

# icm (intermodal container markings)

icm generates or validates single data or whole data sets of intermodal container markings according to [ISO 6346](https://en.wikipedia.org/wiki/ISO_6346).

See examples for [`generate` command](docs/icm_generate.md#examples) and [`validate` command](docs/icm_validate.md#examples).

## Demo

![Demo](docs/gif/demo.gif)

## Installation

### macOS with [Homebrew](https://brew.sh)

```
brew install --cask mrclmr/tap/icm
```

### Linux with [Homebrew on Linux](https://docs.brew.sh/Homebrew-on-Linux)

```
brew install --cask mrclmr/tap/icm
```

### Windows with [Scoop](https://scoop.sh)

```
scoop bucket add mrclmr-bucket https://github.com/mrclmr/scoop-bucket.git
scoop install icm
```

### Download binary

See binaries in the [Releases](https://github.com/mrclmr/icm/releases) section.

Find help how to generate shell completion and man pages:
```
icm completion -h && icm doc man -h
```

## Development

1. Requirements
    * [Golang latest version](https://golang.org/doc/install)
    * [golangci-lint latest version](https://github.com/golangci/golangci-lint#install-golangci-lint)
    * [GNU Make 4.x.x](https://www.gnu.org/software/make/)

2. To build project execute
    ```
    make
    ```

## Release

1. Dry run with `goreleaser`
    ```
    goreleaser release --clean --skip=validate --skip=publish
    ```

2. Create version tag according to [SemVer](https://semver.org)
    ```
    git tag 'v0.0.1'
    ```

3. Push tag and let GitHub Actions and Goreleaser do the work
    ```
    git push --tags
    ```

## License

icm is released under the MIT license. See [LICENSE](https://github.com/mrclmr/icm/blob/master/LICENSE)
