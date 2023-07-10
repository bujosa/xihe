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

func TrimMatchingStrategy() {
	log.Println("Starting Trim Matching Strategy...")

	countryVersion, err := env.GetString("COUNTRY_VERSION_ID")
	if err != nil {
		panic(err)
	}

	category, err := env.GetString("CATEGORY_ID")
	if err != nil {
		panic(err)
	}

	cars := database.GetCars()

	for _, car := range cars {
		log.Println("Car: " + car.Id)

		price := utils.ConvertPrice(car.Price, car.Currency)

		createCarInput := api.CreateCarInput{
			TrimLevel:        car.Trim,
			InteriorColor:    car.InteriorColor,
			ExteriorColor:    car.ExteriorColor,
			MainPicture:      car.MainPicture,
			ExteriorPictures: car.ExteriorPictures,
			InteriorPictures: car.InteriorPictures,
			Mileage:          car.Mileage,
			CountryVersion:   countryVersion,
			LicensePlate:     car.LicensePlate,
			Categories:       []string{category},
			Dealer:           car.DealerId,
			Provider:         "DEALER",
			CurboSpot:        car.Spot,
			Status:           "AVAILABLE",
			PriceInfo: api.CreatePriceInfoInput{
				BasePrice: price,
				Fee:       utils.FEE,
				Transfer:  utils.TRANSFER,
			},
			VinNumber: "00000000000000000",
		}

		if !car.PicturesUploaded {
			UploadPictures(car, &createCarInput)
		} else {
			log.Println("Pictures already uploaded for car: " + car.Id)
		}

		carUploaded, status := api.CreateCar(createCarInput, car.Id)

		updateCarInfo := database.UpdateCarInfo{
			Car: car,
			Set: bson.M{
				"status":           status,
				"matchingStrategy": "trim",
				"picturesUploaded": true,
				"mainPicture":      createCarInput.MainPicture,
				"exteriorPictures": createCarInput.ExteriorPictures,
				"interiorPictures": createCarInput.InteriorPictures,
				"newId":            carUploaded.Id,
			},
		}

		if status != "success" {
			log.Println("Error creating car: " + car.Id)
			database.UpdateCar(updateCarInfo)
			continue
		}

		database.UpdateCar(updateCarInfo)
		time.Sleep(5 * time.Second)
	}
}
