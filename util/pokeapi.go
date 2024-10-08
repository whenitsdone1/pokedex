package util

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const (
	LocationAreaApiUrl = "https://pokeapi.co/api/v2/location-area"
)

type LocationAreaBatch struct {
	Next     string         `json:"next"`
	Results  []LocationArea `json:"results"`
	Previous string         `json:"previous"`
}
type Parseable interface {
	Parse(data []byte) error
}
type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type PokemonDetails struct {
	PokemonEncounters []struct {
		Pokemon Pokemon `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
type PokeDexInformation struct {
	Exp    int `json:"base_experience"`
	Weight int `json:"weight"`
	Height int `json:"height"`
	Forms  []struct {
		Name string `json:"name"`
	} `json:"forms"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func (p *PokeDexInformation) GetStatsMap() map[string]int {
	stats := make(map[string]int)
	for _, stat := range p.Stats {
		stats[stat.Stat.Name] = stat.BaseStat
	}
	return stats
}

func (p *PokeDexInformation) Parse(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *PokeDexInformation) GetTypes() []string {
	var types []string
	for _, t := range p.Types {
		types = append(types, t.Type.Name)
	}
	return types
}

func ParsePokedexDetails(url string, cache *Cache) (PokeDexInformation, error) {
	result, err := Parse[*PokeDexInformation](url, cache)
	if err != nil {
		return PokeDexInformation{}, err
	}
	if result == nil {
		return PokeDexInformation{}, fmt.Errorf("error parsing result %v", err)
	}
	return *result, nil
}

func Parse[T Parseable](url string, cache *Cache) (T, error) {
	var item T
	if val, ok := cache.Get(url); ok {
		fmt.Println("found in cache, retrieving...")
		parsed, ok := val.(T)
		if ok {
			return parsed, nil
		}
	}
	fmt.Println("not found in cache, hitting api..")
	jsonData, err := GetJson(url)
	if err != nil {
		return item, fmt.Errorf("error fetching data: %w", err)
	}

	switch any(item).(type) { //convert the nil value T to a concrete instance of the appropiate struct, allows our .parse methods to function correctly
	case *LocationAreaBatch:
		item = any(&LocationAreaBatch{}).(T)
	case *PokemonDetails:
		item = any(&PokemonDetails{}).(T)
	case *PokeDexInformation:
		item = any(&PokeDexInformation{}).(T)
	default:
		return item, fmt.Errorf("unsupported type for parsing")
	}
	err = item.Parse(jsonData)
	if err != nil {
		return item, fmt.Errorf("error parsing data: %w", err)
	}
	cache.Add(url, item)
	return item, nil
}

func (l *LocationAreaBatch) Parse(data []byte) error {
	err := json.Unmarshal(data, l)
	if err != nil {
		return nil
	}
	slices.SortFunc(l.Results, CompareLocations)
	return nil
}

func (p *PokemonDetails) Parse(data []byte) error {
	return json.Unmarshal(data, p)
}

func ParseLocationAreas(toParse string, cache *Cache) (LocationAreaBatch, error) {
	result, err := Parse[*LocationAreaBatch](toParse, cache)
	if err != nil {
		return LocationAreaBatch{}, err
	}
	if result == nil {
		return LocationAreaBatch{}, fmt.Errorf("error parsing result %v", err)
	}
	return *result, nil
}

func ParseLocations(toParse string, cache *Cache) ([]Pokemon, error) {
	result, err := Parse[*PokemonDetails](toParse, cache)
	if err != nil {
		return nil, err
	}
	var pokemon []Pokemon
	details := *result
	for _, p := range details.PokemonEncounters {
		pokemon = append(pokemon, p.Pokemon)
	}
	return pokemon, nil
}

func GetJson(toParse string) ([]byte, error) {
	response, err := http.Get(toParse)
	if err != nil {
		return nil, errors.New("error making get request")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading body")
	}
	return body, nil
}

func CompareLocations(x, y LocationArea) int {
	a, errA := strconv.Atoi(getNumber(x.Url))
	if errA != nil {
		log.Fatal("error in parsing of URLs")
	}
	b, errB := strconv.Atoi(getNumber(y.Url))
	if errB != nil {
		log.Fatal("error in parsing of URLs")
	}
	return a - b
}
func getNumber(x string) string { //retrieve the number of the location area from the url
	x = strings.TrimPrefix(x, "https://pokeapi.co/api/v2/location-area/")
	x = strings.TrimSuffix(x, "/")
	return x
}

func ExtractNames(p []Pokemon) []string {
	toCompare := []string{}
	for _, v := range p {
		toCompare = append(toCompare, strings.TrimSpace(v.Name))
	}
	slices.SortFunc(toCompare, func(a, b string) int {
		return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
	})
	return toCompare
}

func DisplayPokemonInformation(p PokeDexInformation) {
	fmt.Printf("Name: %v\n", p.Forms[0].Name)
	fmt.Printf("Height: %v\n", p.Height)
	fmt.Printf("Weight: %v\n", p.Weight)
	fmt.Printf("Stats:\n")
	stats := p.GetStatsMap()
	fmt.Printf("-hp: %v\n", stats["hp"])
	fmt.Printf("-attack: %v\n", stats["attack"])
	fmt.Printf("-defense: %v\n", stats["defense"])
	fmt.Printf("-special-attack: %v\n", stats["special-attack"])
	fmt.Printf("-special-defense: %v\n", stats["special-defense"])
	fmt.Printf("-speed: %v\n", stats["speed"])
	fmt.Printf("-Types:\n")
	types := p.GetTypes()
	fmt.Printf("- %v\n", types[0])
	if len(types) == 2 {
		fmt.Printf("- %v\n", types[1])
	}
}

func (p *PokeDexInformation) CatchChance() bool {
	clamp := func(x int) int {
		max := 60 //set max to 60 rather than natural 99 to make catching pokemon easier
		min := 1
		if x > max {
			return max
		}
		if x < min {
			return min
		}
		return x
	}
	chance := clamp(p.Exp)
	probability := 1.0 - float64(chance)/100.0
	rand := rand.Float64()
	return rand < probability
}
