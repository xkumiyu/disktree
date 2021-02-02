# Disk Tree

[![Github Actions][ci-status]][ci]

[ci]: https://github.com/xkumiyu/disktree/actions
[ci-status]: https://github.com/xkumiyu/disktree/workflows/test/badge.svg

Disk Tree is a CLI that displays the file size according to the directory structure,
like the `tree` command.

![screenshot](https://user-images.githubusercontent.com/6437204/103475169-20a26180-4dee-11eb-94eb-fdfd1310dd98.png)

## Installation

Install using homebrew:

```sh
brew install xkumiyu/tap/disktree
```

If you want to use the HEAD, you can do the following:

```sh
go get github.com/xkumiyu/disktree/cmd/dtree
```

## Usage

You can use `dtree` command.

```sh
dtree [OPTIONS] PATH
```

### Options

- `-L, --max-depth INT`: show only to max-depth
- `-s, --sort STRING`: select sort: name(default), size, and files
- `--min-size STRING`: show only files/dirs larger than min-size
- `--color STRING`: set the colorization: auto(default), always, and never

### Example

```sh
dtree -L 2 -s size -min-size 1M ~/
```

## Comparison

The processing time when running in a large directory (>1M files) is as follows:

| command | time |
| :-- | :-- |
| dtree | 45.6s +/- 4.8s |
| tree -a -s | 117.2s +/- 3.8s |
