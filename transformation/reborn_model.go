package transformation

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func RebornModel(ctx context.Context) {
	log.Print("Starting Reborn Model Script... This purifies the database from cars that are not in the model list")
	cars := database.BaseGetCars(ctx, database.Filter{
		Match: bson.M{
			"modelMatchLayer": 1,
			"modelMatched":    false,
			"setTrimName":     false,
		}})

	for _, car := range cars {
		modelName, trimName := fixModelAndTrimName(utils.EliminateExtraSpace(car.Model), car.Brand, car.Trim, car.Year)

		updateCarInfo := database.UpdateCarInfo{
			Car: car,
			Set: bson.M{
				"model":           modelName,
				"trimName":        trimName,
				"modelSlug":       utils.Slug([]string{car.Brand, modelName}),
				"setTrimName":     true,
				"modelMatchLayer": 5,
			},
		}

		database.UpdateCar(ctx, updateCarInfo)
	}

	Model(ctx, 5)

	newCars := database.BaseGetCars(ctx, database.Filter{
		Match: bson.M{
			"modelMatchLayer": 5,
			"modelMatched":    false,
			"setTrimName":     true,
		}})

	for _, car := range newCars {
		model := api.CreateModelInput{
			Name:   car.Model,
			Brand:  car.BrandId,
			Status: "PUBLISHED",
		}

		_, status := api.CreateModel(ctx, model)

		if status == "failed" {
			log.Println("Failed to create model: " + car.Model)
		}

		time.Sleep(500 * time.Millisecond)

	}

	Model(ctx, 5)
	log.Println("Finished Create Model Script")
}

func fixModelAndTrimName(model string, brand string, trim string, year int) (string, string) {
	var modelWords []string = strings.Split(model, " ")

	if modelWords[0] == "4" && brand == "Toyota" && len(modelWords) > 2 {
		return "4Runner", strings.Join(modelWords[2:], " ")
	}

	if modelWords[0] == "F" && brand == "Ford" && len(modelWords) > 2 {
		return modelWords[0] + modelWords[1], strings.Join(modelWords[2:], " ")
	}

	switch len(modelWords) {
	case 1:
		return model, "Base " + strconv.Itoa(year)
	case 2:
		return modelWords[0], modelWords[1]
	default:
		return modelWords[0], strings.Join(modelWords[1:], " ")
	}
}
