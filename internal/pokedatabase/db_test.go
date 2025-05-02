package pokedatabase

import (
	"log"
	"testing"

	"github.com/timeskeletor/pokedexcli/internal/pokeapi"
)

func TestIsRegistered(t *testing.T) {
    db := New()
    defer db.Close()

    log.Println("🧪 [TestIsRegistered] Starting test")

    // Test with a non-existent Pokémon
    log.Println("🧪 [TestIsRegistered] Checking non-existent Pokémon")
    found, wasCaught, err := db.isRegistered("999")
    if err != nil {
        t.Fatalf("isRegistered failed: %v", err)
    }
    if found || wasCaught {
        t.Errorf("Expected found=false, wasCaught=false, got found=%v, wasCaught=%v", found, wasCaught)
    }
    log.Println("🧪 [TestIsRegistered] Non-existent Pokémon check passed")

    // Insert a test Pokémon
    pkmn := pokeapi.PokemonSpecies{
        ID:             999,
        Name:           "testmon",
        CaptureRate:    45,
        BaseHappiness:  70,
        GenderRate:     4,
        IsBaby:         false,
        IsLegendary:    false,
        IsMythical:     false,
    }
    log.Println("🧪 [TestIsRegistered] Registering test Pokémon")
    if err := db.RegisterPokemon(pkmn, true, false); err != nil {
        t.Fatalf("RegisterPokemon failed: %v", err)
    }
    log.Println("🧪 [TestIsRegistered] Test Pokémon registered")

    // Test with an existing Pokémon
    log.Println("🧪 [TestIsRegistered] Checking existing Pokémon")
    found, wasCaught, err = db.isRegistered("999")
    if err != nil {
        t.Fatalf("isRegistered failed: %v", err)
    }
    if !found || !wasCaught {
        t.Errorf("Expected found=true, wasCaught=true, got found=%v, wasCaught=%v", found, wasCaught)
    }
    log.Println("🧪 [TestIsRegistered] Existing Pokémon check passed")
}