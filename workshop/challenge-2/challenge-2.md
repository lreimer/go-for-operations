# Challenge 2 - Kubernetes sidecar with Go

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

### The main execution

Create your main function that contains the log shipping loop.
Use [time.Tick](https://golang.org/pkg/time/#Tick) with your scan interval variable as ticker (look at the example for sample usage) and call your logfile scanning method in the ticker loop. 
The logfile scanning method should read all files in the log directory and print the file content of files that match the file pattern. For listing the files of a directory you can use `ioutil.ReadDir`. When looping through the files check if its a file using `file.IsDir` and use the pattern regex to find out if the pattern matches. To print the contents you can use [ioutil.ReadFile](https://golang.org/pkg/io/ioutil/#example_ReadFile).

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

### Building and Running

Use the makefile and provided Dockerfile inside the `workshop/challenge-2` folder. 

```bash
$ make build
$ make docker

$ docker run -it -v `pwd`:/logs k8s-logship-sidecar:v1.0.0
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
$ kubectl apply -f nginx-deployment.yaml
```

Now, you can check if the nginx logs are printed as logs of the sidecar with the following command:

```bash
$ kubectl logs <nginx pod> -c logshipper
```