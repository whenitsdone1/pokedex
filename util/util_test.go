package util

import (
	"fmt"
	"testing"
)

func TestParseLocationAreas(t *testing.T) {
	t.Run("parse fields correctly", func(t *testing.T) {
		url := "https://pokeapi.co/api/v2/location-area"
		locations, err := ParseLocationAreas(url, DataStore)

		if err != nil {
			t.Fatalf("unexpected error %v encounterd", err)
		}

		for _, element := range locations.Results { //ingesting the first page properly, but isnt yet getting every page
			fmt.Println(element.Url)
		}
		fmt.Println(len(locations.Results))
	})
}

// 	got := locations
// 	// if want != 19 {
// 	// 	t.Errorf("incorrect parsing - got length %v", want)
// 	// }
// 	fmt.Printf("The name is %s", got.Next)
// })
