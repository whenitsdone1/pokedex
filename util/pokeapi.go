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
	Next    string         `json:"next"`
	Results []LocationArea `json:"results"` //?
}

const (
	LocationAreaApiUrl = "https://pokeapi.co/api/v2/location-area"
	//LocationApiUrl     = "https://pokeapi.co/api/v2/location/" - not needed yet
)

type LocationArea struct {
	Name string `json:"name"` //will capture more information as needed
	Url  string `json:"url"`
}

var JsonDB []LocationArea //should look at removing this - global variables == bad

func ParseLocationAreas(toParse string) ([]LocationArea, error) { //recursively get location areas and add to DB
	response, err := http.Get(toParse)
	if err != nil {
		return nil, errors.New("error making get request")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading body")
	}

	var batches LocationAreaBatch
	err = json.Unmarshal(body, &batches)

	if err != nil {
		return nil, errors.New("error parsing json")
	}

	JsonDB := append(JsonDB, batches.Results...)

	if batches.Next != "null" && batches.Next != "" { //maybe a channel could wait to get next 20?
		next, err := ParseLocationAreas(batches.Next)

		if err != nil {
			return nil, errors.New("error parsing json")
		}
		JsonDB = append(JsonDB, next...)
	}
	slices.SortFunc(JsonDB, CompareLocations)
	return JsonDB, nil

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
