SKA
=======
[![Go Report Card](https://goreportcard.com/badge/github.com/xdefrag/ska)](https://goreportcard.com/report/github.com/xdefrag/ska) [![Build Status](https://travis-ci.com/xdefrag/ska.svg?branch=master)](https://travis-ci.com/xdefrag/ska) [![codecov](https://codecov.io/gh/xdefrag/ska/branch/master/graph/badge.svg)](https://codecov.io/gh/xdefrag/ska)

**SKA** is simple scaffolding tool like [yeoman](https://github.com/yeoman/yo) but simpler and like [helm](https://github.com/helm/helm) templates but for everything.

Templates powered by [go template](https://golang.org/pkg/html/template/) package and has this structure:
````
.
+--~/.local/share/ska
|  +--your_template
|  |  +--values.toml   // Values for templates
|  |  +--templates     // Actual templates
|  |  |  +--main.go
|  |  |  +--Makefile
|  |  |  ...

````

## Usage
````sh
$ ska your_template
````
$EDITOR will be opened with values.toml copy. After you save and quit templates will be recursively executed and copied to current folder.

## Install
````sh
$ go install https://github.com/xdefrag/ska
````
