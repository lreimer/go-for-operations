# Challenge 2 - Kubernetes sidecars with Go

In this challenge, we will write a sidecar for Kubernetes in Go. A sidecar is a container that is a supplement to your application container inside your pod. The sidecar shares storage volumes, networking and other things with the application container, thus it can be used for intercepting traffic, sharing files with the application or adding monitoring. 

What we will build is a log shipping sidecar. It will read out the log files from the application by mounting the sidecar to the same volume as the application and transferring the images. 

## Exercise: Build a log shipping sidecar in Go (90 min)

### Initialize the Go module

Create a go module and the code files we need, running the following commands in this challenge folder:

```bash
go mod init k8s-logship-sidecar
touch main.go
touch config.go
```

### Configuration via environment variables
We build a simple log shipping sidecar that is configured using the following environment variables:
```
LOG_DIRECTORY # directory where logs lay / example:/logs
LOG_FILE_PATTERN # file pattern of log files / example: .+.log
LOG_SCAN_INTERVAL # scanning interval in seconds / example: "10"
``` 

Your code needs to read the environment variables and store them in Go variables. Put the config parsing in `config.go`. Use the `os` library to lookup environment, choose appropriate Go types for the variables. For the scan interval, use `time.Duration`. The log directory should be created if it does not exist. The file pattern should be stored as regex - use the `regexp` library. Store these variables as global variables. 

One example, the LOG_DIRECTORY environment variable, to help you get started:
```golang
// Directory to watch for changes
var Directory = logDirectory()

func logDirectory() string {
	name, ok := os.LookupEnv("LOG_DIRECTORY")
	if !ok || len(name) == 0 {
		log.Fatal("LOG_DIRECTORY environment variable not set.")
	}

	// Create the directory if it does not exist
	if _, err := os.Stat(name); os.IsNotExist(err) {
		os.Mkdir(name, os.ModePerm)
	}

	return name
}
```

Define the other two in the same way by creating a var which gets the method assigned that reads out the environment variable. 

### The main execution

Create your main function that contains the log shipping loop.
Use [time.Tick](https://golang.org/pkg/time/#Tick) with your scan interval variable as ticker (look at the example for sample usage) and call your logfile scanning method in the ticker loop. 
The logfile scanning method should read all files in the log directory and print the file content of files that match the file pattern. For listing the files of a directory you can use `ioutil.ReadDir`. When looping through the files check if its a file using `file.IsDir` and use the pattern regex to find out if the pattern matches. To print the contents you can use [ioutil.ReadFile](https://golang.org/pkg/io/ioutil/#example_ReadFile).

A log shipper should have graceful shutdown implemented. To do this in Go, we will use a Go routine that listens to the OS termination signal. 
Integrate the following code in your main function:

```golang
	go func() {
    // create channel with os.Signal type as listener for the SIGTERM and SIGINT syscalls
		var gracefulStop = make(chan os.Signal, 1)
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)

    // blocking wait for a termination signal
		sig := <-gracefulStop
		log.Println("Received stop signal:", sig)

    // bonus: add handling here for sending remaining logs...
		os.Exit(0)
	}()
```

The tick loop is implemented as followed:
```golang
	tick := time.Tick(ScanInterval)
	for range tick {
		scanForLogfiles()
	}
```
The `time.Tick` method returns an endless channel which sends a message on the channel on every interval.

Now implement the `scanForLogfiles()` method.

```golang
func scanForLogfiles() {
	log.Printf("Scanning for files %v in %v", FilePattern, Directory)
	files, err := ioutil.ReadDir(Directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
    // for every file:
    // check if its a real file using `f.IsDir` and if it matches the FilePattern (try the MatchString method):
    // read out the contents using `ioutil.ReadFile and print them to the stdout
		// ...
	}
}

```


If you have problems, check the solution in the `k8s-logship-sidecar` in the root of the repo.

### Dockerize the sidecar

Do deploy the sidecar alongside an application container, we need to build a docker container.
Let's look at a multi-stage Dockerfile:
```
FROM golang:1.15.2 as builder # start with golang container as build image

WORKDIR /build

COPY . /build
RUN make build # run make build to build the sidecar

FROM gcr.io/distroless/static-debian10 # actual runtime base container  
COPY --from=builder /build/k8s-logship-sidecar / # we copy the created binary to the runtime image

# set default environment variables
ENV LOG_DIRECTORY=/logs
ENV LOG_FILE_PATTERN=.+.gz
ENV LOG_SCAN_INTERVAL=10

# configure entrypoint
ENTRYPOINT ["/k8s-logship-sidecar"]
CMD [""]
```

The advantage of the multi-stage docker build is that we have a smaller runtime image as the golang build tools are only needed for building the binary. 

### Building the docker container

Use the provided Dockerfile inside the `workshop/challenge-2` folder. 

```bash
# to run the operator in minikube we need to configure the docker daemon to use the minikube context
eval $(minikube -p minikube docker-env)

# create the docker container
docker build -t k8s-logship-sidecar:v1.0.0 .
```

### Deployment of log shipping sidecar with nginx 

For testing, we need to deploy the nginx again but with the log shipping sidecar.
Let's look at the deployment yaml:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
        # need to add a volume mount for the nginx logs to share them with the logshipper
        volumeMounts:
          - name: applogs
            mountPath: /var/log/nginx
    # log shipper added as additional container
      - name: logshipper
        image: k8s-logship-sidecar:v1.0.0
        # configuration of environment variables
        env:
          - name: LOG_DIRECTORY
            value: /logs
          - name: LOG_FILE_PATTERN
            value: .+.log
          - name: LOG_SCAN_INTERVAL
            value: "10"
        # mounting the volume under /logs as specified in the LOG_DIRECTORY environment variable
        volumeMounts:
        - name: applogs
          mountPath: /logs
          readOnly: true
    # create volume as emptyDir which is a temporary directory on the host 
      volumes:
        - name: applogs
          emptyDir: {}
```

Deploy it with the following command (use the supplied nginx-deployment.yaml):

```bash
kubectl apply -f nginx-deployment.yaml
```

Now, you can check if the nginx logs are printed as logs of the sidecar with the following command:

```bash
kubectl logs <nginx pod> -c logshipper -f
```

Try making calls against the nginx to see if there is access logs:
```bash
curl $(minikube ip):$(kubectl get service nginx-service -o jsonpath="{.spec.ports[0].nodePort}")
```
