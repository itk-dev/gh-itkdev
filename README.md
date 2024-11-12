# gh-itkdev

[GitHub CLI](https://docs.github.com/en/github-cli) extension for ITK Development.

## Installation

```shell
gh extension install itk-dev/gh-itkdev
```

**Note**: This may not (yet) work (cf. <https://github.com/itk-dev/gh-itkdev/actions/workflows/release.yml>) and you
may have to [install from source](#installing-from-source).

### Installing from source

Assuming [Go is installed](https://go.dev/doc/install), you can [install from
source](https://cli.github.com/manual/gh_extension_install) by running

``` shell name="build-and-install"
gh extension remove itkdev
rm -fr /tmp/gh-itkdev
git clone https://github.com/itk-dev/gh-itkdev /tmp/gh-itkdev
cd /tmp/gh-itkdev
task build
gh extension install .
cd -
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
task --list-all
```
