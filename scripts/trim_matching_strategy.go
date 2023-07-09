package scripts

import (
	"log"
	"time"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/env"
	"github.com/bujosa/xihe/utils"
)

func TrimMatchingStrategy() {
	utils.SetLogFile("trim_matching_strategy.txt")
	log.Println("Starting Trim Matching Strategy...")

	cars := database.GetCars()

	countryVersion, err := env.GetString("COUNTRY_VERSION_ID")
	if err != nil {
		panic(err)
	}

	category, err := env.GetString("CATEGORY_ID")
	if err != nil {
		panic(err)
	}

	for _, car := range cars {
		log.Println("Car: " + car.Id)

		price := utils.ConvertPrice(car.Price, car.Currency)

		createCatInput := api.CreateCarInput{
			TrimLevel:        car.Trim,
			InteriorColor:    car.InteriorColor,
			ExteriorColor:    car.ExteriorColor,
			MainPicture:      car.MainPicture,
			ExteriorPictures: []string{},
			InteriorPictures: []string{},
			Mileage:          car.Mileage,
			CountryVersion:   countryVersion,
			LicensePlate:     car.LicensePlate,
			Categories:       []string{category},
			Dealer:           car.DealerId,
			Provider:         "dealer",
			CurboSpot:        car.Spot,
			PriceInfo: api.CreatePriceInfoInput{
				BasePrice: price,
				Fee:       utils.FEE,
				Transfer:  utils.TRANSFER,
			},
		}

		if !car.PicturesUploaded {
			UploadPictures(car, &createCatInput)
			utils.SetLogFile("trim_matching_strategy.log")
		} else {
			log.Println("Pictures already uploaded for car: " + car.Id)
		}

		carUploaded, status := api.CreateCar(createCatInput, car.Id)

		updateCarInfo := database.UpdateCarInfo{
			Car:              car,
			Status:           status,
			MatchingStrategy: "trim",
			PicturesUploaded: true,
		}

		if status != "success" {
			log.Println("Error creating car: " + car.Id)
			updateCarInfo.NewId = ""
			database.UpdateCar(updateCarInfo)
			continue
		}

		updateCarInfo.NewId = carUploaded.Id
		database.UpdateCar(updateCarInfo)

		time.Sleep(30 * time.Second)
	}
}
