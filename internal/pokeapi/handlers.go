package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) GenericGetList(pageURL *string, path  string) (PaginationResponse, error) {
    var url string
    if pageURL != nil {
        url = *pageURL
    } else {
        url = baseURL + path
    }
	
	if cachedData, found := c.cache.Get(url); found {
		res := PaginationResponse{}
		err := json.Unmarshal(cachedData, &res)
		if err == nil {
			return res, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PaginationResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PaginationResponse{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return PaginationResponse{}, err
	}

	res := PaginationResponse{}
	err = json.Unmarshal(dat, &res)
	if err != nil {
		return PaginationResponse{}, err
	}

	c.cache.Add(url, dat)
	return res, nil
}

func (c *Client) GetLocation(locationName string) (Location, error) {
	url := baseURL + "/location-area/" + locationName

	if val, ok := c.cache.Get(url); ok {
		locationResp := Location{}
		err := json.Unmarshal(val, &locationResp)
		if err != nil {
			return Location{}, err
		}
		return locationResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Location{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Location{}, err
	}

	locationResp := Location{}
	err = json.Unmarshal(dat, &locationResp)
	if err != nil {
		return Location{}, err
	}

	c.cache.Add(url, dat)

	return locationResp, nil
}

func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	if val, ok := c.cache.Get(url); ok {
		res := Pokemon{}
		err := json.Unmarshal(val, &res)
		if err != nil {
			return Pokemon{}, err
		}
		return res, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	res := Pokemon{}
	err = json.Unmarshal(dat, &res)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, dat)

	return res, nil
}