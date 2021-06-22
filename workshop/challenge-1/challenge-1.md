# Challenge 1 - Command line dev tools

## Simple CLIs with plain Go (90 min)

The Go base library has a lot of useful libraries, so we do not even need to search for a special library to create our first CLI.


### Go modules

To initialize a go project, one needs to create a go module. The `go.mod` file defines the package dependencies of a project. It can be initialized with the go CLI and the *mod* subcommand.

Execute the following commands to initialize the go module and create the necessary go code files which we fill in the next step.
```
$ go mod init go-calc
$ touch main.go
$ touch calc.go
```
### Exercise: Create a calculator

Our first sample CLI is a simple command line calculator.
The calculator should be able to add two numbers. 
Example usage:  
```bash
go-calc add <left-operand> <right-operand>
```

Put the command and argument parsing code in the `main.go`, and the calculation code in the `calc.go`.

As the command arguments are strings, you need to use the `strconv` library to convert a string to an int. 
Example use:
```golang
    a := "10"
    b, err := strconv.ParseInt(a, 10, 64) // converts the string to a base10 int with bitlength 64
    // ... need to handle the error as conversion could fail
    fmt.Println(b) // prints 10
```

For parsing flags and subcommands we will use the `flag` library. Create a subcommand for one of the four mathematic operations and a version subcommand. 

