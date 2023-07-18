package transformation

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const TRIM_SOURCE = "trimlevels"

func Trim(ctx context.Context) {
	log.Print("Starting Trim Transformation ... \n")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"modelMatched": true,
				"trimMatched":  false,
			},
		},
		{
			"$lookup": bson.M{
				"from": TRIM_SOURCE,
				"let": bson.M{
					"modelId":  "$modelObject._id",
					"year":     "$year",
					"trimName": "$trimName",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$or": []bson.M{
									{
										"$and": []bson.M{
											{
												"$eq": []interface{}{
													"$carModel",
													"$$modelId",
												},
											},
											{
												"$eq": []interface{}{
													"$name",
													"$$trimName",
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
									{
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
							"deleted":      1,
						},
					},
				},
				"as": "trimObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$trimObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"trimMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$ifNull": []interface{}{
								"$trimObject",
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

	BaseTransformation(ctx, pipeline, utils.CARS_PROCESSED_COLLECTION, utils.DATABASE)

	log.Print("Trim Transformation Finished \n")
}
