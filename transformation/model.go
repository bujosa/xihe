package transformation

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const MODEL_SOURCE = "models"

// Model function is for lookup the model in the models collection
func Model(ctx context.Context) {
	log.Print("Starting model transformation... \n")

	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"pictureUploaded": bson.M{
					"$ifNull": bson.A{
						"$pictureUploaded",
						false,
					},
				},
				"interiorPictures": bson.M{
					"$ifNull": bson.A{
						"$interiorPictures",
						bson.A{},
					},
				},
				"exteriorPictures": bson.M{
					"$ifNull": bson.A{
						"$exteriorPictures",
						bson.A{},
					},
				},
				"modelSlug": bson.M{
					"$ifNull": bson.A{
						"$modelSlug",
						bson.M{
							"$concat": []interface{}{
								"$brandObject.slug",
								"-",
								bson.M{
									"$replaceAll": bson.M{
										"input":       "$modelSlug",
										"find":        " ",
										"replacement": "-",
									},
								},
							},
						},
					},
				},
				"trimMatched": bson.M{
					"$ifNull": bson.A{
						"$trimMatched",
						false,
					},
				},
				"fueltypeMatched": bson.M{
					"$ifNull": bson.A{
						"$fueltypeMatched",
						false,
					},
				},
				"transmissionMatched": bson.M{
					"$ifNull": bson.A{
						"$transmissionMatched",
						false,
					},
				},
				"bodyStyleMatched": bson.M{
					"$ifNull": bson.A{
						"$bodyStyleMatched",
						false,
					},
				},
				"driveTrainMatched": bson.M{
					"$ifNull": bson.A{
						"$drivetrainMatched",
						false,
					},
				},
				"setTrimName": bson.M{
					"$ifNull": bson.A{
						"$setTrimName",
						false,
					},
				},
				"uploaded": bson.M{
					"$ifNull": bson.A{
						"$uploaded",
						false,
					},
				},
				"modelMatched": bson.M{
					"$ifNull": bson.A{
						"$modelMatched",
						false,
					},
				},
				"modelMatchLayer": bson.M{
					"$ifNull": bson.A{
						"$modelMatchLayer",
						-1,
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         MODEL_SOURCE,
				"localField":   "modelSlug",
				"foreignField": "slug",
				"as":           "modelObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$modelObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"modelMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$ifNull": []interface{}{
								"$modelObject",
								false,
							},
						},
						"then": true,
						"else": "$modelMatched",
					},
				},
				"modelMatchLayer": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$ifNull": []interface{}{
								"$modelObject",
								false,
							},
						},
						"then": 1,
						"else": "$modelMatchLayer",
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
}

func UnMatchedModelLayerTwo(ctx context.Context) {

	log.Println("Starting unmatched model layer one transformation...")

	// In model slug pick the first two words and search in the models collection
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"modelMatched": false,
			},
		},
		{
			"$addFields": bson.M{
				"trimName": "$modelSlug",
				"modelSlug": bson.M{
					"$concat": []interface{}{
						bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{
									"$split": []interface{}{
										"$modelSlug",
										"-",
									},
								},
								0,
							},
						},
						"-",
						bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{
									"$split": bson.A{
										"$modelSlug",
										"-",
									},
								},
								1,
							},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from": MODEL_SOURCE,
				"let":  bson.M{"modelSlug": "$modelSlug"},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": bson.A{
									bson.M{"$eq": bson.A{"$slug", "$$modelSlug"}},
									bson.M{"$eq": bson.A{"$deleted", false}},
								},
							},
						},
					},
				},
				"as": "modelObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$modelObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"modelMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$ifNull": []interface{}{
								"$modelObject",
								false,
							},
						},
						"then": true,
						"else": false,
					},
				},
				"modelMatchLayer": 2,
			},
		},
		{
			"$match": bson.M{
				"modelMatched": true,
			},
		},
		{
			"$project": bson.M{
				"modelObject":     1,
				"modelMatched":    1,
				"modelMatchLayer": 1,
				"trimName":        1,
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

func UnMatchedModelLayerThree(ctx context.Context) {
	log.Println("Starting unmatched model layer two transformation...")

	// Check first model slug split - have more or equal to 3 words
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"modelMatched": false,
			},
		},
		{
			"$addFields": bson.M{
				"trimName": "$modelSlug",
				"modelSlugCount": bson.M{
					"$size": bson.M{
						"$split": bson.A{
							"$modelSlug",
							"-",
						},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"modelSlugCount": bson.M{
					"$gte": 3,
				},
			},
		},
		{
			"$addFields": bson.M{
				"modelSlug": bson.M{
					"$concat": []interface{}{
						bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{
									"$split": bson.A{
										"$modelSlug",
										"-",
									},
								},
								0,
							},
						},
						"-",
						bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{
									"$split": bson.A{
										"$modelSlug",
										"-",
									},
								},
								1,
							},
						},
						"-",
						bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{
									"$split": bson.A{
										"$modelSlug",
										"-",
									},
								},
								2,
							},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from": MODEL_SOURCE,
				"let":  bson.M{"modelSlug": "$modelSlug"},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": bson.A{
									bson.M{"$eq": bson.A{"$slug", "$$modelSlug"}},
									bson.M{"$eq": bson.A{"$deleted", false}},
								},
							},
						},
					},
				},
				"as": "modelObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$modelObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"modelMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$ifNull": []interface{}{
								"$modelObject",
								false,
							},
						},
						"then": true,
						"else": false,
					},
				},
				"modelMatchLayer": 3,
			},
		},
		{
			"$match": bson.M{
				"modelMatched": true,
			},
		},
		{
			"$project": bson.M{
				"modelObject":     1,
				"modelMatched":    1,
				"modelMatchLayer": 1,
				"trimName":        1,
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
