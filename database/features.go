package database

import (
	"context"
	"log"

	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Feature struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
	Slug string `bson:"slug"`
}

func GetFeatures(ctx context.Context) []Feature {
	dbUri := ctx.Value(utils.DbUri).(string)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		dbUri,
	))
	if err != nil {
		log.Println("Error connecting to database: " + dbUri)
		return nil
	}

	db := client.Database(utils.DATABASE)
	coll := db.Collection(utils.FEATURE_SOURCE_COLLECTION)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"deleted": false,
				"status":  "PUBLISHED",
			},
		},
		{
			"$project": bson.M{
				"id":   1,
				"name": 1,
				"slug": 1,
			},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Error getting features with error " + err.Error())
		panic(err)
	}

	features := []Feature{}

	for cursor.Next(ctx) {
		var doc Feature
		err := cursor.Decode(&doc)
		if err != nil {
			log.Println("Error decoding feature with error " + err.Error())
			panic(err)
		}
		features = append(features, doc)
	}

	return features
}
