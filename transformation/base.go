package transformation

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BaseTransformation(ctx context.Context, pipeline []bson.M, collection string, database string) {
	dbUri := ctx.Value(utils.DbUri).(string)

	// Set up the client and connect to the database.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		dbUri,
	))
	if err != nil {
		log.Println("Error connecting to database: " + dbUri)
		panic(err)
	}

	// Get the database and collection.
	db := client.Database(database)
	coll := db.Collection(collection)

	// Execute the aggregation.
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Error executing aggregation: " + err.Error())
		panic(err)
	}

	for cursor.Next(ctx) {
		var doc bson.D
		err := cursor.Decode(&doc)
		if err != nil {
			log.Println("Error decoding document: " + err.Error())
			panic(err)
		}
	}
}
