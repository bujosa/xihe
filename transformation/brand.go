package transformation

import (
	"context"
	"fmt"

	"github.com/bujosa/xihe/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DATABASE = "supercarros"
const COLLECTION = "cars"
const SOURCE = "brands"
const DB_URL_ENV_KEY = "SUPER_CARROS_DATABASE_URL"

func Brand() {
	dbUri, err := env.GetString(DB_URL_ENV_KEY)
	print(dbUri)
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
	coll := db.Collection(COLLECTION)

	// Define a filter that matches all documents in the collection.
	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"brandSlug": bson.M{
					"$toLower": "$brand",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from": SOURCE,
				"localField": "brandSlug",
				"foreignField": "slug",
				"as": "result",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$result",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"brandId": bson.M{
					"$ifNull": []interface{}{"$result._id", nil},
				},
			},
		},
		{
			"$match": bson.M{
				"brandId": bson.M{
					"$exists": true,
					"$ne": nil,
				},
			},
		},
		{
			"$out": "results",
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
		fmt.Println("Operation successful.")
	}
}