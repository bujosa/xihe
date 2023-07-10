package database

import (
	"context"
	"log"
	"time"

	"github.com/bujosa/xihe/env"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dealer struct {
	Id                       string  `bson:"_id"`
	Name                     string  `bson:"name"`
	Slug                     string  `bson:"slug"`
	Address                  string  `bson:"address"`
	Latitude                 float64 `bson:"latitude"`
	Longitude                float64 `bson:"longitude"`
	TelephoneNumberSanitized string  `bson:"telephoneNumberSanitized"`
	City                     string  `bson:"city"`
	Spot                     string  `bson:"spot"`
	Uploaded                 bool    `bson:"uploaded"`
	Dealer                   string  `bson:"dealer"`
}

type UpdateDealerInfo struct {
	Id  string `bson:"_id"`
	Set bson.M `bson:"$set"`
}

func GetDealers() []Dealer {
	dbUri, err := env.GetString(utils.DB_URL_ENV_KEY)
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(
		dbUri,
	))
	if err != nil {
		panic(err)
	}

	db := client.Database(utils.DATABASE)
	coll := db.Collection(utils.DEALERS_COLLECTION)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"uploaded": false,
			},
		},
		{
			"$project": bson.M{
				"dealer":                   1,
				"uploaded":                 1,
				"name":                     1,
				"slug":                     1,
				"address":                  1,
				"latitude":                 1,
				"longitude":                1,
				"telephoneNumberSanitized": 1,
				"city":                     1,
				"spot":                     1,
			},
		},
	}

	cursor, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}

	dealers := []Dealer{}

	for cursor.Next(context.TODO()) {
		var doc Dealer
		err := cursor.Decode(&doc)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		dealers = append(dealers, doc)
	}

	return dealers
}

func UpdateDealer(updateDealerInfo UpdateDealerInfo) {
	log.Println("Updating dealer: " + updateDealerInfo.Id)

	dbUri, err := env.GetString(utils.DB_URL_ENV_KEY)
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(
		dbUri,
	))
	if err != nil {
		panic(err)
	}

	db := client.Database(utils.DATABASE)
	coll := db.Collection(utils.DEALERS_COLLECTION)

	objectId, err := utils.ToObjectId(updateDealerInfo.Id)
	if err != nil {
		log.Println("Error converting id to ObjectId: " + updateDealerInfo.Id)
		return
	}

	filter := bson.M{"_id": objectId}

	update := bson.M{
		"$set": updateDealerInfo.Set,
	}

	update["$set"].(bson.M)["updatedAt"] = time.Now().UTC()

	_, err = coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Print("Error updating dealer: " + updateDealerInfo.Id)
	}

	log.Println("Updated dealer: " + updateDealerInfo.Id)
}
