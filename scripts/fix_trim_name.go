package scripts

import (
	"context"
	"log"
	"strings"

	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func FixTrimNameForModelMatchLayer(ctx context.Context, layer int) {
	log.Print("Starting fix trim name for model match layer " + string(rune(layer)))
	cars := database.GetCarsForModelMatchLayer(ctx, 2)

	for _, car := range cars {

		// car.Model eliminar primera palabra y añadir a trimName la restante
		words := strings.Split(car.ModelSlug, "-")
		trimName := strings.Join(words[layer:], " ")
		trimName = utils.Title(trimName)

		updateCarInfo := database.UpdateCarInfo{
			Car: car,
			Set: bson.M{
				"trimName":    trimName,
				"setTrimName": true,
			},
		}

		database.UpdateCar(ctx, updateCarInfo)
	}
}
