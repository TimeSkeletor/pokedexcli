package pokedatabase

import (
	"fmt"
	"log"
	"strings"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

const (
	dbName = "pokedex.db"
)

func SetDb() {
	if err := withConnTxFn(createTables, getTables()); err != nil {
		log.Fatalf("failed to setup DB: %v", err)
	}
}

func getConn() (*sqlite.Conn, error) {
	conn, err := sqlite.OpenConn(dbName, sqlite.SQLITE_OPEN_CREATE|sqlite.SQLITE_OPEN_READWRITE)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := sqlitex.Exec(conn, "BEGIN;", nil); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return conn, nil
}

func withConnTxFn[T any](fn func(*sqlite.Conn, T) error, arg T) error {
	conn, err := getConn()
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Close()

	if err := fn(conn, arg); err != nil {
		_ = sqlitex.Exec(conn, "ROLLBACK;", nil)
		return err
	}
	return sqlitex.Exec(conn, "COMMIT;", nil)
}

func createTables(conn *sqlite.Conn, tables map[string]table) error {
	for name, table := range tables {
		stmt := prepareTable(name, table)
		err := sqlitex.ExecScript(conn, stmt)

		if err != nil {
			return fmt.Errorf("failed to execute script: %w", err)
		}
	}
	return nil
}

func prepareTable(name string, t table) string {
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", name)

	fields := []string{}
	for fieldName, field := range t.fields {
		fields = append(fields, fmt.Sprintf("    %s %s", fieldName, field.dataType))
	}

	stmt += strings.Join(fields, ",\n")
	stmt += "\n);"
	return stmt
}
