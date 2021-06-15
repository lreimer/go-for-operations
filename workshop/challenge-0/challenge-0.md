# Challenge 0 - Setup the environment

## Install required tools (15 min)

You need the following tools for this workshop installed:

- Golang runtime - install for your OS using this document: https://golang.org/doc/install
- Git for pulling this repo
- Editor with support for Go
    - Visual Studio Code with Go plugin
    - IntelliJ IDEA with Go plugin
    - JetBrains GoLand
- Local Docker installation
    - Windows/Mac: https://docs.docker.com/engine/install/
    - Linux: https://get.docker.com/
- Kubectl https://kubernetes.io/docs/tasks/tools/#kubectl
- Minikube: https://minikube.sigs.k8s.io/docs/start/

## Test your Go Runtime (10 min)

1) Open your editor of choice

2) Create a new file `helloworld.go`.

3) Put the following snippet into this file:

    ```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello World!")
    }
    ```

    or use the sample in /hello-go

4) Run the file with `go run helloworld.go`

5) If you see the output, you are ready to GO!