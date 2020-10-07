# kubectl-choas Plugin using Cobra

This demo shows how to create a Kubernetes plugin by using a Cobra CLI app to 
perform chaos monkey style operations on Kubernetes deployments and pods.

Example usage:
```bash
$ kubectl apply -f nginx-deployment.yaml
$ kubectl chaos nginx-deployment --replicas 2 --namespace default

$ kubectl chaos --help
$ kubectl chaos --version
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
Open the generated `cmd/root.go` file, perform some cleanup.

## Exercise: Add Cobra command line flags

Next, add support for the following command line flags:
- `--replicas 2` to specify the number of replicas
- `--namespace default` to specify the K8s namespace to use
- `--version` to output version information

```golang
rootCmd.Flags().IntVarP(&replicas, "replicas", "r", 1, "Number of replicas")
rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "the namespace to use")
rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Display version info")
```

## Exercise: Add Kubernetes Client API

Next, we need to add the Kubenetes client API, in order to query of a deployment and then
delete the matching amount of pods for this deployment.

```golang
import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)
```

To connect to Kubernetes, the following API calls are required:
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

To query the deployment and delete pods, the following API calls are required:
```golang
    deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), args[0], metav1.GetOptions{})
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    labelMap := deployment.Spec.Selector.MatchLabels

    pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
        LabelSelector: labels.SelectorFromSet(labelMap).String(),
    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    for i, p := range pods.Items {
        if i >= replicas {
            break
        }
        fmt.Println("Deleting pod", p.Name)
        clientset.CoreV1().Pods(namespace).Delete(context.TODO(), p.Name, metav1.DeleteOptions{})
    }
```

## References

- https://github.com/kubernetes/client-go
- https://github.com/kubernetes/client-go/tree/master/examples/out-of-cluster-client-configuration
- https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration
