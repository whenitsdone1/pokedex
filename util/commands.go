package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var UnknownKeyError error

type Counter struct {
	Value int
}

func (c *Counter) Advance() { //global state == bad, how are you meant to do this?
	c.Value += 20
}
func (c *Counter) Recede() {
	c.Value -= 40 //because we call advance at the end of calling map the pointer ends up 40 places ahead of the previous values we are interested in
}

var Pointer = &Counter{}

var LocationAreaState, _ = ParseLocationAreas(LocationAreaApiUrl, DataStore)

// var InitMap, _ = initMap() //how do we check errors here?
var DataStore = NewCache(5)
var FirstCall = true

//hold type information and methods

type cliCommand struct {
	name        string
	description string
	callback    func(map[string]cliCommand) //this is gonna cause problems...
}

func NewCommandMap() map[string]cliCommand { //this can't fail - it would fail prior to compile time
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "display a help message",
			callback:    CommandHelp,
		},
		"exit": {
			name:        "exit",
			description: "exit the pokedex",
			callback:    func(map[string]cliCommand) { CommandExit() }, //anon function that takes required params and calls niladic func
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
			}, //imp
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

func CommandHelp(commands map[string]cliCommand) { //when do we print the pokedex > thing? somewhere in main?
	var loopVar int //there's gotta be a more elegant way to do this
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

func CommandMap(areas LocationAreaBatch) LocationAreaBatch { //state?
	if FirstCall {
		for i := range areas.Results {
			fmt.Println(areas.Results[i].Name)
			FirstCall = false
			return areas
		}
	}
	if areas.Next != "null" && areas.Next != "" { //may need to start returning stuff to update the state holding struct
		next, err := ParseLocationAreas(areas.Next, DataStore)
		if err != nil {
			fmt.Printf("%v oopsie when parsing jason", err)
		}
		for i := range next.Results {
			fmt.Println(next.Results[i].Name)
		}
		return next
	}
	fmt.Println("you've reached the end of the world, we can't go any further!")
	return areas
}

func CommandMapB(areas LocationAreaBatch) LocationAreaBatch {
	if FirstCall {
		fmt.Println("we're back at the start!")
		FirstCall = false
		return areas
	}
	if areas.Previous != "null" && areas.Previous != "" { //maybe a channel could wait to get next 20?
		last, err := ParseLocationAreas(areas.Next, DataStore)
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

//current ideas - implement concurrency when generating map - store more information
//write more extensive testing

// if batches.Next != "null" && batches.Next != "" { //maybe a channel could wait to get next 20?
// 	next, err := ParseLocationAreas(batches.Next, c)
// 	if err != nil {
// 		return nil, errors.New("error parsing json")
// 	}
// 	JsonDB = append(JsonDB, next...)
// }

/* READMEhave made the following modifications -> no longer parses entire api tree, should just make requests as needed
stored global state through package LocationAreaState declaration, implemented untested caching and concurrent elements,
TODO -> test basic parsing and api traversal, then test caching, then implement concurrency */
