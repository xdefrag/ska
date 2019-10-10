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

So you can turn this

````
// <template>/templates/Dockerfile

FROM golang:alpine as builder

WORKDIR /project

{{if len .addons}}
RUN set -xe && \
    apk update && apk upgrade && \
    apk add --no-cache make {{if has "certs" .addons}}ca-certificates {{end}}git curl{{if has "migrate" .addons}} && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v3.5.2/migrate.linux-amd64.tar.gz | tar xvz && \
    cp migrate.linux-amd64 /migrate{{end}}
{{end}}

COPY . .

RUN make dep && \
    make build-{{.svc}}{{if has "health" .addons}} && \
    make build-health{{end}}

FROM scratch

{{if has "certs" .addons}}COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/{{end}}
{{if has "migrate" .addons}}COPY --from=builder /migrate /migrate
COPY --from=builder /project/internal/{{.svc}}/migrations /migrations{{end}}
COPY --from=builder /project/dist/{{.svc}} /{{.svc}}
{{if has "assets" .addons}}COPY --from=builder /project/internal/notify/assets /assets{{end}}
{{if has "health" .addons}}COPY --from=builder /project/dist/health /health

HEALTHCHECK CMD ["/health"]{{end}}

{{if len .expose}}EXPOSE{{if has "metrics" .expose}} 8080{{end}}{{if has "grpc" .expose}} 8086{{end}}{{end}}

CMD ["/{{.svc}}{{if .cmd}} {{.cmd}}{{end}}"]
````
````
// <template>/values.toml

svc = "example"

addons = ["certs"]

expose = ["metrics"]
````

into this

````
FROM golang:alpine as builder

WORKDIR /project

RUN set -xe && \
    apk update && apk upgrade && \
    apk add --no-cache make ca-certificates git curl

COPY . .

RUN make dep && \
    make build-example

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /project/dist/example /example

EXPOSE 8080
````

with one simple command. [More examples](https://github.com/xdefrag/ska/tree/master/examples).

## Usage
````sh
$ ska your_template
````
$EDITOR will be opened with values.toml copy. After you save and quit templates will be recursively executed and copied to current folder.  
SKA will work out of the box with any console editor (vim, emacs), for others such as vscode or atom [see this comment](https://github.com/xdefrag/ska/issues/27#issuecomment-500422334).

## Install
````sh
$ go install https://github.com/xdefrag/ska
````
