package util

import (
	"fmt"
	"testing"
	"time"
)

func TestParseLocationAreas(t *testing.T) {
	t.Run("parse fields correctly", func(t *testing.T) {
		url := "https://pokeapi.co/api/v2/location-area"
		locations, err := ParseLocationAreas(url, DataStore)

		if err != nil {
			t.Fatalf("unexpected error %v encounterd", err)
		}

		for _, element := range locations.Results {
			fmt.Println(element.Url)
		}
		fmt.Println(len(locations.Results))
	})
}
func TestCache(t *testing.T) { //TODO: Implement better test coverage
	Check := NewCache(100 * time.Hour)
	baseURL := "https://pokeapi.co/api/v2/location-area"

	urls := []string{
		baseURL,
		baseURL + "?offset=20&limit=20",
		baseURL + "?offset=40&limit=20",
	}

	for _, url := range urls {
		for i := 0; i < 2; i++ {
			result, err := ParseLocationAreas(url, Check)
			if err != nil {
				t.Fatalf("Error parsing location areas: %v", err)
			}
			fmt.Printf("Parsed URL: %s, Next URL: %s\n", url, result.Next)
		}
	}

	fmt.Printf("\nNumber of entries in cache: %d\n", len(Check.Entries))

	for key, value := range Check.Entries {
		fmt.Printf("\nKey: %s\n", key)
		fmt.Printf("Created At: %s\n", value.createdAt.Format(time.RFC3339))
		fmt.Printf("Number of results: %d\n", len(value.val.Results))
		fmt.Printf("Next URL: %s\n", value.val.Next)
		fmt.Printf("First location area: %s\n", value.val.Results[0].Name)
	}
}
