package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/whenitsdone1/pokedex/util"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin) //create a standard input scanner
	commands := util.NewCommandMap()      //initalize the commands map
	fmt.Print("pokedex > ")
	for { //event loop
		if scanner.Scan() { //if we recieve input
			input := scanner.Text()
			sanitized := util.SanitizeInput(input)
			entered := strings.Fields(sanitized)
			util.ParseCommand(entered, commands) //clean input and find the command the input corresponds with
			fmt.Print("pokedex > ")
		} else { //failure state, if scanner fails
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "error reading input", err)
			}
		}
	}
}
