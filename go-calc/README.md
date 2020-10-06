# Basic CLI and Tooling with Go

## Writing the Go-Calc CLI

Issue the following commands on the command line terminal.
```bash
$ go mod init github.com/lreimer/go-for-operations/go-calc
$ touch main.go

# open folder and `main.go` file in an IDE of your choice
$ code .
```

Next, open the folder and `main.go` file in an IDE of your choice and
add the following code snippet to the file.

```golang
package main

import "fmt"

func main() {
}
```

Save the file, open a terminal and build your first binary Go application.
```bash
$ go run main.go version

$ go build -o go-calc 

# omit the symbol table, debug information and the DWARF table
$ go build -o go-calc -ldflags="-s -w -X main.Version=v1.0.0"

$ ./go-calc version
```

## Testing Go Applications

First, create an empty file called `calc_test.go` and add the following code:
```golang
package main

import "testing"

func TestAdd(t *testing.T) {
	result := Add("1", "1")
	if result != 2 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", result, 2)
	}
}
```

## Building Go Applications

### Using Makefiles

To save some typing on your Go commands the good old `make` tool is helpful and easy to use. Create
a `Makefile` with the following content:
```
VERSION = v1.0.0

.PHONY: build

test:
	@go test

build:
	@go build -o go-calc -ldflags="-s -w -X main.version=$(VERSION)"

install: test build
```

### Using GoReleaser

### Using Docker

Building Go applications using Docker is also very easy and helpful, if you want to
distribute and run your application containerized. Create a `Dockerfile` with the following content:
```
FROM golang:1.15.2 as builder

WORKDIR /build

COPY . /build
RUN make build

FROM gcr.io/distroless/static-debian10
COPY --from=builder /build/go-calc /

ENTRYPOINT ["/go-calc"]
CMD ["version"]
```

The build the Docker image and run it as follows:
```
$ make docker
$ docker build -t go-calc:v1.0.0 .
$ docker images

$ docker run go-calc:v1.0.0
$ docker run go-calc:v1.0.0 add 1 1
```