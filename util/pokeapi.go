package util

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type LocationAreaBatch struct {
	Next    string         `json:"next"`
	Results []LocationArea `json:"results"` //?
}

const (
	LocationAreaApiUrl = "https://pokeapi.co/api/v2/location-area"
	LocationApiUrl     = "https://pokeapi.co/api/v2/location/"
)

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

//call will need to be recursive

func GetNextLocations() {}

func GetPreviousLocations() {}

func LocationCache() {}

func ParseLocationAreas(toParse string) []LocationArea {
	response, err := http.Get(toParse)
	if err != nil {
		log.Fatalf("error making http request %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("error reading response body: %v", err)
	}
	var batches LocationAreaBatch

	err = json.Unmarshal(body, &batches)
	if err != nil {
		log.Fatalf("error parsing json %v", err)
	}

	var tempJson []LocationArea
	if batches.Next != "null" { //maybe a channel could wait to get next 20?
		tempJson = (parseCollection(batches.Results, tempJson))
		ParseLocationAreas(batches.Next)
	}
	tempJson = parseCollection(batches.Results, tempJson) //one more

	return tempJson

}

func parseCollection[T any](collection []T, result []T) []T {
	result = append(result, collection...) //very cool - ... passes elements one by one to the function!
	return result
}

////test this tmrw - my thinking is, first collect a big list of all listed locations using recursion - once I have that and understand whats happening, can do more stuff
