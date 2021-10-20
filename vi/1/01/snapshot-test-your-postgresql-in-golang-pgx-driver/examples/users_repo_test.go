package examples

import (
	"context"
	"testing"
)

func Test_FetchUsersInvoices(t *testing.T) {
	mockDB := &MockDB{
		// toggle recording to true for created new snapshot,
		// then use it with assertStatementWithTemplate

		// recording: true,
		t: t,
	}
	FetchUsersInvoices(context.Background(), mockDB, &UsersInvoicesFilter{})

	mockDB.assertNumberOfSQueryArgs("stmts.0.stmt.SelectStmt.whereClause.BoolExpr.args")
	mockDB.assertStatementWithTemplate("fa015d1a59499f78.rec.json")
}
