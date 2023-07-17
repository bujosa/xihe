package scripts

import (
	"context"
	"log"
	"time"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func UploadTrims(ctx context.Context) {
	log.Println("Starting Upload Trims...")

	// In future I will match features cars with features database appart from newFeatures
	features := database.GetFeatures(ctx)

	newFeatures := []string{}

	for i := 0; i < len(features); i++ {
		if features[i].Slug == "bolsa-de-aire-de-chofer" || features[i].Slug == "aire-acondicianado" {
			newFeatures = append(newFeatures, features[i].Id)
			features = append(features[:i], features[i+1:]...)
			i--
		}
	}

	cars := database.BaseGetCars(ctx, database.Filter{
		Match: bson.M{
			"setTrimName": true,
			"trimMatched": false,
		}})

	for _, car := range cars {

		price := utils.ConvertPrice(car.Price, car.Currency)

		createTrimInput := api.CreateTrimInput{
			Name:         car.TrimName,
			Alias:        car.TrimName,
			Model:        car.ModelId,
			Year:         car.Year,
			BodyStyle:    car.BodyStyleId,
			Transmission: car.TransmissionId,
			DriveTrain:   car.DriveTrainId,
			FuelType:     car.FuelTypeId,
			Features:     newFeatures,
			Status:       "PUBLISHED",
			Price:        price,
		}

		_, err := api.CreateTrim(ctx, createTrimInput)

		time.Sleep(500 * time.Millisecond)

		if err == "failed" {
			log.Println("Error creating trim:", err)
			continue
		}
	}

}
