package examples

import (
	"context"
	"io/fs"
	"os"
	"regexp"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/kinbiko/jsonassert"
	pg_query "github.com/pganalyze/pg_query_go/v2"
	"github.com/tidwall/gjson"
)

type MockDB struct {
	recording     bool
	t             *testing.T
	statementLogs []string
	queryArgsLogs [][]interface{}
	pgx.Tx
}

func (m *MockDB) Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error) {
	json, err := pg_query.ParseToJSON(sql)
	if err != nil {
		panic(err)
	}

	m.statementLogs = append(m.statementLogs, json)
	m.queryArgsLogs = append(m.queryArgsLogs, arguments)

	if m.recording {
		m.recordingStmtLog(sql, json)
	}

	return &MockRows{
		// this is for no data returned
		mockNextFunc: func() bool {
			return false
		},
		// scan need to check nothing for this case
		mockScanFunc: func(dest ...interface{}) error {
			return nil
		},
		// this is for no error returned
		mockErrFunc: func() error {
			return nil
		},
	}, nil
}

func (m *MockDB) recordingStmtLog(sql, parsedJSON string) {
	id, err := pg_query.Fingerprint(sql)
	if err != nil {
		panic(err)
	}

	recName := id + ".rec.json"
	os.WriteFile(recName, []byte(makeLocationFieldEasierToCompare(parsedJSON)), fs.FileMode(0644))
	panic("panic on recording mode, record file name: " + recName)
}

func (m *MockDB) assertNumberOfSQueryArgs(jsonPath string) {
	actual, expected := len(m.queryArgsLogs[0]),
		len(gjson.Get(m.statementLogs[0], jsonPath).Array())
	if actual != expected {
		m.t.Errorf("WHERE conditions have %d but function call send %d", expected, actual)
	}
}

func (m *MockDB) assertStatementWithTemplate(template string) {
	recContent, err := os.ReadFile("fa015d1a59499f78.rec.json")
	if err != nil {
		m.t.Fatal("error reading recorded statement", err)
	}

	ja := jsonassert.New(m.t)
	ja.Assertf(m.statementLogs[0], string(recContent))
}

type MockRows struct {
	pgx.Rows
	mockNextFunc func() bool
	mockScanFunc func(dest ...interface{}) error
	mockErrFunc  func() error
}

func (m *MockRows) Next() bool {
	return m.mockNextFunc()
}

func (m *MockRows) Scan(dest ...interface{}) error {
	return m.mockScanFunc(dest...)
}

func (m *MockRows) Err() error {
	return m.mockErrFunc()
}

func (m *MockRows) Close() {
	// no-op func
}

var locationPattern = regexp.MustCompile("location\":[0-9]+")

func makeLocationFieldEasierToCompare(j string) string {
	return locationPattern.ReplaceAllString(j, `location":"<<PRESENCE>>"`)
}
