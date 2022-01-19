# Proof of Concept for gocovmerge with two processes + killserver

This proof of concept is to show that we can start an http server with go test, and test it from the outside (like an integration test) to produce coverage results, and then merge two results together to find a final percent. 

For simplicity I just use one body of code with 2 endpoints that we will pretend are separate services. We will test each endpoint individually and produce the coverage, and merge it.

In reality you should modify the -coverpkg to include all code from all relevant services. Gocovmerge cannot merge different sets of code.


## Build It

```
go test -c ./ -cover -covermode=count -coverpkg=./...
```

## Get First Coverage Profile

```bash
ADDR=":1234" ./covertest.test -test.coverprofile covA.out &
# 2021/12/27 15:37:47 Starting server at :1234

curl -s http://localhost:1234/a
# A

curl -s http://localhost:1234/kill
# 2021/12/27 15:38:09 http: Server closed
# PASS
# coverage: 68.8% of statements in ./...
```

## Get Second Coverage Profile

```bash
ADDR=":1234" ./covertest.test -test.coverprofile covB.out &
# 2021/12/27 15:37:47 Starting server at :1234

curl -s http://localhost:1234/b
# B
# has
# 2 extra lines of code

curl -s http://localhost:1234/kill
# 2021/12/27 15:38:09 http: Server closed
# PASS
# coverage: 81.2% of statements in ./...
```

## Explanation

Our first handler has two fewer lines of code. We started the server, and ran endpoint a, ended the server. This generates a coverage profile that covers a handler, and startup code.

In our second test we test the b handler, which has two extra lines of code.

in A we get `68.8%` coverage, and in B we get `81.2%` coverage.

This is similar to having multiple services, except I'm lazy and used one main file.

## Merge the results

Using this cool package: https://github.com/wadey/gocovmerge

```bash
gocovmerge covA.out covB.out > mergedAB.cov
go tool cover -func=mergedAB.cov    
# covertest/main.go:13:	main		0.0%
# covertest/main.go:20:	Serve		100.0%
# covertest/main.go:36:	HandlerA	100.0%
# covertest/main.go:41:	HandlerB	100.0%
# covertest/main.go:48:	KillHandler	100.0%
# total:			(statements)	93.8%
```

We can also go check the other coverage profiles to verify that they each ran one handler:

```bash
go tool cover -func=covA.out        
# covertest/main.go:13:	main		0.0%
# covertest/main.go:20:	Serve		100.0%
# covertest/main.go:36:	HandlerA	100.0%
# covertest/main.go:41:	HandlerB	0.0%
# covertest/main.go:48:	KillHandler	100.0%
# total:			(statements)	68.8%
```

```bash
go tool cover -func=covB.out        
# covertest/main.go:13:	main		0.0%
# covertest/main.go:20:	Serve		100.0%
# covertest/main.go:36:	HandlerA	0.0%
# covertest/main.go:41:	HandlerB	100.0%
# covertest/main.go:48:	KillHandler	100.0%
# total:			(statements)	81.2%
```