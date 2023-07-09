package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bujosa/xihe/env"
)

type GeocodeResponse struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func Geocode(address string) (float64, float64, error) {
	key, err := env.GetString("GOOGLE_MAPS_API_KEY")
	if err != nil {
		return 0, 0, err
	}

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", url.QueryEscape(address), key)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var data GeocodeResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, 0, err
	}

	print(data.Results)

	if len(data.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found")
	}

	return data.Results[0].Geometry.Location.Lat, data.Results[0].Geometry.Location.Lng, nil
}
