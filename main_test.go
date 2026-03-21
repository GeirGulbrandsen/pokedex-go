package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/geirgulbrandsen/pokedex-go/internal/pokecache"
)

func TestCommandExploreRequiresAreaName(t *testing.T) {
	cfg := &config{Cache: pokecache.NewCache(1 * time.Minute)}

	err := commandExplore(cfg, []string{})
	if err == nil {
		t.Fatal("expected error when no area name is provided")
	}
}

func TestCommandExploreReadsFromCacheAndPrintsPokemon(t *testing.T) {
	cfg := &config{Cache: pokecache.NewCache(1 * time.Minute)}
	areaName := "test-area"
	url := "https://pokeapi.co/api/v2/location-area/test-area"
	body := []byte(`{"pokemon_encounters":[{"pokemon":{"name":"pikachu"}},{"pokemon":{"name":"eevee"}}]}`)
	cfg.Cache.Add(url, body)

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	os.Stdout = w

	cmdErr := commandExplore(cfg, []string{areaName})

	_ = w.Close()
	os.Stdout = oldStdout

	if cmdErr != nil {
		t.Fatalf("expected no error, got: %v", cmdErr)
	}

	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		t.Fatalf("failed reading command output: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Exploring test-area...") {
		t.Fatalf("expected explore header in output, got: %q", output)
	}
	if !strings.Contains(output, "Found Pokemon:") {
		t.Fatalf("expected pokemon header in output, got: %q", output)
	}
	if !strings.Contains(output, " - pikachu") || !strings.Contains(output, " - eevee") {
		t.Fatalf("expected pokemon names in output, got: %q", output)
	}
}
