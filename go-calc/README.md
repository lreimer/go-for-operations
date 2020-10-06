# Basic CLI and Tooling with Go

## Writing the Go-Calc CLI

Issue the following commands on the command line terminal.
```bash
$ go mod init github.com/lreimer/go-for-operations/go-calc
$ touch main.go
$ touch calc.go

# open VS.code for editing the source files.
$ code .
```

Open the `calc.go` file the your IDE and add the following code snippet to the file. 

```golang
package main

import "strconv"

// Add two numbers represented as string
func Add(a string, b string) int64 {
	x, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		panic(err)
	}

	y, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		panic(err)
	}

	return x + y
}
```

Open the `main.go` file in your IDE and add the following code snippet to the file.

```golang
package main

import (
	"flag"
	"fmt"
	"os"
)

var version = ""
var commit = ""

func main() {
	// add CLI subcommand and boolean flag
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addEnabled := addCmd.Bool("enabled", true, "enabled")

	if len(os.Args) < 2 {
		illegalArguments()
	}

	// decide based on first argument
	switch os.Args[1] {
	case "add":
		// parse the remaing arguments
		addCmd.Parse(os.Args[2:])
		if *addEnabled {
			args := addCmd.Args()
			result := Add(args[0], args[1])
			fmt.Printf("%v + %v = %v\n", args[0], args[1], result)
		}
	case "version":
		fmt.Printf("go-calc %v %v\n", version, commit)
	default:
		illegalArguments()
	}
}

func illegalArguments() {
	fmt.Println("Expected 'add' or 'version' subcommands.")
	os.Exit(1)
}
```

Now, build and run the application for the first time.

```bash
$ go run main.go version

$ go build -o go-calc 

# omit the symbol table, debug information and the DWARF table
$ go build -o go-calc -ldflags="-s -w -X main.Version=v1.0.0"

$ ./go-calc version
```

Add another command to the CLI, for example to subtract or multiply two numbers.

## Testing Go Applications

Go brings out-of-the-box support for unit testing your code. First, create an empty file called `calc_test.go` and add the following code:

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

Any files ending with `_test.go` are ignored during build. In order to execute the Go
tests simply run `go test -v`.

## Building Go Applications

### Using Makefiles

To save some typing on your Go commands the good old `make` tool is helpful and easy to use. Create a `Makefile` with the following content:

```
VERSION = v1.0.0

.PHONY: build

clean:
	@go clean

test:
	@go test -v

build:
	# omit the symbol table, debug information and the DWARF table
	@go build -o go-calc -ldflags="-s -w -X main.version=$(VERSION)"

all: clean build test
```

### Using GoReleaser

The GoReleaser utility can be used to cross-compile and publish Go applications for many target platforms
and distribution formats such as Docker images, Brew taps and archives.

```bash
# brew install goreleaser
# curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh

$ goreleaser init
$ goreleaser --snapshot --skip-publish --rm-dist
```

Open the created `.goreleaser.yml` in your editor and add the following sections and information:
```yaml
project_name: go-calc
builds:
  - 
    # add this to the builds section
    ldflags: -s -w -X main.version={{.Version}}
archives:
  - name_template: '{{ .ProjectName }}-{{ .Version }}-{{ .Os }}-{{ .Arch }}{{ if .Arm}}v{{ .Arm }}{{ end }}'
    # add this to the archives section
    format_overrides:
     - goos: windows
       format: zip
    replacements:
      darwin: Darwin
      linux: Linux
      windows: Windows
      386: i386
	  amd64: x86_64
# add this to also create docker images
dockers:
  - image_templates:
      - lreimer/go-calc:latest
      - lreimer/go-calc:v{{ .Major }}
      - lreimer/go-calc:{{ .Version }}
    skip_push: true
    dockerfile: Dockerfile_goreleaser
    goos: linux
    goarch: amd64
```

Also, create a new file called `Dockerfile_goreleaser` and add the following contents:
```
FROM gcr.io/distroless/static-debian10
COPY go-calc /

ENTRYPOINT ["/go-calc"]
CMD ["version"]
```

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
$ docker build -t go-calc:v1.0.0 .
$ docker images

$ docker run go-calc:v1.0.0
$ docker run go-calc:v1.0.0 add 1 1
```