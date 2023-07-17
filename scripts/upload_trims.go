package scripts

import (
	"context"
	"log"

	"github.com/bujosa/xihe/database"
	"go.mongodb.org/mongo-driver/bson"
)

func UploadTrims(ctx context.Context) {
	log.Println("Starting Upload Trims...")

	features := database.GetFeatures(ctx)

	cars := database.GetCars(ctx, database.Filter{
		Match: bson.M{
			"setTrimName": true,
			"trimMatched": false,
		}})

	for _, car := range cars {
		
}
