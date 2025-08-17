package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io"
    "math/rand"
    "net/http"
    "os"
    "strings"
    "pokedex/internal/cli"
    "pokedex/internal/pokecache"
    "time"
)

// Struct to store caught Pokémon details
type CaughtPokemon struct {
    Name   string
    Height int
    Weight int
    Stats  []struct {
        Name  string
        Value int
    }
    Types []string
}

var (
    nextLocationURL    = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
    prevLocationURL    = ""
    currentLocationURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
    cache              = pokecache.NewCache(5 * time.Second)
    caughtPokemon       = make(map[string]CaughtPokemon)
)

func commandExit(args []string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil // not reached
}

func commandHelp(args []string) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:")
    fmt.Println()
    fmt.Println("help: Displays a help message")
    fmt.Println("exit: Exit the Pokedex")
    fmt.Println("map: Explore the Pokemon world map")
    fmt.Println("mapb: Go back to the previous page of the Pokemon world map")
    fmt.Println("explore <location-area>: List all Pokémon in a location area")
    fmt.Println("catch <pokemon>: Attempt to catch a Pokémon by name")
    return nil
}

func fetchAndPrintLocations(url string) (next string, prev string, err error) {
    var body []byte
    if cached, ok := cache.Get(url); ok {
        fmt.Println("(cache hit)")
        body = cached
    } else {
        resp, err := http.Get(url)
        if err != nil {
            return "", "", fmt.Errorf("failed to fetch locations: %v", err)
        }
        defer resp.Body.Close()
        body, err = io.ReadAll(resp.Body)
        if err != nil {
            return "", "", fmt.Errorf("failed to read response: %v", err)
        }
        cache.Add(url, body)
    }

    var result struct {
        Results []struct {
            Name string `json:"name"`
        } `json:"results"`
        Next     string `json:"next"`
        Previous string `json:"previous"`
    }
    if err := json.Unmarshal(body, &result); err != nil {
        return "", "", fmt.Errorf("failed to parse response: %v", err)
    }

    for _, loc := range result.Results {
        fmt.Println(loc.Name)
    }
    return result.Next, result.Previous, nil
}

func commandMap(args []string) error {
    next, _, err := fetchAndPrintLocations(nextLocationURL)
    if err != nil {
        return err
    }
    prevLocationURL = currentLocationURL
    currentLocationURL = nextLocationURL
    nextLocationURL = next
    return nil
}

func commandMapBack(args []string) error {
    if prevLocationURL == "" {
        fmt.Println("you're on the first page")
        return nil
    }
    _, prev, err := fetchAndPrintLocations(prevLocationURL)
    if err != nil {
        return err
    }
    nextLocationURL = currentLocationURL
    currentLocationURL = prevLocationURL
    prevLocationURL = prev
    return nil
}

func commandExplore(args []string) error {
    if len(args) < 2 {
        fmt.Println("Usage: explore <location-area>")
        return nil
    }
    area := args[1]
    url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", area)

    var body []byte
    if cached, ok := cache.Get(url); ok {
        fmt.Println("(cache hit)")
        body = cached
    } else {
        resp, err := http.Get(url)
        if err != nil {
            return fmt.Errorf("failed to fetch location area: %v", err)
        }
        defer resp.Body.Close()
        body, err = io.ReadAll(resp.Body)
        if err != nil {
            return fmt.Errorf("failed to read response: %v", err)
        }
        cache.Add(url, body)
    }

    var result struct {
        PokemonEncounters []struct {
            Pokemon struct {
                Name string `json:"name"`
            } `json:"pokemon"`
        } `json:"pokemon_encounters"`
    }
    if err := json.Unmarshal(body, &result); err != nil {
        return fmt.Errorf("failed to parse response: %v", err)
    }

    fmt.Printf("Exploring %s...\n", area)
    if len(result.PokemonEncounters) == 0 {
        fmt.Println("No Pokémon found in this location area.")
        return nil
    }

    fmt.Println("Found Pokemon:")
    for _, encounter := range result.PokemonEncounters {
        fmt.Printf(" - %s\n", encounter.Pokemon.Name)
    }
    return nil
}

func commandPokedex(args []string) error {
    if len(caughtPokemon) == 0 {
        fmt.Println("You haven't caught any Pokémon yet.")
        return nil
    }
    fmt.Println("Your Pokedex:")
    for name := range caughtPokemon {
        fmt.Printf(" - %s\n", name)
    }
    return nil
}

