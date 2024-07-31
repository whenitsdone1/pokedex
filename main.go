package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/whenitsdone1/pokedex/util"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin) //TODO: Clean up main function
	commands := util.NewCommandMap()
	fmt.Print("pokedex > ")

	for { //while true...

		if scanner.Scan() {
			input := scanner.Text()
			sanitized := util.SanitizeInput(input)
			entered := strings.Fields(sanitized)
			util.ParseCommand(entered, commands)
			fmt.Print("\npokedex > ")
		} else {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "error reading input", err)
			}
		}
	}
}
