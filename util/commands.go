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
	callback    func(map[string]cliCommand) //TODO: Refactor this to a more appropriate function signature
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
			callback:    func(map[string]cliCommand) { CommandExit() }, //TODO: Make this hacky stuff not needed
		},
		"map": {
			name:        "map",
			description: "display the next locations",
			callback: func(m map[string]cliCommand) { //implement
				LocationAreaState = CommandMap(LocationAreaState)
			},
		},
		"mapb": {
			name:        "map back",
			description: "display the previous locations",
			callback: func(m map[string]cliCommand) {
				LocationAreaState = CommandMapB(LocationAreaState)
			},
		},
	}

	return commands
}

func ParseCommand(in string, commands map[string]cliCommand) {
	command, ok := commands[SanitizeInput(in)] //clean and check if the input exists
	if !ok {
		HandleUnknownKeys(in)
	} else {
		command.callback(commands)
	}
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

func CommandHelp(commands map[string]cliCommand) {
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
	fmt.Printf("%s is not a recognised command, please try again", in)
}

func CommandMap(areas LocationAreaBatch) LocationAreaBatch {
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

func CommandMapB(areas LocationAreaBatch) LocationAreaBatch { //TODO: Test and make sure refactoring is not needed here
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
