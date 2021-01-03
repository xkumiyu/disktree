# Disk Tree

[![Github Actions][ci-status]][ci]

[ci]: https://github.com/xkumiyu/disktree/actions
[ci-status]: https://github.com/xkumiyu/disktree/workflows/Go/badge.svg

Disk Tree is a CLI that displays the file size according to directory structure,
like the `tree` command.

![screenshot](https://user-images.githubusercontent.com/6437204/103475169-20a26180-4dee-11eb-94eb-fdfd1310dd98.png)

## Installation

Install using homebrew:

```sh
brew install xkumiyu/tap/disktree
```

Alternatively, you can use `go get` to install:

```sh
go get github.com/xkumiyu/disktree/cmd/dtree
```

## Usage

You can use `dtree` command.

```sh
dtree /path/to/dir
```

### Options

- `-max-depth <int>`: show only to max-depth
- `-sort <string>`: select sort: name(default), size
- `-no-color`: disable colorization
