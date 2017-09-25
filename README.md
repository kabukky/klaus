# Klaus

[![Build Status](https://travis-ci.org/kabukky/klaus.svg?branch=master)](https://travis-ci.org/kabukky/klaus)

This service is intended as go file-watcher.

Just use the "klaus" serivce as Docker base image and put your go-Files into "/go/src/myproject". Klaus will watch for changes in the current workdir and do a `go build`.

Example use:

```
FROM kaih/klaus

WORKDIR /go/src/myproject

COPY ./src/myproject ./
```

This should NOT used as production environment!