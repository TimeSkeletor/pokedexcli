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

