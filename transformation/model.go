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
						"else": false,
					},
				},
				"trimMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$ifNull": []interface{}{
								"$trimObject",
								false,
							},
						},
						"then": "$trimMatched",
						"else": false,
					},
				},
				"uploaded": bson.M{
					"$ifNull": bson.A{
						"$uploaded",
						false,
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
