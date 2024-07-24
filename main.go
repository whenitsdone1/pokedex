package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/whenitsdone1/pokedex/util"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin) // pass a reader

	commands := util.NewCommandMap()
	fmt.Print("pokedex > ")

	for { //while true...

		if scanner.Scan() {
			input := scanner.Text()
			sanitized := util.SanitizeInput(input)
			util.ParseCommand(sanitized, commands)
			fmt.Print("\npokedex > ")
		} else {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "error reading input", err)
			}
		}
	}
}
