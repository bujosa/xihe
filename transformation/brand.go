package transformation

import (
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const COLLECTION = "cars"
const BRAND_SOURCE = "brands"

func Brand() {
	print("Starting brand transformation...\n")

	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"brand": bson.M{
					"$toLower": "$brand",
				},
				"model": bson.M{
					"$toLower": "$model",
				},
				"interiorColor":  bson.M{
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
				"exteriorColor":  bson.M{
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
				"licensePlate": "$_id",
			},
		},
		{
			"$addFields": bson.M{
				"interiorColor":  bson.M{
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
				"exteriorColor":  bson.M{
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
				"from": BRAND_SOURCE,
				"localField": "brand",
				"foreignField": "slug",
				"as": "brand",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$brand",
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
					"$ne": nil,
				},
			},
		},
		{
			"$out": PROCESSED_DATA,
		},
	}
	BaseTransformation(pipeline, COLLECTION, utils.DATABASE)
}

func BrandToModel(regex string, find string, replacement string) {
	print("Starting brand to model transformation... with regex: " + regex + " find: " + find + " replacement: " + replacement + "\n")

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
									"input": "$model",
									"find": find,
									"replacement": replacement,
					},
				},
			},
		},
		{
			"$merge": bson.M{
				"into": PROCESSED_DATA,
				"on": "_id",
				"whenMatched": "replace",
				"whenNotMatched": "fail",
			},
		},
	}

	BaseTransformation(pipeline, PROCESSED_DATA, utils.DATABASE)
}