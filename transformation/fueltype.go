package transformation

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func Fueltype(ctx context.Context) {
	log.Println("Starting Mapping Fueltype transformation...")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"fueltypeMatched": false,
			},
		},
		{
			"$addFields": bson.M{
				"fueltypeSlug": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$fuelType", "Gasolina"}},
						"gasolina",
						bson.M{
							"$cond": bson.A{
								bson.M{"$eq": bson.A{"$fuelType", "Diesel"}},
								"diesel",
								bson.M{
									"$cond": bson.A{
										bson.M{"$eq": bson.A{"$fuelType", "Híbrido"}},
										"hibrido",
										bson.M{
											"$cond": bson.A{
												bson.M{"$eq": bson.A{"$fuelType", "GLP"}},
												"glp",
												bson.M{
													"$cond": bson.A{
														bson.M{"$eq": bson.A{"$fuelType", "Gas Natural"}},
														"gnc",
														"gasolina",
													},
												},
											},
										},
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
				"from":         utils.FUELTYPE_SOURCE_COLLECTION,
				"localField":   "fueltypeSlug",
				"foreignField": "slug",
				"as":           "fueltypeObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$fueltypeObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"fueltypeObject": bson.M{
					"$exists": true,
				},
			},
		},
		{
			"$addFields": bson.M{
				"fueltypeMatched": true,
			},
		},
		{
			"$project": bson.M{
				"fueltypeObject":  1,
				"fueltypeSlug":    1,
				"fueltypeMatched": 1,
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
