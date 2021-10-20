package examples

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v4"
)

type Invoice struct {
	InvoiceID      string
	UserID         string
	BillingAddress string
	Amount         int
}

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

const FetchUsersInvoicesStmt = `SELECT
	invoices.invoice_id,
	invoices.user_id,
	users.billing_address,
	invoices.amount
FROM users LEFT JOIN invoices ON users.user_id = invoices.user_id
WHERE users.user_id = ANY($1) AND invoices.created_at >= $2`

func FetchUsersInvoices(ctx context.Context, tx pgx.Tx, filter Filter) ([]*Invoice, error) {
	rows, err := tx.Query(ctx, FetchUsersInvoicesStmt, filter.Args()...)
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
