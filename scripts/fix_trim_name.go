package scripts

import (
	"context"
	"log"
	"strings"

	"github.com/bujosa/xihe/database"
	"go.mongodb.org/mongo-driver/bson"
)

func FixTrimNameForModelMatchLayerTwo(ctx context.Context) {
	log.Print("Starting fix trim name for model match layer two... ")
	cars := database.GetCarsForModelMatchLayerTwo(ctx)

	for _, car := range cars {

		// car.Model eliminar primera palabra y añadir a trimName la restante
		words := strings.Split(car.ModelSlug, "-")
		trimName := strings.Join(words[2:], " ")
		trimName = strings.Title(trimName)

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
