package scripts

import (
	"log"
	"time"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/env"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func UploadDealers() {
	// Get dealers from database
	log.Println("Starting dealer upload...")
	dealers := database.GetDealers()

	geoCode, err := env.GetString("GEOCODE")
	if err != nil {
		panic(err)
	}

	// Upload dealers to api
	for _, dealer := range dealers {

		if dealer.TelephoneNumberSanitized == "" {
			dealer.TelephoneNumberSanitized = "8090000000"
		}

		dealer.TelephoneNumberSanitized = utils.TransformTelephoneNumber(dealer.TelephoneNumberSanitized)

		createDealerInput := api.CreateDealerInput{
			Name:            dealer.Name,
			Address:         geoCode + " " + utils.ReplaceNewLine(dealer.Address),
			Latitude:        dealer.Latitude,
			Longitude:       dealer.Longitude,
			City:            dealer.City,
			Spot:            dealer.Spot,
			TelephoneNumber: dealer.TelephoneNumberSanitized,
		}

		newDealer, err := api.CreateDealer(createDealerInput, dealer.Id)
		if err == "failed" {
			continue
		}

		// Update dealer
		updateDealerInfo := database.UpdateDealerInfo{
			Id: dealer.Id,
			Set: bson.M{
				"uploaded":                 true,
				"dealer":                   newDealer.Id,
				"telephoneNumberSanitized": dealer.TelephoneNumberSanitized,
			},
		}

		database.UpdateDealer(updateDealerInfo)

		time.Sleep(3 * time.Second)
	}
}
