package transformation

import (
	"context"

	"github.com/bujosa/xihe/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TRIM_SOURCE = "trimlevels"

func Trim() {
	print("Starting trim transformation... \n")

	dbUri, err := env.GetString(DB_URL_ENV_KEY)
	if err != nil {
		panic(err)
	}

	// Set up the client and connect to the database.
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(
		dbUri,
	))
	if err != nil {
		panic(err)
	}

	// Get the database and collection.
	db := client.Database(DATABASE)
	coll := db.Collection(PROCESSED_DATA)

	// Define a filter that matches all documents in the collection.
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"modelMatched": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": TRIM_SOURCE,
				"let": bson.M{
					"modelId": "$model._id",
					"year": "$year",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{
										"$eq": []interface{}{
											"$carModel",
											"$$modelId",
										},
									},
									{
										"$eq": []interface{}{
											"$year",
											"$$year",
										},
									},
								},
							},
						},
					},
					{
						"$limit": 1,
					},
					{
						"$project": bson.M{
							"_id": 1,
							"id": 1,
							"slug": 1,
							"name": 1,
							"carModel": 1,
							"bodyStyle": 1,
							"driveTrain": 1,
							"fuelType": 1,
							"transmission": 1,
							"features": 1,
						},
					},
				},
				"as": "trim",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$trim",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"trimMatched": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$ifNull": []interface{}{
								"$trim",
								false,
							},
						},
						"then": true,
						"else": false,
					},
				},
			},
		},
		{ 
			"$merge": bson.M{
				"into": PROCESSED_DATA,
				"on": "_id",
				"whenMatched": "merge",
				"whenNotMatched": "fail",
			},
		},
	}

	// Execute the aggregation.
	cursor, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var doc bson.D
		err := cursor.Decode(&doc)
		if err != nil {
			panic(err)
		}
	}
}