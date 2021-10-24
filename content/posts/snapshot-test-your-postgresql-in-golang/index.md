+++
date = "2021-10-20T14:28:23+07:00"
author = "nvcnvn"
description = "Write unit-test for your repository layer using query snapshot."
title = "Snapshot test your Postgresql in Golang pgx driver"
categories = ["DevSecOps", "Testing"]
tags = ["Postgresql", "unit-test", "Golang"]
slug = "snapshot-test-your-postgresql-in-golang-pgx-driver"
+++

If your unit-test understand SQL syntax, you can cover many behaviors of your code without the need of starting a real DB. 
Want to understand Postgresql syntax? just simply import Postgresql parser to your program, 
https://github.com/pganalyze/pg_query_go helps you to do that.  
#### How about using a containerized DB?
Sometimes, the repository layer doesn't have much logic (maybe concat some `WHERE` conditions), only propagate the SQL 
statement to DB.  

Testing with a simple containerized DB is a good option where you can:
* verify the syntax, and
* checking the logic, where the real complexity happen with all the fake data and logic

Good option? Yes but not the best option for all the cases. For example, the [UI snapshot test](https://jestjs.io/docs/snapshot-testing)   technique has its place whenever we don't want to start a real browser for testing your UI code.
*Snapshot* test your SQL statement can have their place also. Hopefully, I can show you some of its advantages.

#### A better tool for you
Before going further, https://github.com/cockroachdb/copyist is a good tool for you if you're using the std sql interface. 
We're using pgx custom interface so we need to write some testing helper but the idea is the same. 

#### OK, now the snapshot test
In this post, we use a simple example query but you can see the real benefit with complex query that make you hate your 
ORM lib.  
Here is our example:
```go
type UsersInvoicesFilter struct {
	IDs       []string
	CreatedAt time.Time
}

const FetchUsersInvoicesStmt = `SELECT
	invoices.invoice_id,
	invoices.user_id,
	users.billing_address,
	invoices.amount
FROM users LEFT JOIN invoices ON users.user_id = invoices.user_id
WHERE users.user_id = ANY($1) AND invoices.created_at >= $2`

func FetchUsersInvoices(ctx context.Context, tx pgx.Tx, filter UsersInvoicesFilter) ([]*Invoice, error) {
	rows, err := tx.Query(ctx, FetchUsersInvoicesStmt, &filter.IDs, &filter.CreatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []*Invoice
	for rows.Next() {
		invoice := Invoice{}
		err := rows.Scan(
			&invoice.InvoiceID,
			&invoice.UserID,
			&invoice.BillingAddress,
			&invoice.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning invoice data error: %w", err)
		}
		invoices = append(invoices, &invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	return invoices, nil
}
```
All seasoned Go developers must see this code block familiar, we introduce some examples that unit-test can help to cover.
##### Validating PostgresSQL syntax
* increase code coverage percentage just to be looking good :shit:
* avoid "accident change" when refactoring :+1:
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
Now, import `pg_query "github.com/pganalyze/pg_query_go/v2"` and then test the basic syntax with
```go
func Test_parsingUsersInvoicesQuery(t *testing.T) {
	json, err := pg_query.ParseToJSON(FetchUsersInvoicesStmt)
	if err != nil {
		panic(err)
	}

	t.Log(json)
}
```
Then run your test
```
go test -v
--- FAIL: Test_parsingUsersInvoicesQuery (0.00s)
    users_repo_test.go:46: syntax error at or near "WHERE"
    users_repo_test.go:49: 
FAIL
exit status 1
FAIL    examples        0.005s
```
Oops, what's wrong with my `WHERE` condition? I think you already know I intentionally put the wrong join condition 
and with a short query like this, it's very easy to spot. But I guess at some time you must have written something 
more complex than CRUD queries, this is where it can be helpful. Besides that, this test returns the result in a small 
fraction of a second.  
Correcting the `JOIN invoices WHERE` to `JOIN invoices ON` and returning the test, a formatted version of the JSON can be 
found here https://gist.github.com/nvcnvn/d73d441b7878c47e85a654eef61819db  
The output shows how Postgresql parses and manages the SQL syntax tree, not something trivial in general but in our case 
you can compare the tree with the original query since it is fairly simple.
##### Not only text compare but understand SQL semantic
Sometimes, adding a new line, space, and tab does not change the semantic of the query, we have a small helper for this issue.
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
Because this is a legit parser, it comes with a location of each token, we don't care about the location. 
https://github.com/kinbiko/jsonassert is with `"<<PRESENCE>>"` is very handy for this case. Using together with a 
simple regex replaces all, we can have this JSON "template" stored for comparing later, the flow is:
* run the test the first time, record the JSON string
* modify all the location fields to check their presence only (ignore the value)
* then for each later test we use the query passed and compare it with the stored template

https://github.com/manabie-com/manabie-com.github.io/blob/main/content/posts/snapshot-test-your-postgresql-in-golang/examples/helper_test.go#L33

##### avoid mistake when refactor
Assuming this code working already:
```go
type UsersInvoicesFilter struct {
	IDs       []string
	CreatedAt time.Time
}

func FetchUsersInvoices(ctx context.Context, tx pgx.Tx, filter UsersInvoicesFilter) ([]*Invoice, error) {
	rows, err := tx.Query(ctx, FetchUsersInvoicesStmt, &filter.IDs, &filter.CreatedAt)
```

One day, you come up with an idea that we should have a generic interface for the filter in general, something like:
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
See the issue?  
I have misused `filter.Args()`, the correct usage should be `filter.Args()...`. This mistake can be avoided by unit-test. 
The idea is understand the SQL tree structure, then let it check if the number of args send to 
`tx.Query` method matches the number of `WHERE` conditions.  
https://github.com/tidwall/gjson can be used for this, example:
```go
	actual, expected := len(m.queryArgsLogs[0]),
		len(gjson.Get(m.statementLogs[0], "stmts.0.stmt.SelectStmt.whereClause.BoolExpr.args").Array())
	if actual != expected {
		m.t.Errorf("WHERE conditions have %d but function call send %d", expected, actual)
	}
```
We have the mock written in [helper_test.go](https://github.com/manabie-com/manabie-com.github.io/blob/main/content/posts/snapshot-test-your-postgresql-in-golang/examples/helper_test.go) for recording the Query call args.  

Let run the test with the bad code:
```
go test -v       
=== RUN   Test_FetchUsersInvoices
    helper_test.go:68: WHERE conditions have 2 but function call send 1
--- FAIL: Test_FetchUsersInvoices (0.01s)
FAIL
exit status 1
FAIL    examples        0.013s
```
##### more than just increasing code coverage percentage
With a fairly small amount of test code (not counting the helper that can be reused by the whole team), we can have some 
code coverage already, without the need of starting any docker container.  
```
go test -v -cover
=== RUN   Test_FetchUsersInvoices
--- PASS: Test_FetchUsersInvoices (0.01s)
PASS
coverage: 53.3% of statements
ok      examples        0.012s
```
But the Manabie backend team use these kinds of test as an automatic code-review, for example when we have a rule that 
all the queries touching a table need to have a special filter, this comes to handy since we can code that rule to the 
`MockDB` to check it.   
I think this is post long enough, please refer to our example folder https://github.com/manabie-com/manabie-com.github.io/tree/main/content/posts/snapshot-test-your-postgresql-in-golang/examples 
and try it. Maybe write a unit test to check the number of `Scan` args match with the `SELECT` target.  
