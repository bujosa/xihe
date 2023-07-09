package scripts

import (
	"log"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/env"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func UploadDealers() {
	// Get dealers from database
	utils.SetLogFile("dealer.txt")
	log.Println("Starting dealer upload...")
	dealers := database.GetDealers()

	geoCode, err := env.GetString("GEOCODE")
	if err != nil {
		panic(err)
	}

	// Upload dealers to api
	for _, dealer := range dealers {
		createDealerInput := api.CreateDealerInput{
			Name:      dealer.Name,
			Adress:    geoCode + " " + dealer.Address,
			Latitude:  dealer.Latitude,
			Longitude: dealer.Longitude,
			City:      dealer.City,
			Spot:      dealer.Spot,
		}

		newDealer, err := api.CreateDealer(createDealerInput, dealer.Id)
		if err == "failed" {
			continue
		}

		// Update dealer
		updateDealerInfo := database.UpdateDealerInfo{
			Id: dealer.Id,
			Set: bson.M{
				"uploaded": true,
				"dealer":   newDealer.Id,
			},
		}

		database.UpdateDealer(updateDealerInfo)
	}
}
