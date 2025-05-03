package pokedatabase

import (
    "context"
    "fmt"
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
    mu   sync.RWMutex
    once  sync.Once
}

func (db *Database) SetConn() error {
    var err error
    db.once.Do(func() {
        db.mu.Lock()
        defer db.mu.Unlock()
        db.pool, err = sqlitex.Open(dbName, sqlite.SQLITE_OPEN_CREATE|sqlite.SQLITE_OPEN_READWRITE, 10)
    })
    return err
}
func (db *Database) CloseConn() error {
    db.mu.Lock()
    defer db.mu.Unlock()

    if db.pool != nil {
        err := db.pool.Close()
        db.pool = nil
        return err
    }
    return nil
}

func (db *Database) SetDb(ctx context.Context) error {
    if err := db.SetConn(); err != nil {
        return err
    }
    return db.WithTransaction(ctx, func(conn *sqlite.Conn) error {
        return db.createTables(conn)
    })
}

func (db *Database) ExecQuery(ctx context.Context, query string, resultFn func(stmt *sqlite.Stmt) error, args ...interface{}) error {
    if err := ctx.Err(); err != nil {
        return fmt.Errorf("context error: %w", err)
    }
    conn, err := db.getConn(ctx)
    if err != nil {
        return err
    }
    defer db.pool.Put(conn)
    if len(args) != strings.Count(query, "?") {
        return fmt.Errorf("query expects %d arguments, got %d", strings.Count(query, "?"), len(args))
    }
    return sqlitex.Exec(conn, query, resultFn, args...)
}


func (db *Database) WithTransaction(ctx context.Context, fn func(*sqlite.Conn) error) error {
    if err := ctx.Err(); err != nil {
        return fmt.Errorf("context error: %w", err)
    }
    conn, err := db.getConn(ctx)
    if err != nil {
        return err
    }
    defer db.pool.Put(conn)
    if err := sqlitex.ExecTransient(conn, "BEGIN;", nil); err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    if err := fn(conn); err != nil {
        sqlitex.ExecTransient(conn, "ROLLBACK;", nil)
        return err
    }
    return sqlitex.ExecTransient(conn, "COMMIT;", nil)
}

func (db *Database) getConn(ctx context.Context) (*sqlite.Conn, error) {
    db.mu.Lock()
    if db.pool == nil {
        db.mu.Unlock()
        return nil, fmt.Errorf("database pool not initialized; call SetConn first")
    }
    conn := db.pool.Get(ctx)
    db.mu.Unlock()
    if conn == nil {
        return nil, fmt.Errorf("failed to get connection from pool")
    }
    return conn, nil
}

func (db *Database) createTables(conn *sqlite.Conn) error {
    for name, table := range db.getTables() {
        stmt := prepareTable(name, table)
        if err := sqlitex.ExecScript(conn, stmt); err != nil {
            return fmt.Errorf("failed to execute script for table %s: %w", name, err)
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