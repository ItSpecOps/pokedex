# Pokedex

Pokedex is a command-line application that allows users to explore the Pokémon world, catch Pokémon, and manage their own Pokédex. It interacts with the [PokéAPI](https://pokeapi.co/) and uses in-memory caching for fast repeated queries.

## Project Structure

```
pokedex
├── cmd
│   └── main.go              # Entry point and REPL logic
├── internal
│   ├── cli
│   │   └── cli.go           # CLI command abstraction
│   ├── pokecache
│   │   ├── pokecache.go     # In-memory cache for API responses
│   │   └── pokecache_test.go
│   └── pokedex
│       └── pokedex.go       # (Legacy) Example pokedex logic
├── go.mod                   # Module definition and dependencies
└── README.md                # Documentation for the project
```

## Setup Instructions

1. Clone the repository:
   ```
   git clone <repository-url>
   cd pokedex
   ```

2. Initialize the Go module:
   ```
   go mod tidy
   ```

3. Run the application:
   ```
   go run cmd/main.go
   ```

## Usage

The CLI supports the following commands:

- `help`  
  Displays a help message with available commands.

- `exit`  
  Exit the Pokedex.

- `map`  
  Explore the Pokémon world map, showing 20 location areas at a time.

- `mapb`  
  Go back to the previous page of the Pokémon world map.

- `explore <location-area>`  
  List all Pokémon in a given location area.

- `catch <pokemon>`  
  Attempt to catch a Pokémon by name. The chance of catching depends on the Pokémon's base experience.

- `inspect <pokemon>`  
  Show details about a caught Pokémon (name, height, weight, stats, types). Only works for Pokémon you have caught.

- `pokedex`  
  List all Pokémon you have caught.

## Example Session

```
Pokedex > map
canalave-city-area
eterna-city-area
...
Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 - magikarp
...
Pokedex > catch pidgey
Throwing a Pokeball at pidgey...
pidgey was caught!
You may now inspect it with the inspect command.
Pokedex > inspect pidgey
Name: pidgey
Height: 3
Weight: 18
Stats:
  -hp: 40
  -attack: 45
  -defense: 40
  -special-attack: 35
  -special-defense: 35
  -speed: 56
Types:
  - normal
  - flying
Pokedex > pokedex
Your Pokedex:
 - pidgey
```

## Notes

- The application uses an in-memory cache for API responses, making repeated queries fast.
- Only Pokémon you have caught can be inspected.
- Use `help` at any time to see available commands.