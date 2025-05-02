package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/timeskeletor/pokedexcli/config"
    "github.com/timeskeletor/pokedexcli/internal/pokeapi"
    "github.com/timeskeletor/pokedexcli/internal/pokedatabase"
)

func main() {
    pokeClient := pokeapi.NewClient(5*time.Second, time.Minute*5)
    cfg := &config.Config{
        PokeapiClient: pokeClient,
        PokeDb:        &pokedatabase.Database{},
        CaughtPkmn:    make(map[string]pokeapi.Pokemon),
    }

    // Initialize database
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := cfg.PokeDb.SetConn(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize database connection: %v\n", err)
        os.Exit(1)
    }
    if err := cfg.PokeDb.SetDb(ctx); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to set up database: %v\n", err)
        os.Exit(1)
    }
    defer cfg.PokeDb.CloseConn()

    startRepl(cfg)
}