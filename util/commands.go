package util

import (
	"errors"
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
var InitMap, _ = initMap() //how do we check errors here?

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
			callback:    func(m map[string]cliCommand) { CommandMap(InitMap) }, //implement
		},
		"mapb": {
			name:        "map back",
			description: "display the previous locations",
			callback:    func(m map[string]cliCommand) { CommandMapB(InitMap) }, //imp
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

func initMap() ([]LocationArea, error) {
	InitalMap, err := ParseLocationAreas(LocationAreaApiUrl)
	if err != nil {
		return nil, errors.New("error parsing json")
	}
	return InitalMap, nil
}

func CommandMap(areas []LocationArea) {
	if Pointer.Value >= len(InitMap) {
		fmt.Println("you've reached the end of the world, we can't go any further!")
		return
	}
	for i := Pointer.Value; i < Pointer.Value+20; i++ {
		fmt.Println(areas[i].Name)
	}
	if (Pointer.Value + 20) < len(InitMap) {
		Pointer.Advance()
	}
}

func CommandMapB(areas []LocationArea) {
	if Pointer.Value == 0 {
		fmt.Println("we're back at the start!")
		return
	}
	if (Pointer.Value - 20) <= 0 { //handle index error in the negative case
		Pointer.Value = 0
	} else {
		Pointer.Recede()
	}
	for i := Pointer.Value; i < Pointer.Value+20; i++ {
		fmt.Println(areas[i].Name)
	}
}

//current ideas - implement concurrency when generating map - store more information
//write more extensive testing
