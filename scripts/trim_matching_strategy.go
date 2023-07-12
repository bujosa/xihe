package scripts

import (
	"context"
	"log"
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
			database.UpdateCar(ctx, updateCarInfo)
			time.Sleep(4 * time.Second)
			continue
		}

		updateCarInfo.Set["uploaded"] = true
		database.UpdateCar(ctx, updateCarInfo)
		time.Sleep(4 * time.Second)
	}
}
