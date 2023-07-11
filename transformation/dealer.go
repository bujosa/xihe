package transformation

import (
	"log"

	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/env"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func InitDealerTransformation() {
	log.Println("Init dealer transformation...")

	dealers := database.GetDealers()

	for _, dealer := range dealers {
		slug := utils.Slug([]string{dealer.Name})
		geoCode, err := env.GetString("GEOCODE")
		if err != nil {
			panic(err)
		}

		lat, lng, _ := utils.Geocode(geoCode + " " + dealer.Address)

		updateDealerInfo := database.UpdateDealerInfo{
			Id: dealer.Id,
			Set: bson.M{
				"slug":      slug,
				"latitude":  lat,
				"longitude": lng,
			},
		}

		database.UpdateDealer(updateDealerInfo)
	}

	log.Println("Dealer transformation finished!")
}

func DealerTransformation() {
	log.Println("Starting Dealer transformation...")

	city, err := env.GetString("CITY_ID")
	if err != nil {
		panic(err)
	}

	spot, err := env.GetString("SPOT_ID")
	if err != nil {
		panic(err)
	}

	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"uploaded": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$uploaded", false}},
						"$uploaded",
						false,
					},
				},
				"latitude": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$latitude", nil}},
						"$latitude",
						0,
					},
				},
				"longitude": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$longitude", nil}},
						"$longitude",
						0,
					},
				},
				"targetNumber": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$telephoneNumber", nil}},
						"$telephoneNumber",
						"$whatssapp",
					},
				},
				"telephoneNumberSanitized": bson.M{
					"$cond": []interface{}{
						bson.M{"$ifNull": bson.A{"$telephoneNumberSanitized", nil}},
						"$telephoneNumberSanitized",
						bson.M{"$reduce": bson.M{
							"input": bson.M{
								"$map": bson.M{
									"input": bson.M{
										"$regexFindAll": bson.M{
											"input": "$targetNumber",
											"regex": "\\d",
										},
									},
									"as": "char",
									"in": "$$char.match",
								},
							},
							"initialValue": "",
							"in": bson.M{
								"$concat": []interface{}{"$$value", "$$this"},
							},
						}},
					},
				},
				"slug": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$slug", nil}},
						"$slug",
						nil,
					},
				},
				"city": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$city", nil}},
						"$city",
						city,
					},
				},
				"spot": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$spot", nil}},
						"$spot",
						spot,
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         utils.DEALERS_SOURCE_COLLECTION,
				"localField":   "slug",
				"foreignField": "slug",
				"as":           "dealer",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$dealer",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"dealerObject": "$dealer",
				"dealer": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$dealer", nil}},
						"$dealer.id",
						nil,
					},
				},
				"uploaded": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$dealer", nil}},
						true,
						"$uploaded",
					},
				},
			},
		},
		{
			"$out": utils.DEALERS_COLLECTION,
		},
	}
	BaseTransformation(pipeline, utils.DEALERS_COLLECTION, utils.DATABASE)
}

func DealerIntoCarTransformation() {
	log.Println("Starting Dealer into Car transformation...")

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         utils.DEALERS_COLLECTION,
				"localField":   "dealer",
				"foreignField": "name",
				"as":           "dealerObject",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$dealerObject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"spot": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$spot", nil}},
						"$spot",
						"$dealerObject.spot",
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"dealerObject": "$dealerObject.dealerObject",
				"dealerId":     "$dealerObject.dealer",
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
