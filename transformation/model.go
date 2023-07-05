package transformation

import (
	"context"

	"github.com/bujosa/xihe/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MODEL_SOURCE = "models"
const PROCESSED_DATA = "cars_processed"

func Model() {
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
			"$addFields": bson.M{
				"model": bson.M{
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
				"localField": "model",
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
			},
		},
		{"$out": PROCESSED_DATA},
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