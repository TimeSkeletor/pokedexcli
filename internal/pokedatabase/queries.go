package pokedatabase

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
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

func (db *Database) isRegistered(identifier string) (found bool, wasCaught bool, err error) {
    log.Printf("🔍 [isRegistered] Start: identifier=%s", identifier)

    conn := db.pool.Get(nil)
    if conn == nil {
        log.Printf("❌ [isRegistered] Failed to get connection from pool")
        return false, false, fmt.Errorf("failed to get connection from pool")
    }
    defer db.pool.Put(conn)

	stmt := `
	SELECT id, name, caught
	FROM pokemon 
	WHERE number = ? OR name = ?
	LIMIT 1
	`

    var id int64
    var name string
    var caughtInt int64
    callbackCount := 0

    err = sqlitex.Exec(conn, stmt, func(stmt *sqlite.Stmt) error {
        callbackCount++
        log.Printf("🔍 [isRegistered] Callback invoked %d time(s)", callbackCount)
        if ok, err := stmt.Step(); err != nil {
            log.Printf("❌ [isRegistered] Error stepping through result: %v", err)
            return err
        } else if !ok {
            log.Printf("ℹ️ [isRegistered] No Pokémon found in DB")
            return nil
        }

        id = stmt.ColumnInt64(0)
        name = stmt.ColumnText(1)
        caughtInt = stmt.ColumnInt64(2)
        found = true
        wasCaught = caughtInt != 0
        log.Printf("✅ [isRegistered] Found Pokémon: ID=%d, Name=%s, WasCaught=%v", id, name, wasCaught)
        return nil
    }, identifier, identifier)

    if err != nil {
        log.Printf("❌ [isRegistered] Error during lookup: %v", err)
        return false, false, fmt.Errorf("failed to find pokemon: %v", err)
    }

    log.Printf("🔍 [isRegistered] Completed: callbackCount=%d, found=%v, wasCaught=%v", callbackCount, found, wasCaught)
    return found, wasCaught, nil
}

func updateCaughtStatus(conn *sqlite.Conn, number int, caughtAt time.Time) error {
	log.Printf("📌 Updating caught status for Pokémon #%d at %s", number, caughtAt.Format(time.RFC3339))

	stmt := `
	UPDATE pokemon
	SET caught = 1, caught_at = ?
	WHERE number = ?
	`
	err := sqlitex.Exec(conn, stmt, nil, caughtAt, number)
	if err != nil {
		log.Printf("❌ Failed to update caught status: %v", err)
	}
	return err
}

func catchQuery(conn *sqlite.Conn, pkmn pokeapi.PokemonSpecies, caught bool, isShiny bool) error {
    log.Printf("📥 [catchQuery] Inserting Pokémon %s (#%d) into DB [Caught: %v, Shiny: %v]", pkmn.Name, pkmn.ID, caught, isShiny)

    stmt := `INSERT INTO pokemon 
    (number, name, gender, capture_rate, base_happiness, is_baby, is_legendary, is_mythical, is_shiny, sprite_path, caught, caught_at) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

    gender := getGender(pkmn.GenderRate)
    pkmnSpritePath := asset_path + strconv.Itoa(pkmn.ID) + img_type

    var caughtAt interface{}
    if caught {
        caughtAt = time.Now()
    } else {
        caughtAt = nil
    }

    err := sqlitex.ExecTransient(
        conn,
        stmt,
        nil, // No callback needed for INSERT
        pkmn.ID,
        pkmn.Name,
        gender,
        pkmn.CaptureRate,
        pkmn.BaseHappiness,
        pkmn.IsBaby,
        pkmn.IsLegendary,
        pkmn.IsMythical,
        isShiny,
        pkmnSpritePath,
        caught,
        caughtAt,
    )

    if err != nil {
        log.Printf("❌ [catchQuery] Failed to insert Pokémon %s (#%d): %v", pkmn.Name, pkmn.ID, err)
    } else {
        log.Printf("✅ [catchQuery] Successfully inserted Pokémon %s (#%d)", pkmn.Name, pkmn.ID)
    }

    return err
}

func (db *Database) RegisterPokemon(pkmn pokeapi.PokemonSpecies, caught bool, isShiny bool) error {
    log.Printf("🚀 [RegisterPokemon] Registering Pokémon: %s (#%d) | Caught: %v | Shiny: %v", pkmn.Name, pkmn.ID, caught, isShiny)

    return withConnTxFn(db.pool, func(conn *sqlite.Conn, _ struct{}) error {
        found, wasCaught, err := db.isRegistered(strconv.Itoa(pkmn.ID))
        if err != nil {
            log.Printf("❌ [RegisterPokemon] isRegistered failed: %v", err)
            return err
        }

        if !found {
            log.Printf("📄 [RegisterPokemon] Pokémon not found. Proceeding with insertion...")
            return catchQuery(conn, pkmn, caught, isShiny)
        }

        log.Printf("📄 [RegisterPokemon] Pokémon already registered. Checking caught status...")
        if caught && !wasCaught {
            log.Printf("📌 [RegisterPokemon] Updating caught status...")
            return updateCaughtStatus(conn, pkmn.ID, time.Now())
        }

        log.Printf("ℹ️ [RegisterPokemon] No update needed for %s (#%d) — Already caught: %v", pkmn.Name, pkmn.ID, wasCaught)
        return nil
    }, struct{}{})
}