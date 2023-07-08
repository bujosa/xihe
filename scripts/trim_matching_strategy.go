package scripts

import (
	"log"
	"os"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/env"
)

func TrimMatchingStrategy() {
	// Add register file
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.SetPrefix("[INFO] ")
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

		createCatInput := api.CreateCarInput{
			TrimLevel:        car.Trim,
			InteriorColor:    car.InteriorColor,
			ExteriorColor:    car.ExteriorColor,
			MainPicture:      car.MainPicture,
			ExteriorPictures: []string{},
			InteriorPictures: []string{},
			Mileage: 		car.Mileage,
			CountryVersion: 	countryVersion,
			LicensePlate: 	car.LicensePlate,
			Categories: 	[]string{category},
		}

		if !car.PicturesUploaded {
			UploadPictures(car, &createCatInput)
		} else {
			log.Println("Pictures already uploaded for car: " + car.Id)
		}

		carUploaded, status := api.CreateCar(createCatInput, car.Id)

		updateCarInfo := database.UpdateCarInfo{
			Car: car,
			Status: status,
			MatchingStrategy: "trim",
			PicturesUploaded: true,
		}

		if status != "success" {
			log.SetPrefix("[ERROR] ")
			log.Println("creating car: " + car.Id)
			updateCarInfo.NewId = ""
			database.UpdateCar(updateCarInfo)
			continue
		}

		updateCarInfo.NewId = carUploaded.Id
		database.UpdateCar(updateCarInfo)
	}
}