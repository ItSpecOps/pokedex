package pokedex

type Pokemon struct {
    Name   string
    Type   string
    Ability string
}

type Pokedex struct {
    entries map[string]Pokemon
}

func NewPokedex() *Pokedex {
    return &Pokedex{
        entries: make(map[string]Pokemon),
    }
}

func (p *Pokedex) AddPokemon(pokemon Pokemon) {
    p.entries[pokemon.Name] = pokemon
}

func (p *Pokedex) GetPokemon(name string) (Pokemon, bool) {
    pokemon, exists := p.entries[name]
    return pokemon, exists
}

func (p *Pokedex) ListPokemons() []Pokemon {
    pokemons := []Pokemon{}
    for _, pokemon := range p.entries {
        pokemons = append(pokemons, pokemon)
    }
    return pokemons
}