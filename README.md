# Pokedex

Pokedex is a simple application that allows users to manage a collection of Pokémon entries. This project serves as a command-line interface for adding, retrieving, and managing Pokémon data.

## Project Structure

```
pokedex
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   └── pokedex
│       └── pokedex.go   # Core functionality for managing Pokémon entries
├── go.mod               # Module definition and dependencies
└── README.md            # Documentation for the project
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

## Usage Examples

- To add a Pokémon entry, use the appropriate command in the CLI.
- To retrieve a Pokémon entry, use the corresponding command.

For more detailed usage instructions, refer to the documentation in the `cmd/main.go` file.