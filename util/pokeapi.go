package util

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type LocationAreaBatch struct {
	Next     string         `json:"next"`
	Results  []LocationArea `json:"results"` //?
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
	InCache := func(key string) bool {
		_, ok := c.Entries[key] //call cache
		return ok
	}
	if InCache(toParse) {
		return c.Entries[toParse].val, nil
	}
	Json, _ := GetJson(toParse)
	var batches LocationAreaBatch
	err := json.Unmarshal(Json, &batches)

	if err != nil {
		var Zero LocationAreaBatch
		return Zero, errors.New("error parsing json")
	}
	slices.SortFunc(batches.Results, CompareLocations)
	c.Add(toParse, batches) // add to cache
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

//0. refactor like half the app so that a cache makes anysense
//1.fix cache
//2. dont download every result on startup

//idea - just return next and write to a channel so every next call parses next page rather than parsing as a big db
//back would need to access previous, will need to add a field to json parsing
