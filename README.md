# Recog-Go: Pattern Recognition using Rapid7 Recog

This is a Go implementation of the [Recog](https://github.com/rapid7/recog/) library and fingerprint database from Rapid7.

Recog-Go is open source, please see the [LICENSE](https://raw.githubusercontent.com/runZeroInc/recog-go/master/LICENSE) file for more information.

To build and install:
```
$ git clone https://github.com/giovanni-bellini-argo/recog-go.git /path/to/recog
$ go install .
```

# Purpose

This repo has as goal to build a tool to make recog-go a standalone to implement in other projects, as well as to make some custom variations

# variations
1. modifyed the pattern matching strings (regex) to match without the ^ and $ anchors constraints