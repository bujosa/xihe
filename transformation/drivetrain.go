package transformation

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func DriveTrain(ctx context.Context) {
	log.Println("Starting Mapping DriveTrain transformation...")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"driveTrainMatched": false,
			},
		},
		{
			"$addFields": bson.M{
				"driveTrainSlug": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$driveTrain", "4WD"}},
						"traccion-trasera-y-delantera-4wd",
						bson.M{
							"$cond": bson.A{
								bson.M{"$eq": bson.A{"$driveTrain", "Trasera"}},
								"traccion-trasera-rwd",
								bson.M{
									"$cond": bson.A{
										bson.M{"$eq": bson.A{"$driveTrain", "Delantera"}},
										"traccion-delantera-fwd",
										bson.M{
											"$cond": bson.A{
												bson.M{"$eq": bson.A{"$driveTrain", "2WD"}},
												"traccion-en-dos-ruedas-2wd",
												bson.M{
													"$cond": bson.A{
														bson.M{"$eq": bson.A{"$driveTrain", "AWD"}},
														"traccion-en-todas-las-ruedas-awd",
														bson.M{
															"$cond": bson.A{
																bson.M{"$eq": bson.A{"$driveTrain", "6x4"}},
																"6x4",
																"full",
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
			},
		},
		{
			"$lookup": bson.M{
				"from":         utils.BODYSTYLE_SOURCE_COLLECTION,
				"localField":   "driveTrainSlug",
				"foreignField": "slug",
				"as":           "driveTrainObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$driveTrainObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"driveTrainObject": bson.M{
					"$ne": nil,
				},
			},
		},
		{
			"$addFields": bson.M{
				"driveTrainMatched": true,
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
