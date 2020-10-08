# K8s Log Shipping Sidecar

This demo is a simple sidecar container than scans a directory for files
to be processed and shipped.

## Basic Implementation

```bash
$ go mod init github.com/lreimer/go-for-operations/k8s-logship-sidecar
$ touch main.go
```

```golang
package main

package (
    "log"
)

func main() {
    log.Println("Started Log Shipping Sidecar")
}
```

## Building and Running

```bash
$ make build
$ make docker

$ docker run -it -v `pwd`:/logs k8s-logship-sidecar:v1.0.0

$ kubectl apply -f nginx-deployment.yaml
```
