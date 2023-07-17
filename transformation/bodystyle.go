package transformation

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func BodyStyle(ctx context.Context) {
	log.Println("Starting Mapping Bodystyle transformation...")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"bodyStyleMatched": false,
			},
		},
		{
			"$addFields": bson.M{
				"bodyStyleSlug": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$bodyStyle", "Volteo"}},
						"cargo-van",
						bson.M{
							"$cond": bson.A{
								bson.M{"$eq": bson.A{"$bodyStyle", "Furgoneta"}},
								"van",
								bson.M{
									"$cond": bson.A{
										bson.M{"$eq": bson.A{"$bodyStyle", "Camión"}},
										"truck",
										bson.M{
											"$cond": bson.A{
												bson.M{"$eq": bson.A{"$bodyStyle", "MiniVán"}},
												"minivan",
												bson.M{
													"$cond": bson.A{
														bson.M{"$eq": bson.A{"$bodyStyle", "Coupé/Deportivo"}},
														"coupe",
														bson.M{
															"$cond": bson.A{
																bson.M{"$eq": bson.A{"$bodyStyle", "Coupé"}},
																"coupe",
																bson.M{
																	"$cond": bson.A{
																		bson.M{"$eq": bson.A{"$bodyStyle", "Sedán"}},
																		"sedan",
																		bson.M{
																			"$cond": bson.A{
																				bson.M{"$eq": bson.A{"$bodyStyle", "Hatchback"}},
																				"hatchback",
																				bson.M{
																					"$cond": bson.A{
																						bson.M{"$eq": bson.A{"$bodyStyle", "Jeepeta"}},
																						"muvsuv",
																						bson.M{
																							"$cond": bson.A{
																								bson.M{"$eq": bson.A{"$bodyStyle", "Jeep"}},
																								"jeep",
																								bson.M{
																									"$cond": bson.A{
																										bson.M{"$eq": bson.A{"$bodyStyle", "Convertible"}},
																										"convertible",
																										bson.M{
																											"$cond": bson.A{
																												bson.M{"$eq": bson.A{"$bodyStyle", "StationWagon"}},
																												"wagon",
																												bson.M{
																													"$cond": bson.A{
																														bson.M{"$eq": bson.A{"$bodyStyle", "Four Wheel"}},
																														"todo-terrero-atv",
																														bson.M{
																															"$cond": bson.A{
																																bson.M{"$eq": bson.A{"$bodyStyle", "Camioneta"}},
																																"pick-up",
																																"other",
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
			},
		},
		{
			"$lookup": bson.M{
				"from":         utils.BODYSTYLE_SOURCE_COLLECTION,
				"localField":   "bodyStyleSlug",
				"foreignField": "slug",
				"as":           "bodyStyleObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$bodyStyleObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"bodyStyleObject": bson.M{
					"$exists": true,
				},
			},
		},
		{
			"$addFields": bson.M{
				"bodyStyleMatched": true,
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
