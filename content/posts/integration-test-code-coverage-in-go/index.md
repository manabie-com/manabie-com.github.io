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

Code coverage is a handy measurement. It's useful to help know your code is being tested, which helps you to know your code is working. It can also be used as a check before github pull requests can be merged.

In Golang, code coverage is easy to generate when testing packages, simply test them with `-cover` flags.

It's a bit trickier in integration tests because we have to exit all of the go services to make them generate coverage reports.

The steps we ended up taking to find the final percent was to  

- compile services with `go test -c`,
- add an http killswitch
- start test services & test them
- stop the services by calling the http killswitch endpoint
- collect coverage reports from containers
- merge them with gocovmerge.
- get a shiny code coverage percentage.


## Part I: Compiling services with go test instrumentation

We want to run our service with `go test` so we can use it's `cover` flags to enable code coverage output.

To do this, we need to first create a test function that works like our main, so that we can use `go test` to run our server.

That is to say, take:

```go
func main() {
    // start the service
}
```

and make it:


```go
// main.go

func main() {
    runService()
}

func runService() {
  // start the service.
}
```

```go
// main_test.go

func TestRun(t *testing.T) {
    runService()
}
```

For the coverage to be output, TestRun function needs to finish. We can't just kill the process, or go test will not output a coverage profile.

So we can set up an HTTP server, because it's really easy in Go. It will listen for a request which is our signal to let the function complete. Later, we can run curl to tell our services to stop.

Our simple HTTP server calls a context's CancelFunc when it receives a request.

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

	fmt.Println("Started KillServer")

	err := s.server.ListenAndServe()

	if err == http.ErrServerClosed {
		fmt.Println("KillServer Closed")
	} else {
		fmt.Println("KillServer Error:", err)
	}
}

func (s *killServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

  // cancel the context
	s.cancel()
}

```

And our TestRun uses the context to exit once the context is canceled.

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

Next, we can compile the service with: `go test -c ./ -cover -covermode=count -coverpkg=./...` to create a `.test` binary, and run it with:

```bash
./my-service.test -test.coverprofile my-service.out
# server started...
```

Now our server is running as it normally would, but with go's test tooling built in. All of the code it runs is counted for code coverage now.

In another terminal, we run some tests and then kill it with the HTTP signal.

```bash
./run-tests.sh

curl localhost:19999 # kills the server in the other terminal
```

A file named `my-service.out` should be created with coverage information from our run.


Using `go test -c` to build test binaries is a nice little trick I learned about recently. I suppose we could accomplish the same thing by just running all our servers with the `go test` command directly, instead of compiling them into binary, but then we would need go installed in our containers.

## Part 2: Running Things

Now that our services are built in test mode, we can add them to our Dockerfile, with a little wrapper script to restart them endleslessly.

```Dockerfile
FROM debian

WORKDIR /
# coverage test
COPY ./server.test /server.test
COPY ./server_with_restart.sh /server_with_restart.sh

ENTRYPOINT [ "/server_with_restart.sh" ]
```

server_with_restart.sh

```sh
#!/bin/sh

while true; do
    echo "Test server started. Built with 'go test -cover -c'"

    /server.test \
        -test.coverprofile=cover.out \
        "$@" || exit 1; # exit if process exited with error
    
    # Server was probably stopped by http kill signal
    # to collect code coverage. Restart.        
    echo "Server restarting.."
done
```

Now start those containers, and run your tests.


## Part 3: Collecting Coverage

Now to generate, download and merge our coverage.

To trigger generation, we just have to send an HTTP request to our killserver on port 19999. Not too hard, we can just install curl in our kubernetes container, and run it from the outside with exec.

```bash
kubectl exec my-service -- curl -s localhost:19999
```

Then copy the file out of the container:

```bash
kubectl exec -i my-service -- cat /cover.out > cover/my-service.out
```

Do this for each service, and then merge the coverage profiles together.

```bash
go install github.com/wadey/gocovmerge@latest

gocovmerge cover/*.out > cover/merged.cov

# Output Total Coverage 
go tool cover -func=cover/merged.cov | grep -E '^total\:' | sed -E 's/\s+/ /g'
```

Note: We use a `.cov` extention on the `merged.cov` to make it easy to use `cover/*.out` in scripting.

And that's about it. For us at Manabie, we are using this final coverage percentage in our Github Actions pipeline to ensure that our integration test coverage does not drop. 

Anwyays, that's all for now folks.
