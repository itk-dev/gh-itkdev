# gh-itkdev

[GitHub CLI](https://docs.github.com/en/github-cli) extension for ITK Development.

## Installation

```shell
gh extension install itk-dev/gh-itkdev
```

## Usage

### `changelog`

Manage changelog based on [keep a changelog](https://keepachangelog.com/en/1.1.0/):

```shell
gh itkdev changelog --help
```

Create changelog:

```shell
gh itkdev changelog --create
```

Update changelog for a pull request:

```shell
gh itkdev changelog --fucking-changelog
```

Update changelog for a release (`«tag»`):

```shell
gh itkdev changelog --release «tag»
```

## Development

``` shell
task
```
