package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bujosa/xihe/env"
	"github.com/bujosa/xihe/utils"
)

type CreateDealerInput struct {
	CurboSpot        string               `json:"curboSpot"`
	Name    		 string               `json:"name"`
	Adress   		 string               `json:"adress"`
	Latitude  		 float64              `json:"latitude"`
	Longitude 		 float64              `json:"longitude"`
	City   		 	 string               `json:"city"`
}

type Dealer struct {
	Id string `json:"id"`
}

type CreateDealerResponse struct {
	Data struct {
		CreateDealer Dealer `json:"createDealer"`
	} `json:"data"`
}

func CreateDealer(createDealerInput CreateDealerInput, id string) (Dealer, utils.StatusRequest) {
	log.Println("Creating dealer... with ID: " + id)

	url, err:= env.GetString("PRODUCTION_API_URL")
	if err != nil {
		panic(err)
	}
	token, err := env.GetString("SESSION_SECRET")
	if err != nil {
		panic(err)
	}

	mutation := `
		mutation CreateDealer($input: CreateDealerInput!) {
			createDealer(input: $input) {
				id
			}
		}
	`

	request := GraphqlRequest{
		Query: mutation,
		Variables: map[string]interface{}{
			"input": createDealerInput,
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshalling request body %s\n", err)
		return Dealer{}, utils.StatusRequest("failed")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating HTTP request %s\n", err)
		return Dealer{}, utils.StatusRequest("failed")
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request %s\n", err)
		return Dealer{}, utils.StatusRequest("failed")
	}
	defer response.Body.Close()

	var responseData CreateDealerResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		log.Printf("Error decoding response body %s\n", err)
		return Dealer{}, utils.StatusRequest("failed")
	}

	if response.StatusCode != 200 {
		log.Printf("Error status code %d\n", response.StatusCode)
		log.Printf("Error data response %s\n", response.Body)
		return Dealer{}, utils.StatusRequest("failed")
	} else {
		log.Printf("Dealer created with ID: %s\n", responseData.Data.CreateDealer.Id)
		return responseData.Data.CreateDealer, utils.StatusRequest("success")
	}
}