# Hello World in Go

Issue the following commands on the command line terminal.
```bash
$ go mod init github.com/lreimer/go-for-operations/hello-go
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
        fmt.Println("Hello World 😀")
}
```

Save the file, open a terminal and build your first binary Go application.
```bash
$ go run main.go

$ go build
$ go build -o hello-go main.go

$ ./hello-go
```