Use a [FlagSet](https://golang.org/pkg/flag/#FlagSet) for each subcommand. 
Example for using *FlagSet*:
```golang
	// add CLI subcommand and boolean flag
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addEnabled := addCmd.Bool("enabled", true, "enabled") // returned value is pointer to the argument - deference with *addEnabled to get the actual value
```
Here the `enabled` flag is a sample flag which disables calculation.

To retrieve the command and all arguments, use the `os` library and its `Args` variable.   
Parse the subcommands with a simple switch case over `os.Args[1]` which specifies the subcommand. 

In the specific subcommand case you need to then parse the arguments after the subcommand and the flags:
```golang
    // parse remaining arguments 
    addCmd.Parse(os.Args[2:])
    
    // retrieving the arguments which are the operands for the calculation
    args = addCmd.Args()
    // call calculator add method...
```

The *version* subcommand simply prints the version of the CLI. 

Also, add a default case to handle wrong argument usage.

### Solution
Go to the folder `go-calc` and follow the steps of the [README](../../go-calc/README.md) if you need further help. The code in the folder works out-of-the-box and you can compare your solution with it.
Depending on your knowledge level and how fast you are, you can additionally go through the `Testing Go Applications` and `Building Go Applications` parts. 

## Complex CLIs need a Cobra (120 min)

There is way more to an CLI than command and flag parsing. A modern CLI should support auto-completion, global flags vs local flags, automatic config parsing, in-built documentation with help pages and more...
That is a lot of boiler plate to implement by yourself for every CLI. That is why we need a battle-tested "framework" for CLIs that supports all these features. 
The combination of the libraries [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) is exactly that.

Cobra is powering a lot of Cloud Native Dev Tools that are used by platform engineers every day:
- Kubectl
- Helm
- Docker
- Hugo CLI
- ...

With this track record, it is essential tooling for any CLI in Go. 

Let's try building a CLI that can be used as Kubectl plugin. 

### Ready our testing cluster

First, make sure that your local Kubernetes cluster is running using Minikube.
If you have not already started it, run `minikube start`. This will spin up the Kubernetes cluster which we will use to test our plugin.

Check whether the cluster is ready using `kubectl version`.
If you see a server version then the cluster is ready.

Now deploy the sample Nginx which we will use later:
`kubectl apply -f ./nginx-deployment.yaml`

### Extending kubectl with plugins

Kubectl has a simply way to integrate plugins. If you install a Go executable with the prefix `kubectl-`, it is automatically installed as kubectl plugin and can be called with kubectl. The only limitations is that plugins cannot overwrite a already existing kubectl command, for example, you cannot create a plugin with the name `kubectl-create` as it collides with the kubectl command `create`.  

The advantage of being a kubectl plugin is that you are automatically getting all args, flags and environment variables passed to your plugin so you do not need to worry about authentication or other things the kubectl command already has integrated.

### Exercise: Create a Kubernetes Chaos plugin
The purpose of the CLI we will built is to create some chaos on the cluster. 
For a specified deployment, it will kill a specified number of pods of the deployment.
Example usage:
```
kubectl chaos nginx-deployment --replicas 2 --namespace default
```
Where `replicas` is the amount of pods to delete and `namespace` the namespace in which the deployment lies. The argument `nginx-deployment` denotes the deployment for which you want to delete replicas.

We will write the chaos plugin using Cobra and the Kubernetes library *client-go*.

#### Scaffolding

To quickly generate all needed CLI code, we use the scaffolding CLI of Cobra. Install it with the following command:
`go get github.com/spf13/cobra/cobra`

Then, create a Cobra skeleton using the following commands:
```bash
$ cobra init kubectl-chaos --pkg-name chaos --license MIT --author "<your name>"
$ cd kubectl-chaos
$ go mod init kubectl-chaos
$ go mod tidy
```

#### Writing commands 
The skaffolding created a main.go which contains the entrypoint of our CLI and configures the version/commit which is later output when running the version subcommand.
In the `cmd` folder is a file `root.go` created which contains the description and run code for the root command. This one is called when you run `kubectl chaos`. Here is the integration point for adding the chaos handling. 

In the `cobra.Command` structure is the *Run* variable which you need to extend as it is called when the command is executed. You can describe the use and add documentation using the variables *Use*, *Short*, or *Long*. 
See the [documentation of the Command structure](https://pkg.go.dev/github.com/spf13/cobra#Command) for more information.

#### Add flag parsing
Start by adding the needed argument flags to the generated root command `rootCmd`:
- `--replicas 2` to specify the number of replicas
- `--namespace default` to specify the K8s namespace to use

Check the [documentation about flags](https://github.com/spf13/cobra#working-with-flags) on how to add flags to your root command.

Example:
```golang
  rootCmd.Flags().IntVarP(&replicas, "replicas", "r", 1, "Number of replicas")
```

#### Add Config File Parsing with Viper (optional)
Add config file parsing to your CLI using Viper. The user of the CLI can set the namespace or any other flag parameters in a configuration file. Integrate it into your root command with `cobra.OnInitialize(initConfig)` and check the [Viper documentation](https://github.com/spf13/viper#reading-config-files) on how to read config files.
You would reference the configuration file as `rootCmd` flag and read it our in the `initConfig` method. 

#### Add Kubernetes Client API

Next, we need to add the Kubenetes client API, in order to query a deployment and
delete the matching amount of pods for this deployment.

What the following code should do is the following:
1. Connect to the Kubernetes API by configuring a client which reads out the default kubeconfig.
2. Retrieving the deployment given as argument to the CLI in the namespace given by the argument `namespace`.
3. Retrieve all pods of the deployment by using the matching labels. Deployments use *MatchLabels* to map Pods to a specific Deployment. You can use a label selector when listing pods to find the ones that match the given deployment. 
4. Delete the number of pods given by the argument `replicas` by iterating through the list of pods. 

You can either put this code directly in the *Run* variable of the `rootCmd` or extract it into an additional function in a separate file.

```golang
import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)
```

To connect to Kubernetes, the following API calls are required (from [Example](https://github.com/kubernetes/client-go/tree/master/examples/out-of-cluster-client-configuration)):
```golang
    loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
    configOverrides := &clientcmd.ConfigOverrides{}
    kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

    config, err := kubeConfig.ClientConfig()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
```

To query the deployment and delete pods, use the following API calls (see [Documentation](https://pkg.go.dev/k8s.io/client-go/kubernetes)):
```golang
    // Get deployments
    deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), args[0], metav1.GetOptions{})
    // TODO: error handling...

    labelMap := deployment.Spec.Selector.MatchLabels

    // Get list of pods that match deployment
    pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
        LabelSelector: labels.SelectorFromSet(labelMap).String(),
    })
    // TODO: error handling...

    // TODO: handle deletion of pods...
    // Use this API call for deletion of a Pod: 
    clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})

```

#### Installation

If you finished coding and want to try the plugin, you need to first install it.
Go to the root of your Go module and run the following command: 
```bash
go install
```

Now, when your minikube is running and the nginx deployment was applied you can your plugin with the following command:
```bash
kubectl chaos nginx-deployment --replicas 2 --namespace default
```

With `kubectl get pods -n default` you should see pods terminating and new pods being created.

### Solution
Go to the subfolder /kubectl-chaos of this repo and follow the steps of the [README](../../kubectl-chaos/README.md) if you need further help. The code in the folder works out-of-the-box and you can compare your solution with it.
