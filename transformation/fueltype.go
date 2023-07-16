package transformation

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func MappingFueltypeTransformation(ctx context.Context) {
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
						bson.M{"$eq": bson.A{"$fueltype", "Gasolina"}},
						"gasolina",
						bson.M{
							"$cond": bson.A{
								bson.M{"$eq": bson.A{"$fueltype", "Diesel"}},
								"diesel",
								bson.M{
									"$cond": bson.A{
										bson.M{"$eq": bson.A{"$fueltype", "Híbrido"}},
										"hibrido",
										bson.M{
											"$cond": bson.A{
												bson.M{"$eq": bson.A{"$fueltype", "GLP"}},
												"glp",
												bson.M{
													"$cond": bson.A{
														bson.M{"$eq": bson.A{"$fueltype", "Gas Natural"}},
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
				"localField":   "fuelTypeSlug",
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
			"$addFields": bson.M{
				"fueltypeMatched": true,
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
