package api

import (
	"bytes"
	"encoding/json"
	"log"
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
	Dealer           string               `json:"dealer"`
	VinNumber        string               `json:"vinNumber"`
}

type CreatePriceInfoInput struct {
	BasePrice int `json:"basePrice"`
	Fee       int `json:"fee"`
	Transfer  int `json:"transfer"`
}

type Car struct {
	Id string `json:"id"`
}

type CreateCarResponse struct {
	Data struct {
		CreateCar Car `json:"createCar"`
	} `json:"data"`
}

func CreateCar(createCarInput CreateCarInput, id string) (Car, utils.StatusRequest) {
	log.Println("Creating car... with ID: " + id)

	url, err := env.GetString("PRODUCTION_API_URL")
	if err != nil {
		panic(err)
	}
	token, err := env.GetString("SESSION_SECRET")
	if err != nil {
		panic(err)
	}

	mutation := `
		mutation CreateCar($input: CreateCarInput!) {
			createCar(input: $input) {
				id
			}
		}
	`
	requestBody, err := json.Marshal(map[string]interface{}{
		"query": mutation,
		"variables": map[string]interface{}{
			"input": createCarInput,
		},
	})

	if err != nil {
		log.Printf("Error marshalling request body %s\n", err)
		return Car{}, utils.StatusRequest("failed")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating HTTP request %s\n", err)
		return Car{}, utils.StatusRequest("failed")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request %s\n", err)
		return Car{}, utils.StatusRequest("failed")
	}
	defer response.Body.Close()

	// responseReader := response.Body
	// body, err := ioutil.ReadAll(responseReader)
	// if err != nil {
	// 	log.Printf("Error reading response body %s\n", err)
	// 	return Car{}, utils.StatusRequest("failed")
	// }
	// log.Println(string(body))

	var responseData CreateCarResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		log.Printf("Error decoding response body %s\n", err)
		return Car{}, utils.StatusRequest("failed")
	}

	if response.StatusCode != 200 {
		log.Printf("Error status code %d\n", response.StatusCode)
		return Car{}, utils.StatusRequest("failed")
	} else if responseData.Data.CreateCar.Id == "" {
		log.Printf("Error creating car %s\n", err)
		return Car{}, utils.StatusRequest("failed")
	} else {
		log.Printf("Car created with ID: %s\n", responseData.Data.CreateCar.Id)
		return responseData.Data.CreateCar, utils.StatusRequest("success")
	}
}
