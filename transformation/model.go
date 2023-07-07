package transformation

import (
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const MODEL_SOURCE = "models"

func Model() {
	print("Starting model transformation... \n")

	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"modelSlug": bson.M{
					"$concat": []interface{}{
						"$brand.slug",
						"-",
						bson.M{
							"$replaceAll": bson.M{
									"input": "$model",
									"find": " ",
									"replacement": "-",
							},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from": MODEL_SOURCE,
				"localField": "modelSlug",
				"foreignField": "slug",
				"as": "model",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$model",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"modelMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$ifNull": []interface{}{
								"$model",
								false,
							},
						},
						"then": true,
						"else": false,
					},
				},
				"trimMatched": false,
				"uploaded": false,
			},
		},
		{"$out": utils.CARS_PROCESSED},
	}

	BaseTransformation(pipeline, utils.CARS_PROCESSED, utils.DATABASE)
}