# Go EC2 CLI using Cobra and AWS SDK

This demo shows how to create a simple Cobra CLI app to manage EC2 instances
using the Golang AWS SDK. Example usage:

Example usage:
```bash
$ go-ec2 version
$ go-ec2 run --type t2.micro demo-vm
$ go-ec2 describe
$ go-ec2 chaos --count 2

$ go-ec2 start <instanceId>
$ go-ec2 stop <instanceId>
$ go-ec2 terminate <instanceId>
```

# Initial application creation

The initial skaffolding of the Go project and CLI application skeleton is taken care
of using the `cobra` CLI utility program.

```bash
$ cobra init go-ec2 --pkg-name github.com/lreimer/go-for-operations/go-ec2 --license MIT --author "Mario-Leander Reimer"
$ cd go-ec2
$ go mod init github.com/lreimer/go-for-operations/go-ec2
$ go build
```

Finally, add the `Makefile` and a `.goreleaser.yml` to build the binary distribution.

## Exercise: Add Describe Subcommand

Next, we want to add a `describe` subcommand that lists all available EC2 instances in an AWS region.

```bash
$ cobra add describe
```

Open the created `cmd/describe.go` file, perform some cleanup and add the following business logic to the subcommand run function.

```golang
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
```

```golang
    // Create an EC2 session instance
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(Region())},
    )

    // Create EC2 service client
    svc := ec2.New(sess)

    // Call to get detailed information on each instance
    result, err := svc.DescribeInstances(nil)
    if err != nil {
        fmt.Println("Error describing EC2 instances.", err)
    } else {
        fmt.Println("Instances:", result)
    }
```

## Exercise: Add Chaos Subcommand

Next, we want to add a `chaos` subcommand that terminates random tagged EC2 instances in an AWS region.

```bash
$ cobra add chaos
```

Open the created `cmd/chaos.go` file, perform some cleanup and add the required business logic to the subcommand run function.


## Bonus: Add further Subcommands

Next, add the planned subcommands for our CLI application. Once created, perform some
cleanup and remove unwanted code.

```bash
$ cobra add run
$ cobra add start
$ cobra add stop
$ cobra add terminate
```

## References

- https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/
- https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/using-ec2-with-go-sdk.html
- https://github.com/awsdocs/aws-doc-sdk-examples/tree/master/go/example_code/ec2
