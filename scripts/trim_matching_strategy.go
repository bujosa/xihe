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

var missingUpdatedCars []database.UpdateCarInfo
var errorCount int = 0

func TrimMatchingStrategy(ctx context.Context, dealerPublished bool) {
	log.Println("Starting Trim Matching Strategy...")
	countryVersion := ctx.Value(utils.CountryVersionId).(string)
	category := ctx.Value(utils.CategoryId).(string)
	storage := storage.New(ctx)
	updatingMissingUpdatedCars(ctx)

	var cars []database.Car

	if dealerPublished {
		cars = database.GetCarsWithPublishDealers(ctx)
	} else {
		cars = database.GetCars(ctx)
	}

	for _, car := range cars {
		time.Sleep(2 * time.Second)
		log.Println("Car: " + car.Id)

		createCarInput := formatCreateCarInput(countryVersion, category, car)

		if errorCount > 5 {
			errorCount = 0
			log.Println("Waiting 10 minutes to continue...")
			utils.CleanDns()
			time.Sleep(10 * time.Minute)
			TrimMatchingStrategy(ctx, dealerPublished)
		}

		if !car.PicturesUploaded {
			err := UploadPictures(storage, car, &createCarInput)
			if err != nil {
				log.Println("Error uploading pictures for car: " + car.Id)
				errorCount++
				continue
			}
		} else {
			log.Println("Pictures already uploaded for car: " + car.Id)
			errorCount = 0
		}

		carUploaded, status := api.CreateCar(ctx, createCarInput, car.Id)

		updateCar(ctx, car, createCarInput, carUploaded, status)
	}

	updatingMissingUpdatedCars(ctx)

	log.Println("Trim Matching Strategy finished")
}

// FormatCreateCarInput formats the input to create a car
func formatCreateCarInput(countryVersionId string, category string, car database.Car) api.CreateCarInput {
	price := utils.ConvertPrice(car.Price, car.Currency)

	createCarInput := api.CreateCarInput{
		TrimLevel:        car.Trim,
		InteriorColor:    car.InteriorColor,
		ExteriorColor:    car.ExteriorColor,
		MainPicture:      car.MainPicture,
		ExteriorPictures: car.ExteriorPictures,
		InteriorPictures: car.InteriorPictures,
		Mileage:          car.Mileage + 1,
		CountryVersion:   countryVersionId,
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

	return createCarInput
}

// Update car
func updateCar(ctx context.Context, car database.Car, createCarInput api.CreateCarInput, carUploaded api.Car, status utils.StatusRequest) {
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
			errorCount++
		}
		return
	}

	updateCarInfo.Set["uploaded"] = true
	result := database.UpdateCar(ctx, updateCarInfo)
	if !result {
		missingUpdatedCars = append(missingUpdatedCars, updateCarInfo)
		errorCount++
		return
	}
	errorCount = 0
}

// Update missing cars
func updatingMissingUpdatedCars(ctx context.Context) {
	newArray := []database.UpdateCarInfo{}
	if len(missingUpdatedCars) > 0 {
		log.Print("Updating " + string(rune(len(missingUpdatedCars))) + " cars that could not be updated")
		for _, car := range missingUpdatedCars {
			result := database.UpdateCar(ctx, car)
			if !result {
				log.Println("Error updating car: " + car.Car.Id)
			} else {
				log.Println("Car updated: " + car.Car.Id)
				newArray = append(newArray, car)
			}
		}
	}

	missingUpdatedCars = newArray

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
