package pokedatabase

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"github.com/timeskeletor/pokedexcli/internal/pokeapi"
)

const (
	asset_path = "assets/pokemon-sprites/sprites/pokemon/versions/generation-viii/"
	img_type   = ".png"
)

func getGender(genderRate int) int {
	if genderRate == -1 {
		return -1
	}
	r := rand.Intn(8)
	if r < genderRate {
		return 1
	}
	return 0
}

func (db *Database) FetchPokemon(ctx context.Context, table string, number int) (int, string, error) {
    if table != "pokemon" {
        return 0, "", fmt.Errorf("invalid table name: %s", table)
    }
    query := "SELECT number, name FROM pokemon WHERE number = ? LIMIT 1"
    var foundNumber int
    var foundName string
    err := db.ExecQuery(ctx, query, func(stmt *sqlite.Stmt) error {
        foundNumber = int(stmt.GetInt64("number"))
        foundName = stmt.GetText("name")
        return nil
    }, number)
    if err != nil {
        log.Printf("ℹ️ [FetchPokemon] No Pokemon found for number %d: %v", number, err)
        return 0, "", fmt.Errorf("failed to fetch Pokemon: %w", err)
    }
    if foundNumber == 0 {
        log.Printf("ℹ️ [FetchPokemon] No Pokemon found for number %d", number)
        return 0, "", fmt.Errorf("pokemon not found")
    }
    return foundNumber, foundName, nil
}

func (db *Database) CatchPokemon(ctx context.Context, table string, number int) error {
    if table != "pokemon" {
        return fmt.Errorf("invalid table name: %s", table)
    }
    log.Printf("📌 Updating caught status for Pokemon #%d", number)

    // Check if Pokemon exists
    _, _, err := db.FetchPokemon(ctx, table, number)
    if err != nil {
        log.Printf("ℹ️ [CatchPokemon] No Pokemon found for number %d: %v", number, err)
        return fmt.Errorf("pokemon not found: %w", err)
    }

    // Update in a transaction for atomicity
    caughtAt := time.Now().Format("2006-01-02 15:04:05") // SQLite-compatible format
    query := "UPDATE pokemon SET caught = ?, caught_at = ? WHERE number = ?"
    return db.WithTransaction(ctx, func(conn *sqlite.Conn) error {
        return sqlitex.Exec(conn, query, nil, 1, caughtAt, number)
    })
}

func (db *Database) RegisterPokemon(ctx context.Context, pkmn pokeapi.PokemonSpecies, caught bool, isShiny bool) error {
    log.Printf("📥 [RegisterPokemon] Inserting Pokemon %s (#%d) into DB [Caught: %v, Shiny: %v]", pkmn.Name, pkmn.ID, caught, isShiny)

    query := `INSERT INTO pokemon 
    (number, name, gender, capture_rate, base_happiness, is_baby, is_legendary, is_mythical, is_shiny, sprite_path, caught, caught_at) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

    gender := getGender(pkmn.GenderRate)
    pkmnSpritePath := asset_path + fmt.Sprintf("%d%s", pkmn.ID, img_type)
    var caughtAt interface{}
    if caught {
        caughtAt = time.Now().Format("2006-01-02 15:04:05")
    } else {
        caughtAt = nil // NULL for uncaught Pokemon
    }

    err := db.ExecQuery(ctx, query, nil, pkmn.ID, pkmn.Name, gender, pkmn.CaptureRate, pkmn.BaseHappiness, pkmn.IsBaby, pkmn.IsLegendary, pkmn.IsMythical, isShiny, pkmnSpritePath, caught, caughtAt)
    if err != nil {
        log.Printf("❌ [RegisterPokemon] Failed to register Pokemon %s (#%d): %v", pkmn.Name, pkmn.ID, err)
        return fmt.Errorf("failed to register Pokemon: %w", err)
    }
    log.Printf("✅ [RegisterPokemon] Successfully registered Pokemon %s (#%d)", pkmn.Name, pkmn.ID)
    return nil
}