func commandCatch(args []string) error {
    if len(args) < 2 {
        fmt.Println("Usage: catch <pokemon>")
        return nil
    }
    name := args[1]
    if _, ok := caughtPokemon[name]; ok {
        fmt.Printf("You already caught %s!\n", name)
        return nil
    }
    fmt.Printf("Throwing a Pokeball at %s...\n", name)
    url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", name)

    var body []byte
    if cached, ok := cache.Get(url); ok {
        fmt.Println("(cache hit)")
        body = cached
    } else {
        resp, err := http.Get(url)
        if err != nil {
            fmt.Printf("%s could not be found!\n", name)
            return nil
        }
        defer resp.Body.Close()
        body, err = io.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("%s could not be found!\n", name)
            return nil
        }
        cache.Add(url, body)
    }

    var pokeData struct {
        Name           string `json:"name"`
        BaseExperience int    `json:"base_experience"`
        Height         int    `json:"height"`
        Weight         int    `json:"weight"`
        Stats          []struct {
            Stat struct {
                Name string `json:"name"`
            } `json:"stat"`
            BaseStat int `json:"base_stat"`
        } `json:"stats"`
        Types []struct {
            Type struct {
                Name string `json:"name"`
            } `json:"type"`
        } `json:"types"`
    }
    if err := json.Unmarshal(body, &pokeData); err != nil {
        fmt.Printf("%s could not be found!\n", name)
        return nil
    }

    // Lower base_experience = easier to catch. We'll use a simple formula:
    // chance = max(10, 100 - base_experience)
    chance := 100 - pokeData.BaseExperience
    if chance < 10 {
        chance = 10
    }
    rand.Seed(time.Now().UnixNano())
    roll := rand.Intn(100)
    if roll < chance {
        fmt.Printf("%s was caught!\n", pokeData.Name)
        // Store all details for inspect
        stats := make([]struct {
            Name  string
            Value int
        }, len(pokeData.Stats))
        for i, s := range pokeData.Stats {
            stats[i].Name = s.Stat.Name
            stats[i].Value = s.BaseStat
        }
        types := make([]string, len(pokeData.Types))
        for i, t := range pokeData.Types {
            types[i] = t.Type.Name
        }
        caughtPokemon[pokeData.Name] = CaughtPokemon{
            Name:   pokeData.Name,
            Height: pokeData.Height,
            Weight: pokeData.Weight,
            Stats:  stats,
            Types:  types,
        }
        fmt.Println("You may now inspect it with the inspect command.")
    } else {
        fmt.Printf("%s escaped!\n", pokeData.Name)
    }
    return nil
}

func commandInspect(args []string) error {
    if len(args) < 2 {
        fmt.Println("Usage: inspect <pokemon>")
        return nil
    }
    name := args[1]
    p, ok := caughtPokemon[name]
    if !ok {
        fmt.Println("you have not caught that pokemon")
        return nil
    }
    fmt.Printf("Name: %s\n", p.Name)
    fmt.Printf("Height: %d\n", p.Height)
    fmt.Printf("Weight: %d\n", p.Weight)
    fmt.Println("Stats:")
    for _, stat := range p.Stats {
        fmt.Printf("  -%s: %d\n", stat.Name, stat.Value)
    }
    fmt.Println("Types:")
    for _, t := range p.Types {
        fmt.Printf("  - %s\n", t)
    }
    return nil
}

var commands = map[string]cli.Command{
    "exit": {
        Name:        "exit",
        Description: "Exit the Pokedex",
        Callback:    commandExit,
    },
    "help": {
        Name:        "help",
        Description: "Displays a help message",
        Callback:    commandHelp,
    },
    "map": {
        Name:        "map",
        Description: "Explore the Pokemon world map",
        Callback:    commandMap,
    },
    "mapb": {
        Name:        "mapb",
        Description: "Go back to the previous page of the Pokemon world map",
        Callback:    commandMapBack,
    },
    "explore": {
        Name:        "explore",
        Description: "List all Pokémon in a location area. Usage: explore <location-area>",
        Callback:    commandExplore,
    },
    "catch": {
        Name:        "catch",
        Description: "Catch a Pokemon by name. Usage: catch <pokemon>",
        Callback:    commandCatch,
    },
    "inspect": {
        Name:        "inspect",
        Description: "Show details about a caught Pokemon. Usage: inspect <pokemon>",
        Callback:    commandInspect,
    },
    "pokedex": {
        Name:        "pokedex",
        Description: "List all Pokémon you have caught",
        Callback:    commandPokedex,
    },
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        if !scanner.Scan() {
            break // EOF or error
        }
        input := scanner.Text()
        words := cleanInput(input)
        if len(words) == 0 {
            continue
        }
        cmdName := strings.ToLower(words[0])
        cmd, ok := commands[cmdName]
        if !ok {
            fmt.Println("Unknown command")
            continue
        }
        if err := cmd.Callback(words); err != nil {
            fmt.Println(err)
        }
    }
}

// cleanInput lowercases, trims, and splits the input string.
func cleanInput(text string) []string {
    fields := strings.Fields(strings.ToLower(strings.TrimSpace(text)))
    return fields
}

