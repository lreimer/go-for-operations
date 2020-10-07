# kubectl-choas Plugin using Cobra

This demo shows how to create a Kubernetes plugin by using a Cobra CLI app to 
perform chaos monkey style operations on Kubernetes deployments and pods.

Example usage:
```bash
$ kubectl chaos --replicas 2 deployment/nginx
$ kubectl chaos --help
$ kubectl chaos version
```

## Initial application creation

The initial skaffolding of the Go project and CLI application skeleton is taken care
of using the `cobra` CLI utility program.

```bash
$ cobra init kubectl-chaos --pkg-name chaos --license MIT --author "Mario-Leander Reimer"
$ cd kubectl-chaos
$ go mod init github.com/lreimer/go-for-operations/kubectl-chaos

$ go install
$ kubectl chaos
```

Add the `Makefile` and a `.goreleaser.yml` to build the binary distribution.
