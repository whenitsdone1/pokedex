package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type LocationAreaBatch struct {
	Next     string         `json:"next"`
	Results  []LocationArea `json:"results"`
	Previous string         `json:"previous"`
}

const (
	LocationAreaApiUrl = "https://pokeapi.co/api/v2/location-area"
	//LocationApiUrl     = "https://pokeapi.co/api/v2/location/" - not needed yet
)

type LocationArea struct {
	Name string `json:"name"` //will capture more information as needed
	Url  string `json:"url"`
}

func GetJson(toParse string) ([]byte, error) {
	response, err := http.Get(toParse)
	if err != nil {
		return nil, errors.New("error making get request")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading body")
	}
	return body, nil
}

func ParseLocationAreas(toParse string, c *Cache) (LocationAreaBatch, error) { //we need to update next and previous, so need to return LocationAreaBatch
	if val, ok := c.Get(toParse); ok {
		fmt.Println("Using the cache!")
		return val, nil
	}
	fmt.Println("Could not get from Cache, fetching...")
	Json, err := GetJson(toParse)
	if err != nil {
		return LocationAreaBatch{}, err
	}
	var batches LocationAreaBatch
	err = json.Unmarshal(Json, &batches)

	if err != nil {
		return LocationAreaBatch{}, errors.New("error parsing json")
	}

	slices.SortFunc(batches.Results, CompareLocations)
	c.Add(toParse, batches) // add to cache
	fmt.Println("Adding to cache..")
	return batches, nil

}

func CompareLocations(x, y LocationArea) int { //comparison function
	a, errA := strconv.Atoi(getNumber(x.Url))
	if errA != nil {
		log.Fatal("error in parsing of URLs")
	}
	b, errB := strconv.Atoi(getNumber(y.Url))
	if errB != nil {
		log.Fatal("error in parsing of URLs")
	}
	return a - b
}
func getNumber(x string) string { //retrieve the number of the location area from the url
	x = strings.TrimPrefix(x, "https://pokeapi.co/api/v2/location-area/")
	x = strings.TrimSuffix(x, "/")
	return x
}
