package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var UnknownKeyError error

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
		// 	"map":{
		// 		name: "map",
		// 		description: "display the next locations",
		// 		callback: func(m map[string]cliCommand) {}, //implement
		// 	}
		// 	"mapb":{
		// 		name: "map back",
		// 		description: "display the previous locations",
		// 		callback: func(m map[string]cliCommand) {}, //imp
		// 	}
		// }
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

func CommandMap() {} //implement

func CommandMapB() {} //implement
