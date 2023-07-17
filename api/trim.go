package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bujosa/xihe/utils"
)

type CreateTrimInput struct {
	Alias        string   `json:"alias"`
	Name         string   `json:"name"`
	Status       string   `json:"status"`
	Model        string   `json:"carModel"`
	FuelType     string   `json:"fuelType"`
	DriveTrain   string   `json:"driveTrain"`
	Features     []string `json:"features"`
	BodyStyle    string   `json:"bodyStyle"`
	Transmission string   `json:"transmission"`
	Year         int      `json:"year"`
	Price        int      `json:"price"`
}

type Trim struct {
	Id string `json:"id"`
}

type CreateTrimResponse struct {
	Data struct {
		CreateTrim Trim `json:"createTrimLevel"`
	} `json:"data"`
}

func CreateTrim(ctx context.Context, createTrimInput CreateTrimInput) (Trim, utils.StatusRequest) {
	log.Println("Creating model... with Name: " + createTrimInput.Name)

	url := ctx.Value(utils.ProductionApiUrl).(string)
	token := ctx.Value(utils.SessionSecret).(string)

	mutation := `
		mutation CreateTrimLevel($input: CreateTrimLevelInput!) {
			createTrimLevel(input: $input) {
				id
			}
		}
	`

	requestBody, err := json.Marshal(map[string]interface{}{
		"query": mutation,
		"variables": map[string]interface{}{
			"input": createTrimInput,
		},
	})

	if err != nil {
		log.Printf("Error marshalling request body %s\n", err)
		return Trim{}, utils.StatusRequest("failed")
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating HTTP request %s\n", err)
		return Trim{}, utils.StatusRequest("failed")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request %s\n", err)
		return Trim{}, utils.StatusRequest("failed")
	}
	defer response.Body.Close()

	var responseData CreateTrimResponse

	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		log.Printf("Error decoding response body %s\n", err)
		return Trim{}, utils.StatusRequest("failed")
	}

	if response.StatusCode != 200 {
		log.Printf("Error status code %d\n", response.StatusCode)
		log.Printf("Error data response %s\n", response.Body)
		return Trim{}, utils.StatusRequest("failed")
	} else {
		if responseData.Data.CreateTrim.Id == "" {
			log.Printf("Error creating Trim %s\n", err)
			return Trim{}, utils.StatusRequest("failed")
		}
		log.Printf("Trim created with ID: %s\n", responseData.Data.CreateTrim.Id)
		return responseData.Data.CreateTrim, utils.StatusRequest("success")
	}
}
