package transformation

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const BRAND_SOURCE = "brands"

func Brand(ctx context.Context) {
	log.Print("Starting brand transformation...\n")

	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"brandSlug": bson.M{
					"$toLower": "$brand",
				},
				"modelSlug": bson.M{
					"$toLower": "$model",
				},
				"interiorColor": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{
							"$split": []interface{}{
								"$interiorColor",
								"/",
							},
						},
						0,
					},
				},
				"exteriorColor": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{
							"$split": []interface{}{
								"$exteriorColor",
								"/",
							},
						},
						0,
					},
				},
				"licensePlate": bson.M{
					"$ifNull": bson.A{
						"$licensePlate",
						"$_id",
					},
				},
				"modelMatchLayer": -1,
			},
		},
		{
			"$addFields": bson.M{
				"interiorColor": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{
							"$split": []interface{}{
								"$interiorColor",
								" ",
							},
						},
						0,
					},
				},
				"exteriorColor": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{
							"$split": []interface{}{
								"$exteriorColor",
								" ",
							},
						},
						0,
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         BRAND_SOURCE,
				"localField":   "brandSlug",
				"foreignField": "slug",
				"as":           "brandObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$brandObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"brandObject": bson.M{
					"$exists": true,
					"$ne":     nil,
				},
			},
		},
		{
			"$merge": bson.M{
				"into":           utils.CARS_PROCESSED_COLLECTION,
				"on":             "_id",
				"whenMatched":    "keepExisting",
				"whenNotMatched": "insert",
			},
		},
	}
	BaseTransformation(ctx, pipeline, utils.CARS_NON_PROCESSED_COLLECTION, utils.DATABASE)
}

// This function if for clean the model lower field and eliminate some words that are not necessary
func BrandToModel(ctx context.Context, regex string, find string, replacement string) {
	log.Print("Starting brand to model transformation... with regex: " + regex + " find: " + find + " replacement: " + replacement + "\n")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"title": bson.M{
					"$regex": regex,
				},
			},
		},
		{
			"$addFields": bson.M{
				"modelSlug": bson.M{
					"$replaceOne": bson.M{
						"input":       "$modelSlug",
						"find":        find,
						"replacement": replacement,
					},
				},
			},
		},
		{
			"$merge": bson.M{
				"into":           utils.CARS_PROCESSED_COLLECTION,
				"on":             "_id",
				"whenMatched":    "replace",
				"whenNotMatched": "fail",
			},
		},
	}

	BaseTransformation(ctx, pipeline, utils.CARS_PROCESSED_COLLECTION, utils.DATABASE)
}
