package transformation

import (
	"context"

	"github.com/bujosa/xihe/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BaseTransformation(pipeline []bson.M, collection string, database string) {
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
	db := client.Database(database)
	coll := db.Collection(collection)

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