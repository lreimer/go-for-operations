# K8s Microservice Operator

This operator aims to make operations a lot easier by abstracting the usual
`Deployment`, `Service` and `ConfigMap` definitions using a simple and unified
`Microservice` custom resource. The operator will then manage the underlying
Kubernetes resources automagically.

:exclamation: This operator is just a demo. Currently the operator only creates a basic deployment and service object.

```yaml
apiVersion: apps.qaware.de/v1
kind: Microservice
metadata:
  name: microservice-test
  labels:
    app: nginx
spec:
  replicas: 2
  image: nginx:1.17.6
  ports:
    - 80
  serviceType: LoadBalancer
```

# Creating the Operator Skeleton

To create the operator skeleton we will use the Operator SDK. Make sure it is installed
and available on your machine.

```bash
# brew install operator-sdk
$ mkdir k8s-microservice-operator
$ cd k8s-microservice-operator

$ operator-sdk init --project-version="2" --domain qaware.de --license none --owner "Mario-Leander Reimer" --plugins go.kubebuilder.io/v2 --repo github.com/lreimer/go-for-operations/k8s-microservice-operator
$ operator-sdk create api --group apps --version v1 --kind Microservice --resource=true --controller=true

$ curl https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/master/hack/setup-envtest.sh -o setup-envtest.sh
$ chmod +x setup-envtest.sh

$ make install
$ make docker-build docker-push
$ make deploy

$ k9s
$ kubectl apply -f config/samples/apps_v1_microservice.yaml
```

## Creating the Microservice CRD

Next we need to defined the CRD specification for the Microservice resource.
Once done, run `make generate` and `make manifests`

```golang
// MicroserviceSpec defines the desired state of Microservice
type MicroserviceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Minimum=0
	// Replicas is the number of replicas for the microservice deployment
	Replicas int32 `json:"replicas,omitempty"`

	// Image is the Docker image and tag to use for the microservice deployment
	Image string `json:"image,omitempty"`

	// Ports is the list of HTTP container ports for the microservice deployment
	Ports []int32 `json:"ports"`

	// ServiceType is the service type to use for the microservice service
	ServiceType string `json:"serviceType,omitempty"`
}
```

## Implement Reconcile Loop



## References

- https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/
- https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/
- https://sdk.operatorframework.io/docs/building-operators/golang/references/envtest-setup/
