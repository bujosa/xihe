package transformation

import (
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const COLORS_SOURCE = "colors"

func Color() {
	log.Print("Starting color transformation... \n")

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from": COLORS_SOURCE,
				"let": bson.M{
					"searchColor": "$interiorColor",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$or": []interface{}{
									bson.M{
										"$regexMatch": bson.M{
											"input":   "$name",
											"regex":   "$$searchColor",
											"options": "i",
										},
									},
									bson.M{
										"$regexMatch": bson.M{
											"input":   "$name",
											"regex":   "Cream",
											"options": "i",
										},
									},
								},
							},
						},
					},
					{
						"$limit": 1,
					},
				},
				"as": "interiorColor",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$interiorColor",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"interiorColorMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []interface{}{
								"$interiorColor.slug",
								"cream",
							},
						},
						"then": false,
						"else": true,
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from": COLORS_SOURCE,
				"let": bson.M{
					"searchColor": "$exteriorColor",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$or": []interface{}{
									bson.M{
										"$regexMatch": bson.M{
											"input":   "$name",
											"regex":   "$$searchColor",
											"options": "i",
										},
									},
									bson.M{
										"$regexMatch": bson.M{
											"input":   "$name",
											"regex":   "Cream",
											"options": "i",
										},
									},
								},
							},
						},
					},
					{
						"$limit": 1,
					},
				},
				"as": "exteriorColor",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$exteriorColor",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"exteriorColorMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []interface{}{
								"$exteriorColor.slug",
								"cream",
							},
						},
						"then": false,
						"else": true,
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
