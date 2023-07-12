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
				"brand": bson.M{
					"$toLower": "$brand",
				},
				"model": bson.M{
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
				"licensePlate":     "$_id",
				"picturesUploaded": false,
				"interiorPictures": bson.M{
					"$cond": bson.A{
						bson.M{"$or": bson.A{
							bson.M{"$eq": bson.A{"$interiorPictures", nil}},
							bson.M{"$eq": bson.A{"$interiorPictures", bson.A{}}},
						}},
						bson.A{},
						"$interiorPictures",
					},
				},
				"exteriorPictures": bson.M{
					"$cond": bson.A{
						bson.M{"$or": bson.A{
							bson.M{"$eq": bson.A{"$exteriorPictures", nil}},
							bson.M{"$eq": bson.A{"$exteriorPictures", bson.A{}}},
						}},
						bson.A{},
						"$exteriorPictures",
					},
				},
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
				"localField":   "brand",
				"foreignField": "slug",
				"as":           "brand",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$brand",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"brand": bson.M{
					"$ifNull": []interface{}{"$brand", nil},
				},
			},
		},
		{
			"$match": bson.M{
				"brand": bson.M{
					"$exists": true,
					"$ne":     nil,
				},
			},
		},
		{
			"$out": utils.CARS_PROCESSED_COLLECTION,
		},
	}
	BaseTransformation(ctx, pipeline, utils.CARS_NON_PROCESSED_COLLECTION, utils.DATABASE)
}

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
				"model": bson.M{
					"$replaceOne": bson.M{
						"input":       "$model",
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
