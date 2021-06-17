# Challenge 1 - Command line dev tools

## Simple CLIs with plain Go (60 min)

The Go base library has a lot of useful libraries, so we do not even need to search for a special library to create our first CLI.

### Exercise: Create a calculator
Our first sample CLI is a simple command line calculator.
The calculator should be able to do two of the four basic mathematic operations: add, subtract, multiply and divide. 
Split the command and argument parsing code AND the calculation code in two separate files.
For parsing the arguments, use the `flag` library. Create a subcommand for each of the two mathematic operations and a version subcommand. Use a [FlagSet](https://golang.org/pkg/flag/#FlagSet) for each subcommand. Add different flags for the different operations. Parse the commands with a simple switch case. To retrieve the command and all arguments, use the `os` library and its `Args` variable.   

### Solution
Go to the folder `go-calc` and follow the steps of the [README](../../go-calc/README.md) if you need further help. The code in the folder works out-of-the-box and you can compare your solution with it.
Depending on your knowledge level and how fast you are, you can additionally go through the `Testing Go Applications` and `Building Go Applications` parts. 

## Complex CLIs need a Cobra (90 min)

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
`kubectl chaos nginx-deployment --replicas 2 --namespace default` 
Where `replicas` is the amount of pods to delete and `namespace` the namespace in which the deployment lies. 

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

#### Add flag parsing
Start by adding the needed argument flags to the generated root command:
- `--replicas 2` to specify the number of replicas
- `--namespace default` to specify the K8s namespace to use

Check the [documentation about flags](https://github.com/spf13/cobra#working-with-flags) on how to add flags to your root command.

#### Add Config File Parsing with Viper
Add config file parsing to your CLI using Viper. The user of the CLI can set the namespace or any other flag parameters in a configuration file. Integrate it into your root command with `cobra.OnInitialize(initConfig)` and check the [Viper documentation](https://github.com/spf13/viper#reading-config-files) on how to read config files.

#### Add Kubernetes Client API

Next, we need to add the Kubenetes client API, in order to query a deployment and
delete the matching amount of pods for this deployment.

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

### Solution
Go to the subfolder /kubectl-chaos of this repo and follow the steps of the [README](../../kubectl-chaos/README.md) if you need further help. The code in the folder works out-of-the-box and you can compare your solution with it.
