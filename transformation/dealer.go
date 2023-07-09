package transformation

import (
	"log"

	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/env"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const DEALER_COLLECTION = "dealers"
const DEALER_SOURCE = "dealerss"

func InitDealerTransformation() {
	utils.SetLogFile("log.txt")
	log.Println("Init dealer transformation...")

	dealers := database.GetDealers()

	print(len(dealers))

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
				"telephoneNumberSanitized": bson.M{
					"$cond": []interface{}{
						bson.M{"$ifNull": bson.A{"$telephoneNumberSanitized", nil}},
						"$telephoneNumberSanitized",
						bson.M{"$reduce": bson.M{
							"input": bson.M{
								"$map": bson.M{
									"input": bson.M{
										"$regexFindAll": bson.M{
											"input": "$telephoneNumber",
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
				"from":         DEALER_SOURCE,
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
