package pokedatabase_test

import (
	"context"
	"os"
	"testing"

	"crawshaw.io/sqlite"
	"github.com/timeskeletor/pokedexcli/internal/pokedatabase"
)

func cleanupDB() {
    _ = os.Remove("pokedex.db")
}

func TestSetConnAndCloseConn(t *testing.T) {
    db := &pokedatabase.Database{}
    cleanupDB()
    defer cleanupDB()

    if err := db.SetConn(); err != nil {
        t.Fatalf("SetConn failed: %v", err)
    }
    if db.SetConn() != nil {
        t.Log("Calling SetConn again should not fail")
    }
    if err := db.CloseConn(); err != nil {
        t.Fatalf("CloseConn failed: %v", err)
    }
}

func TestSetDbCreatesPokemonTable(t *testing.T) {
    db := &pokedatabase.Database{}
    cleanupDB()
    defer cleanupDB()

    ctx := context.Background()
    if err := db.SetDb(ctx); err != nil {
        t.Fatalf("SetDb failed: %v", err)
    }

    // Check that the pokemon table exists by selecting from it
    err := db.ExecQuery(ctx, `SELECT COUNT(*) FROM pokemon;`, func(stmt *sqlite.Stmt) error {
        stmt.ColumnInt(0) // Just to access the column
        return nil
    })
    if err != nil {
        t.Fatalf("pokemon table does not exist or failed to query: %v", err)
    }

    _ = db.CloseConn()
}

func TestExecQuery_InsertAndSelect_Pokemon(t *testing.T) {
    db := &pokedatabase.Database{}
    cleanupDB()
    defer cleanupDB()

    ctx := context.Background()
    if err := db.SetDb(ctx); err != nil {
        t.Fatalf("SetDb failed: %v", err)
    }

    insert := `INSERT INTO pokemon (number, name, gender, capture_rate, base_happiness, is_baby, is_legendary, is_mythical, is_shiny, sprite_path) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
    err := db.ExecQuery(ctx, insert, nil, 25, "Pikachu", 1, 190, 70, false, false, false, false, "/sprites/pikachu.png")
    if err != nil {
        t.Fatalf("ExecQuery insert failed: %v", err)
    }

    var gotName string
    query := `SELECT name FROM pokemon WHERE number = ?;`
    err = db.ExecQuery(ctx, query, func(stmt *sqlite.Stmt) error {
        gotName = stmt.ColumnText(0)
        return nil
    }, 25)
    if err != nil {
        t.Fatalf("ExecQuery select failed: %v", err)
    }
    if gotName != "Pikachu" {
        t.Errorf("Expected Pikachu, got %s", gotName)
    }

    _ = db.CloseConn()
}

func TestWithTransactionRollback_Pokemon(t *testing.T) {
    db := &pokedatabase.Database{}
    cleanupDB()
    defer cleanupDB()

    ctx := context.Background()
    if err := db.SetDb(ctx); err != nil {
        t.Fatalf("SetDb failed: %v", err)
    }

    // Insert one Pokémon before the transaction
    err := db.ExecQuery(ctx, `INSERT INTO pokemon (number, name) VALUES (?, ?);`, nil, 133, "Eevee")
    if err != nil {
        t.Fatalf("Insert before tx failed: %v", err)
    }

    // Intentionally fail inside transaction to trigger rollback
    _ = db.WithTransaction(ctx, func(conn *sqlite.Conn) error {
        // Insert one
        err := db.ExecQuery(ctx, `INSERT INTO pokemon (number, name) VALUES (?, ?);`, nil, 150, "Mewtwo")
        if err != nil {
            return err
        }
        // Duplicate number (133) causes unique constraint violation
        return db.ExecQuery(ctx, `INSERT INTO pokemon (number, name) VALUES (?, ?);`, nil, 133, "DuplicateEevee")
    })

    var count int
    _ = db.ExecQuery(ctx, `SELECT COUNT(*) FROM pokemon;`, func(stmt *sqlite.Stmt) error {
        count = stmt.ColumnInt(0)
        return nil
    })
    if count != 1 {
        t.Errorf("Expected rollback to leave 1 row, got %d", count)
    }

    _ = db.CloseConn()
}
