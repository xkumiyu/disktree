# Disk Tree

Disk Tree is a CLI that displays the file size according to directory structure, like the `tree` command.

![screenshot](https://user-images.githubusercontent.com/6437204/103475169-20a26180-4dee-11eb-94eb-fdfd1310dd98.png)

## Installation

```sh
go get github.com/xkumiyu/disktree/cmd/dtree
```

## Usage

```sh
dtree <dir>
```

### Options

- `-max-depth <int>`: show only to max-depth
- `-sort <string>`: select sort: name(default), size
