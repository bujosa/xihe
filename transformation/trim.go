package transformation

import (
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const TRIM_SOURCE = "trimlevels"

func Trim() {
	log.Print("Starting trim transformation... \n")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"modelMatched": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": TRIM_SOURCE,
				"let": bson.M{
					"modelId": "$model._id",
					"year":    "$year",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{
										"$eq": []interface{}{
											"$carModel",
											"$$modelId",
										},
									},
									{
										"$eq": []interface{}{
											"$year",
											"$$year",
										},
									},
								},
							},
						},
					},
					{
						"$match": bson.M{
							"deleted": false,
						},
					},
					{
						"$limit": 1,
					},
					{
						"$project": bson.M{
							"_id":          1,
							"id":           1,
							"slug":         1,
							"name":         1,
							"carModel":     1,
							"bodyStyle":    1,
							"driveTrain":   1,
							"fuelType":     1,
							"transmission": 1,
							"features":     1,
						},
					},
				},
				"as": "trim",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$trim",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"trimMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$ifNull": []interface{}{
								"$trim",
								false,
							},
						},
						"then": true,
						"else": false,
					},
				},
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

	BaseTransformation(pipeline, utils.CARS_PROCESSED_COLLECTION, utils.DATABASE)
}
