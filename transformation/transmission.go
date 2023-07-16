package transformation

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func Transmission(ctx context.Context) {
	log.Println("Starting Mapping Transmission transformation...")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"transmissionMatched": false,
			},
		},
		{
			"$addFields": bson.M{
				"transmissionSlug": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$transmission", "Mecánica"}},
						"manual",
						bson.M{
							"$cond": bson.A{
								bson.M{"$eq": bson.A{"$transmission", "Mecánica o Automátic"}},
								"hibrida",
								bson.M{
									"$cond": bson.A{
										bson.M{"$eq": bson.A{"$transmission", "Automática"}},
										"automatica",
										"default",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         utils.TRANSMISSION_SOURCE_COLLECTION,
				"localField":   "transmissionSlug",
				"foreignField": "slug",
				"as":           "transmissionObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$transmissionObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"transmissionObject": bson.M{"$ne": nil},
			},
		},
		{
			"$addFields": bson.M{
				"transmissionMatched": true,
			},
		},
		{
			"$merge": bson.M{
				"into":           utils.CARS_PROCESSED_COLLECTION,
				"on":             "_id",
				"whenMatched":    "merge",
				"whenNotMatched": "fail",
			},
		},
	}
	BaseTransformation(ctx, pipeline, utils.CARS_PROCESSED_COLLECTION, utils.DATABASE)
}
