+++
date = "2021-11-05T17:10:49+07:00"
author = "ds0nt"
description = "After we build our microservices to collect code coverage and run integration tests on them, we merge our code coverage from the kubernetes containers to get a final coverage figure"
title = "Test Coverage of Go Services during Integration Tests"
categories = ["DevSecOps", "Testing"]
tags = ["Kubernetes", "Docker", "integration-test", "Golang", "microservices", "code-coverage]
slug = "integration-test-code-coverage-in-go"
+++
# Test Coverage of Go Services during Integration Tests


In Golang, getting code coverage with `go test` is easy. But it's still rather hard for integration tests. 

Here I want to introduce the method we used at Manabie to collect and measure code coverage on our servers from integration tests.

## About our integration tests

At Manabie we use Kubernetes for container orchestration. To perform integration tests, we deploy our services, and then run a test container with a go program with a whole lot of integration tests. 

On dev environments, we start up minikube, deploy the services, and then run the test container as well. In CI, we start a vcluster instead of a minikube.

Collecting coverage in this setting is a little more complicated, but doable. The main idea is to compile the services with `go test -cover` instead of `go build`, and then get the services to exit in a timely manner after testing is finished.


What's on the agenda:

Part 1:
- compile services with `go test -c`,
- add an http killswitch
Part 2:
- run services & run tests
- stop the services by calling the http killswitch endpoint
Part 3:
- collect coverage reports from containers
- merge them with gocovmerge.
- get a shiny code coverage percentage.


## Part I: Compiling services with go test instrumentation

We want to run our service with `go test` so we can use it's `cover` flags to enable code coverage output.

To do this, we need to first create a test function that works like our `func main()`, so that we can use `go test` to run our server.

For example:

```go
func main() {
    server.Run()
}
```

becomes


```go
func run() {
	server.Run().
}

func main() {
    run()
}

```

so that we can write a test like this that starts the server:

```go
func TestRun(t *testing.T) {
    run()
}
```

For the coverage to be output, TestRun function needs to finish. We can't just kill the process, or go test will not output a coverage profile.

Because the service does not know when the tests are complete, we have to set up a mechanism to stop it remotely. We can set up a simple HTTP server. When it gets a request, it will gracefully terminate our service. Later, we can call it with curl.

Our Kill HTTP server calls a context's CancelFunc when it receives a request.

```go
type killServer struct {
	server http.Server
	cancel context.CancelFunc
}

func newKillServer(addr string, cancel context.CancelFunc) *killServer {
	return &killServer{
		server: http.Server{
			Addr: addr,
		},
		cancel: cancel,
	}
}

func (s *killServer) Start() {
	s.server.Handler = s

	err := s.server.ListenAndServe()
	if err != nil {
		fmt.Println("KillServer Error:", err)
	}
}

func (s *killServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

  // cancel the context
	s.cancel()
}

```

Our TestRun uses the same context to run the service.

```go
func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
  
	killServer := newKillServer(":19999", cancel)
	go killServer.Start()  
	go runService(ctx)

	<-ctx.Done()

	killServer.server.Shutdown(context.Background())
}
```

To test out our killswitch locally, we can compile the service with: `go test -c ./ -cover -covermode=count -coverpkg=./...` to create a `.test` binary, and run it with:

```bash
./my-service.test -test.coverprofile my-service.out
```

In another terminal, we run some tests and then kill it with the HTTP signal.

```bash
./run-tests.sh # here we could make test calls to the server.

curl localhost:19999 # and then kill the remote server when were done
```

Finally, a file named `my-service.out` should be created with coverage information from our server.

Using `go test -c` to build test binaries is a nice little trick I learned about recently. It plays nicely with containers too.

## Part 2: Running in our Kubernetes Cluster

Now we can add them to our Dockerfile. Also, we're going to add a little wrapper script to restart them endlessly. If we let the entrypoint process end, our container and coverage files will be deleted after the killswitch is called, which is not what we want.


```sh
#!/bin/sh
# server_with_restart.sh

# this script passes arguments into server.test with
# the addition of -test.coverprofile=cover.out

while true; do
    echo "Service started in coverage mode"
    /server.test -test.coverprofile=cover.out "$@" || exit 1;
        
    echo "Server restarting.."
done
```

We need a Dockerfile with the `server.test` binary, `curl`, and the restart script.

```Dockerfile
FROM alpine

WORKDIR /

RUN apk --no-cache add curl

COPY ./server_with_restart.sh /server_with_restart.sh
COPY ./server.test /server.test

ENTRYPOINT [ "/server_with_restart.sh" ]
```

Next, build the new docker image and deploy it in your kubernetes cluster, possibly with helm. Make sure you don't so override the Dockerfile's entrypoint. However, you can pass args as usual into the deployment container, which will pass into the server.test command.

Next, run your tests. In our case, we built an integration testing program named `gandalf` that we run in our cluster:

```
helm install gandalf ./deployments/gandalf
kubectl exec -it gandalf -- /gandalf-tests
```

Coverage will be recorded by the server.test.

## Part 3: Collecting Coverage

Now to generate, download and merge our coverage.

To output the coverage we just have to send an HTTP request to our killserver on port 19999.

```bash
kubectl exec my-service -- curl -s localhost:19999
```

Then copy the file to our local filesystem:

```bash
kubectl exec -i my-service -- cat /cover.out > cover/my-service.out
```

Do this for each service being tested, and then merge the coverage profiles together.

```bash
go install github.com/wadey/gocovmerge@latest

gocovmerge cover/*.out > cover/merged.cov

# Output Total Coverage 
go tool cover -func=cover/merged.cov | grep -E '^total\:' | sed -E 's/\s+/ /g'
```

Note: We use a `.cov` extention on the `merged.cov` to make it easy to use `cover/*.out` in scripting.

Inspect the contents of your new `merged.cov` file and try not to get drunk with power.

### Conclusion

Now our build pipeline prints out our integration test coverage percentage collected from all our go services, and combined into a final figure. 

For us at Manabie, we output this percent in our CI logs, and we set a rule that prevents pull requests from being mergeable if they decrease this percent.