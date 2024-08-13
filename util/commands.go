package util

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

var UnknownKeyError error
var LocationAreaState LocationAreaBatch
var DataStore = NewCache(10 * time.Second) //Set cache duration here
var AvailablePokemon []string

type cliCommand struct {
	name        string
	description string
	callback    func(args []string, cmd map[string]cliCommand)
}

func NewCommandMap() map[string]cliCommand {
	//closures match the function signatures with the signature outlined in callback
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "display a help message",
			callback:    CommandHelp,
		},
		"exit": {
			name:        "exit",
			description: "exit the pokedex",
			callback:    func(args []string, cmd map[string]cliCommand) { CommandExit() },
		},
		"map": {
			name:        "map",
			description: "display the next locations",
			callback: func(args []string, m map[string]cliCommand) {
				LocationAreaState = CommandMap(args, LocationAreaState)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "display the previous locations",
			callback: func(args []string, cmd map[string]cliCommand) {
				LocationAreaState = CommandMapB(args, LocationAreaState)
			},
		},
		"explore": {
			name:        "explore",
			description: "look for pokemon!",
			callback: func(args []string, cmd map[string]cliCommand) {
				AvailablePokemon = CommandExplore(args, LocationAreaState)
			},
		},
		"catch": {
			name:        "catch",
			description: "try to catch the pokemeon!",
			callback: func(args []string, cmd map[string]cliCommand) {
				CommandCatch(args, LocationAreaState) //implement
			},
		},
	}
	return commands
}

func ParseCommand(in []string, commands map[string]cliCommand) []string { //handle commans
	if len(in) == 0 {
		fmt.Println("enter a command to start!")
		return nil
	}
	command, ok := commands[SanitizeInput(in[0])] //clean and check if the input exists
	if !ok {
		HandleUnknownKeys(in[0])
	} else {
		command.callback(in, commands)
	}
	return in
}

func SanitizeInput(in string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`) //forbidden characters
	cleansed_in := func(in string) string {
		return nonAlphanumericRegex.ReplaceAllString(in, "")
	}(in)
	lowered_in := strings.ToLower(cleansed_in)
	sanitized := strings.TrimSpace(lowered_in)
	return sanitized
}

func CommandHelp(args []string, commands map[string]cliCommand) {
	var loopVar int
	for _, entry := range commands {
		if entry.name == "help" {
			continue
		}
		if loopVar < len(commands)-1 {
			fmt.Printf("%v - %v\n", entry.name, entry.description)
		} else { //don't create an extra new line on the last iteration
			fmt.Printf("%v - %v", entry.name, entry.description)
		}
	}
}

func CommandExit() {
	fmt.Println("exiting, byebye!")
	os.Exit(0)
}

func HandleUnknownKeys(in string) {
	fmt.Printf("%s is not a recognised command, please try again\n", in)
}

func CommandMap(args []string, areas LocationAreaBatch) LocationAreaBatch {
	url := LocationAreaApiUrl //set the url to the first page
	if areas.Next != "" {     //if the areas argument to this function has a next (i.e: if we are not at the start)
		url = areas.Next //use the next url
	}
	next, err := ParseLocationAreas(url, DataStore)
	if err != nil {
		return LocationAreaBatch{}
	}

	for _, area := range next.Results {
		fmt.Println(area.Name)
	}

	if next.Next == "" {
		fmt.Println("You've reached the end of the world, we can't go any further!")
	}
	return next
}

func CommandMapB(args []string, areas LocationAreaBatch) LocationAreaBatch {
	if areas.Previous == "null" || areas.Previous == "" {
		fmt.Println("we're still at the start!")
		return areas
	}
	if areas.Previous != "null" && areas.Previous != "" {
		last, err := ParseLocationAreas(areas.Previous, DataStore)
		if err != nil {
			fmt.Printf("%v oopsie when parsing jason", err)
		}
		for i := range last.Results {
			fmt.Println(last.Results[i].Name)
		}
		return last
	}
	fmt.Println("we're back at the start!")
	return areas
}

func CommandExplore(args []string, area LocationAreaBatch) []string {
	if len(args) == 1 {
		fmt.Println("where should we explore?")
		return nil
	}
	in := args[1]
	arg := SanitizeInput(args[1]) //should be the name of a location
	if len(LocationAreaState.Results) == 0 {
		fmt.Println("we need to start exploring to catch pokemon!")
		return nil
	}

	for _, v := range LocationAreaState.Results {
		raw_name := v.Name
		sanitizedName := SanitizeInput(v.Name)
		if arg == sanitizedName {
			locations, _ := ParseLocations(v.Url, DataStore)
			ordered := ExtractNames(locations)
			if len(ordered) == 0 {
				fmt.Println("didn't find any pokemon :(")
				return nil
			}
			fmt.Printf("exploring %s:\n", raw_name)
			fmt.Println("found pokemon:")
			for i := range ordered {
				fmt.Println("-" + ordered[i])
			}
			return ordered
		}
	}
	fmt.Printf("invalid location: %s entered\n", in)
	return nil
}

func CommandCatch(args []string, area LocationAreaBatch) {
	pokemon := strings.TrimSpace(strings.ToLower(args[1]))
	for _, v := range AvailablePokemon {
		if v == pokemon {
			pokeInfo, err := ParsePokedexDetails(pokemon.Url, cache)
			fmt.Printf("throwing a pokeball at %v\n", pokemon)
			chance := rand.Intn(10) + 1
			time.Sleep(1 * time.Second)
			if chance >= 5 {
				fmt.Printf("%v was caught!\n", pokemon)
				return
			} else {
				fmt.Printf("%v got away!\n", pokemon)
				return
			}
		}
	}
	fmt.Println("that pokemon isn't around here...")
}

//TODO: implement random chance of catching using base experience, implement inspect command, implement pokedex command
