package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bujosa/xihe/utils"
)

type CreateModelInput struct {
	Name   string `json:"name"`
	Brand  string `json:"brand"`
	Status string `json:"status"`
}

type Model struct {
	Id string `json:"id"`
}

type CreateModelResponse struct {
	Data struct {
		CreateModel Model `json:"createModel"`
	} `json:"data"`
}

func CreateModel(ctx context.Context, createModelInput CreateModelInput) (Model, utils.StatusRequest) {
	log.Println("Creating model... with Name: " + createModelInput.Name)

	url := ctx.Value(utils.ProductionApiUrl).(string)
	token := ctx.Value(utils.SessionSecret).(string)

	mutation := `
		mutation CreateModel($input: CreateModelInput!) {
			createModel(input: $input) {
				id
			}
		}
	`

	requestBody, err := json.Marshal(map[string]interface{}{
		"query": mutation,
		"variables": map[string]interface{}{
			"input": createModelInput,
		},
	})

	if err != nil {
		log.Printf("Error marshalling request body %s\n", err)
		return Model{}, utils.StatusRequest("failed")
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating HTTP request %s\n", err)
		return Model{}, utils.StatusRequest("failed")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request %s\n", err)
		return Model{}, utils.StatusRequest("failed")
	}
	defer response.Body.Close()

	var responseData CreateModelResponse

	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		log.Printf("Error decoding response body %s\n", err)
		return Model{}, utils.StatusRequest("failed")
	}

	if response.StatusCode != 200 {
		log.Printf("Error status code %d\n", response.StatusCode)
		log.Printf("Error data response %s\n", response.Body)
		return Model{}, utils.StatusRequest("failed")
	} else {
		if responseData.Data.CreateModel.Id == "" {
			log.Printf("Error creating model %s\n", err)
			return Model{}, utils.StatusRequest("failed")
		}
		log.Printf("Model created with ID: %s\n", responseData.Data.CreateModel.Id)
		return responseData.Data.CreateModel, utils.StatusRequest("success")
	}
}
