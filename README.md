# Disk Tree

Disk Tree is a CLI that displays the file size according to directory structure, like the `tree` command.

## Installation

```sh
go get github.com/xkumiyu/disktree
```

## Usage

```sh
dtree .
```

### Output Example

```text
  2M ./ [26 files]
|--  23K .git/ [18 files]
|   |--   0B FETCH_HEAD
|   |--  21B HEAD
|   |-- 187B config
|   |--  73B description
|   |--  22K hooks/ [13 files]
|   |   |-- 478B applypatch-msg.sample
|   |   |-- 896B commit-msg.sample
|   |   |--   4K fsmonitor-watchman.sample
|   |   |-- 189B post-update.sample
|   |   |-- 424B pre-applypatch.sample
|   |   |--   1K pre-commit.sample
|   |   |-- 416B pre-merge-commit.sample
|   |   |--   1K pre-push.sample
|   |   |--   4K pre-rebase.sample
|   |   |-- 544B pre-receive.sample
|   |   |--   1K prepare-commit-msg.sample
|   |   |--   2K push-to-checkout.sample
|   |   `--   3K update.sample
|   |-- 240B info/ [1 files]
|   |   `-- 240B exclude
|   |--   0B objects/
|   |   |--   0B info/
|   |   `--   0B pack/
|   `--   0B refs/
|       |--   0B heads/
|       `--   0B tags/
|-- 883B .gitignore
|--  55B CHANGELOG.md
|--   1K LICENSE
|-- 118B Makefile
|-- 223B README.md
|--  44B go.mod
`--   2K main.go
```
