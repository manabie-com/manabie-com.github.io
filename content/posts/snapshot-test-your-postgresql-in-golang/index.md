+++
author = "nvcnvn"
description = "Write unit-test for you repository layer using query snapshot."
title = "Snapshot test your Postgresql in Golang pgx driver"
categories = ["DevSecOps", "Testing"]
tags = ["Postgresql", "unit-test", "Golang"]
slug = "snapshot-test-your-postgresql-in-golang-pgx-driver"
+++

#### How about using a containerized DB?
In some applications, the repository layer don't have much logic (maybe concat some `WHERE` conditions), only 
propagate the SQL statement to DB - where the real complexity happen with all the fake data and logic.  
Testing with a simple containerized DB is a good option where you can:
* verify the syntax, and
* checking the logic

Good option? Yes but not the best option for all the cases. For example the [UI snapshot test](https://jestjs.io/docs/snapshot-testing)   technique have their place whenever we don't want to start a real browser for testing your UI code, "snapshot test" 
your SQL statement can have their place also. I hope can show you some of it advantage of it.

#### A better tool for you
Before going further, I think https://github.com/cockroachdb/copyist is a good tool for you if you're using the std sql 
interface. In Manabie, we're using pgx custom interface so we need to write some testing helper but the idea is the same. 

#### OK, now the snapshot test
The goal is making SQL syntax "easier" to work with in your unit-test code. This is where I love Postgresql and the open 
source community, you want to understand Postgresql syntax? just simply import Postgresql parser to your program, 
https://github.com/pganalyze/pg_query_go help you to do that.
##### A lazy fully snapshot test
* increase code coverage percentage just too be looking good :shit:
* avoid "accident change" when refactor :+1:
* check Postgresql syntax :+1:

Example use case: a function to fetch all user's invoices. In this example we have 2 tables:
* users(**user_id**, email, billing_address)
* invoices(**invoice_id**, user_id, amount, created_at)

A single query to fetch the users and their invoice data with 2 conditions:
```sql
SELECT
	invoices.invoice_id,
	invoices.user_id,
	users.billing_address,
	invoices.amount
FROM users JOIN invoices WHERE users.user_id = invoices.user_id
WHERE users.user_id = ANY($1) AND invoices.created_at >= $2
```
Now, import `pg_query "github.com/pganalyze/pg_query_go/v2"` and then testing the basic syntax with
```go
func Test_parsingUsersInvoicesQuery(t *testing.T) {
	json, err := pg_query.ParseToJSON(FetchUsersInvoicesStmt)
	if err != nil {
		panic(err)
	}

	t.Log(json)
}
```
Then run you test
```
go test -v
--- FAIL: Test_parsingUsersInvoicesQuery (0.00s)
    users_repo_test.go:46: syntax error at or near "WHERE"
    users_repo_test.go:49: 
FAIL
exit status 1
FAIL    examples        0.005s
```
Oops, whats wrong with my `WHERE` condition? I think you already know I intention put the wrong join condition, and with 
a short query like this its very easy to spot. But I guess at some time you must have written something more complex 
than a CRUD queries, this is where it can be helpful, beside that, this test return the result in a small fraction of a 
second.  
Correcting the `JOIN invoices WHERE` to `JOIN invoices ON` and return the test, a formatted version can be found here 
https://gist.github.com/nvcnvn/d73d441b7878c47e85a654eef61819db  
The output show how Postgresql pare and managing the SQL syntax tree, not some thing trivial in general but in our case 
you can compare the tree with our original query since it fairly simple.
##### More than a text compare
Some time, adding new line, space and tab not changing the sematic of the query, we have small helper for this issue.
```json
                           "ColumnRef":{
                              "fields":[
                                 {
                                    "String":{
                                       "str":"users"
                                    }
                                 },
                                 {
                                    "String":{
                                       "str":"billing_address"
                                    }
                                 }
                              ],
                              "location":49
                           }
```
Because this is a legit parser, it come with a location of each token, we don't really care about the location. 
https://github.com/kinbiko/jsonassert is with `"<<PRESENCE>>"` is very handy for this case. Using together with a 
simple regex replace all, we can have this JSON "template" stored for compare later, the flow is:
* run the test the first time, record the JSON string
* modify all the location field to check their presence only (ignore the value)
* then for each later test we use the query passed and compare with the stored template

##### avoid messing with "...interface{}"
Assuming this code working already:
```go
type UsersInvoicesFilter struct {
	IDs       []string
	CreatedAt time.Time
}

func FetchUsersInvoices(ctx context.Context, tx pgx.Tx, filter UsersInvoicesFilter) ([]*Invoice, error) {
	rows, err := tx.Query(ctx, FetchUsersInvoicesStmt, &filter.IDs, &filter.CreatedAt)
```

One day, you come up with an idea that we should have a generic interface for filter in general, something like:
```go
type Filter interface {
	Args() []interface{}
}

type UsersInvoicesFilter struct {
	IDs       []string
	CreatedAt time.Time
}

func (f *UsersInvoicesFilter) Args() []interface{} {
	return []interface{}{&f.IDs, &f.CreatedAt}
}
```
and then, to use it
```go
// FetchUsersInvoices(context.Background(), mockDB, &UsersInvoicesFilter{})
func FetchUsersInvoices(ctx context.Context, tx pgx.Tx, filter Filter) ([]*Invoice, error) {
	rows, err := tx.Query(ctx, FetchUsersInvoicesStmt, filter.Args())
```
See the issue? I think you see it now because we're on a post talking about testing. But let write a test for a bad day. 
The idea is since our unit test can understand the tree structure, then we let it check if the number of args send to 
`tx.Query` method matching the number of `WHERE` conditions. https://github.com/tidwall/gjson can be use for this.  
Short example:
```go
	actual, expected := len(m.queryArgsLogs[0]),
		len(gjson.Get(m.statementLogs[0], "stmts.0.stmt.SelectStmt.whereClause.BoolExpr.args").Array())
	if actual != expected {
		m.t.Errorf("WHERE conditions have %d but function call send %d", expected, actual)
	}
```
We have the mock written in `helper_test.go` for recording the Query call args.