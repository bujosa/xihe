package scripts

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/storage"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func TrimMatchingStrategy(ctx context.Context, dealerPublished bool) {
	log.Println("Starting Trim Matching Strategy...")

	countryVersion := ctx.Value(utils.CountryVersionId).(string)

	category := ctx.Value(utils.CategoryId).(string)

	storage := storage.New(ctx)

	var cars []database.Car

	if dealerPublished {
		cars = database.GetCarsWithPublishDealers(ctx)
	} else {
		cars = database.GetCars(ctx)
	}

	missingUpdatedCars := []database.UpdateCarInfo{}

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
			err := UploadPictures(storage, car, &createCarInput)
			if err != nil {
				log.Println("Error uploading pictures for car: " + car.Id)
				time.Sleep(5 * time.Second)
				continue
			}
		} else {
			log.Println("Pictures already uploaded for car: " + car.Id)
		}

		if createCarInput.Mileage == 0 {
			createCarInput.Mileage = 1
		}

		carUploaded, status := api.CreateCar(ctx, createCarInput, car.Id)

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
			result := database.UpdateCar(ctx, updateCarInfo)
			if !result {
				missingUpdatedCars = append(missingUpdatedCars, updateCarInfo)
			}
			time.Sleep(5 * time.Second)
			continue
		}

		updateCarInfo.Set["uploaded"] = true
		result := database.UpdateCar(ctx, updateCarInfo)
		if !result {
			missingUpdatedCars = append(missingUpdatedCars, updateCarInfo)
		}
		time.Sleep(5 * time.Second)
	}

	if len(missingUpdatedCars) > 0 {
		for i, car := range missingUpdatedCars {
			result := database.UpdateCar(ctx, car)
			if !result {
				log.Println("Error updating car: " + car.Car.Id)
			} else {
				log.Println("Car updated: " + car.Car.Id)
				missingUpdatedCars = append(missingUpdatedCars[:i], missingUpdatedCars[i+1:]...)
			}
		}
	}

	if len(missingUpdatedCars) > 0 {
		// Add the remaining cars into a output file
		log.Println("There are " + string(rune(len(missingUpdatedCars))) + " cars that could not be updated")

		nameFile := "error_updating_" + time.Now().Format("2006-01-02_15-04-05") + ".txt"
		file, err := os.Create("logs/" + nameFile)
		if err != nil {
			log.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		log.Println("There are", len(missingUpdatedCars), "cars that could not be updated")

		for _, car := range missingUpdatedCars {
			data, err := bson.Marshal(car.Set)
			if err != nil {
				log.Println("Error marshaling car data:", err)
				continue
			}

			line := car.Car.Id + "," + string(data) + "\n"
			_, err = io.WriteString(file, line)
			if err != nil {
				log.Println("Error writing to file:", err)
				return
			}
		}

		log.Println("The list of cars that could not be updated has been saved to error_updating.txt")
	}

}
