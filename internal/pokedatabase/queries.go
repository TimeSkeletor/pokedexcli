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
	asset_path = "assets/pokemon-sprites/sprites/pokemon/versions/generation-viii/icons/"
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

func (db *Database) FetchPokemon(ctx context.Context, table string, number int) int {
    if table != "pokemon" {
        panic(fmt.Sprintf("invalid table name: %s", table))
    }
    query := "SELECT name, caught FROM pokemon WHERE number = ? LIMIT 1"
    var name string
    var caught int
    err := db.ExecQuery(ctx, query, func(stmt *sqlite.Stmt) error {
        name = stmt.GetText("name")
        caught = int(stmt.GetInt64("caught"))
        log.Printf("ℹ️ [FetchPokemon] Pokémon %s found, caught: %d", name, caught)
        return nil
    }, number)
    if err != nil {
        log.Printf("ℹ️ [FetchPokemon] Error querying Pokémon #%d: %v", number, err)
        return 0
    }
    return caught
}

func (db *Database) CatchPokemon(ctx context.Context, table string, number int, isShiny bool) error {
    if table != "pokemon" {
        panic(fmt.Sprintf("invalid table name: %s", table))
    }
    log.Printf("📌 [CatchPokemon] Catching Pokemon #%d. [Shiny: %v]", number, isShiny)

    // Update in a transaction for atomicity
    caughtAt := time.Now().Format("2006-01-02 15:04:05") // SQLite-compatible format
    query := "UPDATE pokemon SET caught = ?, is_shiny = ?, caught_at = ? WHERE number = ?"
    return db.WithTransaction(ctx, func(conn *sqlite.Conn) error {
        return sqlitex.Exec(conn, query, nil, 1, isShiny, caughtAt, number)
    })
}

func (db *Database) RegisterPokemon(ctx context.Context, pkmn pokeapi.PokemonSpecies) {
    log.Printf("📥 [RegisterPokemon] Inserting Pokemon %s (#%d) into DB", pkmn.Name, pkmn.ID)

    query := `INSERT INTO pokemon 
    (number, name, gender, capture_rate, base_happiness, is_baby, is_legendary, is_mythical, is_shiny, sprite_path, caught, caught_at) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

    gender := getGender(pkmn.GenderRate)
    pkmnSpritePath := asset_path + fmt.Sprintf("%d%s", pkmn.ID, img_type)

    err := db.ExecQuery(ctx, query, nil, pkmn.ID, pkmn.Name, gender, pkmn.CaptureRate, pkmn.BaseHappiness, pkmn.IsBaby, pkmn.IsLegendary, pkmn.IsMythical, 0, pkmnSpritePath, 0, nil)
    if err != nil {
        return
    }
    log.Printf("✅ [RegisterPokemon] Successfully registered Pokemon %s (#%d)", pkmn.Name, pkmn.ID)
}