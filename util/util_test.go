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
		if locationBatch, ok := value.val.(*LocationAreaBatch); ok {
			fmt.Printf("Number of results: %d\n", len(locationBatch.Results))
			fmt.Printf("Next URL: %s\n", locationBatch.Next)
			fmt.Printf("First location area: %s\n", locationBatch.Results[0].Name)
		}
	}
}

func TestLocationParsing(t *testing.T) {
	t.Run("Do we get anything?", func(t *testing.T) {
		try, catch := ParseLocations("https://pokeapi.co/api/v2/location-area/2/", DataStore)
		if catch != nil {
			t.Errorf("Caught %v", catch)

		}
		for _, e := range try {
			fmt.Printf("- %s\n", e.Name)
			fmt.Printf("@ %s\n", e.Url)
		}
		j := ExtractNames(try)
		fmt.Print(j)
	})
}
