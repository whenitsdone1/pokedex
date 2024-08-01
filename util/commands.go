package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

var UnknownKeyError error
var LocationAreaState LocationAreaBatch
var DataStore = NewCache(30 * time.Second) //Set cache duration here

type cliCommand struct {
	name        string
	description string
	callback    func(args []string, cmd map[string]cliCommand) //TODO: Refactor this to a more appropriate function signature
}

func NewCommandMap() map[string]cliCommand { //Does this need to be a function ?
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "display a help message",
			callback:    CommandHelp,
		},
		"exit": {
			name:        "exit",
			description: "exit the pokedex",
			callback:    func(args []string, cmd map[string]cliCommand) { CommandExit() }, //TODO: Make this hacky stuff not needed
		},
		"map": {
			name:        "map",
			description: "display the next locations",
			callback: func(args []string, m map[string]cliCommand) { //implement
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
				CommandExplore(args, LocationAreaState)
			},
		},
	}

	return commands
}

func ParseCommand(in []string, commands map[string]cliCommand) []string { //?
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
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`) //we need to clean the input prior to parsing
	cleansed_in := func(in string) string {
		return nonAlphanumericRegex.ReplaceAllString(in, "")
	}(in)
	lowered_in := strings.ToLower(cleansed_in)
	sanitized := strings.TrimSpace(lowered_in)
	return sanitized
}

func CommandHelp(args []string, commands map[string]cliCommand) {
	var loopVar int //TODO: Find a less hacky way to do this
	for _, entry := range commands {
		if entry.name == "help" {
			continue
		}
		if loopVar < len(commands)-1 {
			fmt.Printf("%v - %v\n", entry.name, entry.description)
			loopVar++
		} else {
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
	url := LocationAreaApiUrl //? are we setting the target URL to the start everytime? TODO: Understand this better
	if areas.Next != "" {
		url = areas.Next
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

func CommandMapB(args []string, areas LocationAreaBatch) LocationAreaBatch { //TODO: Test and make sure refactoring is not needed here
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

func CommandExplore(args []string, area LocationAreaBatch) {
	if len(args) == 1 {
		fmt.Println("where should we explore?")
		return
	}
	arg := SanitizeInput(args[1]) //should be the name of a location
	if len(LocationAreaState.Results) == 0 {
		fmt.Println("we need to start exploring to catch pokemon!")
		return
	}

	for _, v := range LocationAreaState.Results {
		raw_name := v.Name
		sanitizedName := SanitizeInput(v.Name)
		if arg == sanitizedName {
			locations, _ := ParseLocations(v.Url, DataStore)
			ordered := ExtractNames(locations)
			if len(ordered) == 0 {
				fmt.Println("didn't find any pokemon :(")
				return
			}
			fmt.Printf("exploring %s:\n", raw_name)
			fmt.Println("found pokemon:")
			for i := range ordered {
				fmt.Println("-" + ordered[i])
			}
			return
		}
	}
	fmt.Println("hmm coudn't find that location")
}
