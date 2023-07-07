package grapqhql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bujosa/xihe/env"
	"github.com/bujosa/xihe/utils"
)

type CreateCarInput struct {
	TrimLevel        string               `json:"trimLevel"`
	InteriorColor    string               `json:"interiorColor"`
	ExteriorColor    string               `json:"exteriorColor"`
	MainPicture      string               `json:"mainPicture"`
	ExteriorPictures []string             `json:"exteriorPictures"`
	InteriorPictures []string             `json:"interiorPictures"`
	CurboSpot        string               `json:"curboSpot"`
	Provider         string               `json:"provider"`
	Status           string               `json:"status"`
	PriceInfo        CreatePriceInfoInput `json:"priceInfo"`
	Mileage          int                  `json:"mileage"`
	CountryVersion   string               `json:"countryVersion"`
	LicensePlate     string               `json:"licensePlate"`
	Categories       []string             `json:"categories"`
}

type CreatePriceInfoInput struct {
	BasePrice float64 `json:"basePrice"`
	Fee       float64 `json:"fee"`
	Transfer  float64 `json:"transfer"`
}

type Car struct {
	Id string `json:"id"`
}

type CreateCarResponse struct {
	Data struct {
		CreateCar Car `json:"createCar"`
	} `json:"data"`
}

func CreateCar(createCarInput CreateCarInput) (Car, utils.StatusRequest) {
	url, err:= env.GetString("PRODUCTION_API_URL")
	if err != nil {
		panic(err)
	}
	token, err := env.GetString("SESSION_SECRET")
	if err != nil {
		panic(err)
	}

	mutation := `
		mutation CreateCar($input: CreateCarInput!) {
			createCar(input: $createCarInput) {
				id
			}
		}
	`

	request := GraphqlRequest{
		Query: mutation,
		Variables: map[string]interface{}{
			"input": createCarInput,
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Error marshalling request body %s\n", err)
		return Car{}, utils.StatusRequest("failed")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error creating HTTP request %s\n", err)
		return Car{}, utils.StatusRequest("failed")
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request %s\n", err)
		return Car{}, utils.StatusRequest("failed")
	}
	defer response.Body.Close()

	var responseData CreateCarResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		fmt.Printf("Error decoding response body %s\n", err)
		return Car{}, utils.StatusRequest("failed")
	}

	if response.StatusCode != 200 {
		fmt.Printf("Error status code %d\n", response.StatusCode)
		return Car{}, utils.StatusRequest("failed")
	} else {
		return responseData.Data.CreateCar, utils.StatusRequest("success")
	}
}