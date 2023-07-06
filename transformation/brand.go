package transformation

import (
	"context"

	"github.com/bujosa/xihe/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DATABASE = "supercarros"
const COLLECTION = "cars"
const BRAND_SOURCE = "brands"
const DB_URL_ENV_KEY = "SUPER_CARROS_DATABASE_URL"

func Brand() {
	print("Starting brand transformation...\n")

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
	coll := db.Collection(COLLECTION)

	// Define a filter that matches all documents in the collection.
	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"brand": bson.M{
					"$toLower": "$brand",
				},
				"model": bson.M{
					"$toLower": "$model",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from": BRAND_SOURCE,
				"localField": "brand",
				"foreignField": "slug",
				"as": "brand",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$brand",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"brand": bson.M{
					"$ifNull": []interface{}{"$brand", nil},
				},
			},
		},
		{
			"$match": bson.M{
				"brand": bson.M{
					"$exists": true,
					"$ne": nil,
				},
			},
		},
		{
			"$out": PROCESSED_DATA,
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

func BrandToModel(regex string, find string, replacement string) {
	print("Starting brand to model transformation... with regex: " + regex + " find: " + find + " replacement: " + replacement + "\n")


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

	// Define The pipeline.
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"title": bson.M{
					"$regex": regex,
				},
			},
		},
		{
			"$addFields": bson.M{
				"model": bson.M{
					"$replaceOne": bson.M{
									"input": "$model",
									"find": find,
									"replacement": replacement,
					},
				},
			},
		},
		{
			"$merge": bson.M{
				"into": PROCESSED_DATA,
				"on": "_id",
				"whenMatched": "replace",
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