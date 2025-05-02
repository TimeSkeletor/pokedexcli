package pokedatabase

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

const (
	dbName = "pokedex.db"
)

type Database struct {
    pool *sqlitex.Pool
    mu   sync.Mutex // Protect pool initialization
}

func New() *Database {
    db := &Database{}
    if err := db.setDb(); err != nil {
        log.Fatalf("failed to setup DB: %v", err)
    }
    return db
}

func (db *Database) setDb() error {
    db.mu.Lock()
    defer db.mu.Unlock()

    if db.pool != nil {
        return nil // Pool already initialized
    }

    pool, err := sqlitex.Open(dbName, sqlite.SQLITE_OPEN_CREATE|sqlite.SQLITE_OPEN_READWRITE, 10)
    if err != nil {
        return fmt.Errorf("failed to open DB pool: %w", err)
    }
    db.pool = pool

    return withConnTxFn(db.pool, createTables, getTables())
}

func (db *Database) Close() error {
    db.mu.Lock()
    defer db.mu.Unlock()

    if db.pool != nil {
        db.pool.Close()
        db.pool = nil
    }
    return nil
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

func withConnTxFn[T any](pool *sqlitex.Pool, fn func(*sqlite.Conn, T) error, arg T) error {
    conn := pool.Get(nil)
    if conn == nil {
        return fmt.Errorf("failed to get connection from pool")
    }
    defer pool.Put(conn)

    if err := sqlitex.ExecTransient(conn, "BEGIN;", nil); err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    if err := fn(conn, arg); err != nil {
        sqlitex.ExecTransient(conn, "ROLLBACK;", nil)
        return err
    }

    return sqlitex.ExecTransient(conn, "COMMIT;", nil)
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
