# Klaus

This service is intended as go file-watcher.

Just use the "klaus" serivce as Docker base image and put your go-Files into "/go/src/myproject". Klaus will watch for changes in the current workdir and do a `go build`.

Example use:

```
FROM klaus

WORKDIR /go/src/myproject

COPY ./src/myproject ./
```