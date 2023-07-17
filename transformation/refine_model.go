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

func RefineModel(ctx context.Context) {
	log.Println("Starting Create Model Script... This purifies the database from cars that are not in the model list")

	var cars []database.Car = database.BaseGetCars(ctx, database.Filter{
		Match: bson.M{
			"modelMatchLayer": -1,
			"modelMatched":    false,
		}})

	for _, car := range cars {
		modelName, trimName := transformModelName(utils.EliminateExtraSpace(car.Model), car.Brand, car.Trim, car.Year)

		updateCarInfo := database.UpdateCarInfo{
			Car: car,
			Set: bson.M{
				"model":       modelName,
				"trimName":    trimName,
				"modelSlug":   utils.Slug([]string{car.Brand, modelName}),
				"setTrimName": true,
			},
		}

		database.UpdateCar(ctx, updateCarInfo)
	}

	Model(ctx)

	var newCars []database.Car = database.BaseGetCars(ctx, database.Filter{
		Match: bson.M{
			"modelMatchLayer": -1,
			"modelMatched":    false,
		}})

	for _, car := range newCars {
		createModelInput := api.CreateModelInput{
			Name:   car.Model,
			Brand:  car.BrandId,
			Status: "PUBLISHED",
		}

		_, status := api.CreateModel(ctx, createModelInput)

		time.Sleep(500 * time.Millisecond)

		if status == "failed" {
			log.Printf("Error creating model %s\n", car.Model)
			continue
		}
	}

	Model(ctx)

	log.Println("Create Model Script finished")
}

func transformModelName(model string, brand string, trim string, year int) (string, string) {
	var modelWords []string = strings.Split(model, " ")

	if isBrandValid(brand) && len(modelWords) > 2 {
		return specialTransformation(brand, model, year)
	}

	if brand == "Kia" && len(modelWords) == 1 {
		return strings.Replace(model, "-", "", 1), "Base " + strconv.Itoa(year)
	}

	switch len(modelWords) {
	case 1:
		return model, "Base " + strconv.Itoa(year)
	case 2:
		return modelWords[0], modelWords[1]
	case 3:
		return modelWords[0], strings.Join(modelWords[1:], " ")
	case 4:
		return modelWords[0] + " " + modelWords[1], strings.Join(modelWords[2:], " ")
	default:
		return modelWords[0], strings.Join(modelWords[1:], " ")
	}
}

func specialTransformation(brand string, model string, year int) (string, string) {
	var modelWords []string = strings.Split(model, " ")

	if brand == "Lexus" {
		return modelWords[0], strings.Join(modelWords[1:], " ")
	}

	if brand == "Chevrolet" {
		if modelWords[0] == "Trail" {
			return "TrailBlazer", strings.Join(modelWords[2:], " ")
		}
	}

	switch modelWords[0] {
	case "Clase":
		return "Clase " + modelWords[1], strings.Join(modelWords[2:], " ")
	case "X":
		return "X" + modelWords[1], strings.Join(modelWords[2:], " ")
	case "4":
		return "4" + modelWords[1], strings.Join(modelWords[2:], " ")
	case "CS":
		return "CS" + modelWords[1], strings.Join(modelWords[2:], " ")
	default:
		return modelWords[0], strings.Join(modelWords[1:], " ")
	}
}

var validBrands = []string{"BMW", "Mercedes-Benz", "Toyota", "Changan", "Lexus", "Chevrolet"}

func isBrandValid(brand string) bool {
	for _, b := range validBrands {
		if b == brand {
			return true
		}
	}
	return false
}
