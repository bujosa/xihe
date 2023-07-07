package database

import (
	"context"
	"log"

	"github.com/bujosa/xihe/env"
	"github.com/bujosa/xihe/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Car struct {
	Id    string `bson:"_id"`
	Year  int    `bson:"year"`
	Title string `bson:"title"`
	Trim string `bson:"trim"`
	InteriorColor string `bson:"interiorColor"`
	ExteriorColor string `bson:"exteriorColor"`
	Mileage int `bson:"mileage"`
	LincesePlate string `bson:"licensePlate"`
	Pictures []string `bson:"pictures"`
	MainPicture string `bson:"mainPicture"`
	Currency string `bson:"currency"`
	Price int `bson:"price"`
	Url string `bson:"url"`
}

func GetCars() []Car {
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
	coll := db.Collection(utils.CARS_PROCESSED_COLLECTION)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"trimMatched": true,
				"uploaded": false,
			},
		},
		{
			"$project": bson.M{
				"year":  1,
				"title": 1,
				"trim": bson.M{
					"$toString": "$trim._id",
				},
				"interiorColor": bson.M{
					"$toString": "$interiorColor._id",
				},
				"exteriorColor": bson.M{
					"$toString": "$exteriorColor._id",
				},
				"mileage": 1,
				"licensePlate": 1,
				"pictures": 1,
				"mainPicture": 1,
				"currency": 1,
				"price": 1,
				"url": 1,
				"dealer": 1,
			},
		},
	}

	cursor, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}

	cars := []Car{}

	for cursor.Next(context.TODO()) {
		var doc Car
		err := cursor.Decode(&doc)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		cars = append(cars, doc)
	}

	return cars
}


func UpdateCar(id string, uploaded bool, status string, curboId string) {
	log.Println("Updating car: " + id)
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
	coll := db.Collection(utils.CARS_PROCESSED_COLLECTION)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"uploaded": uploaded, "status": status, "curboId": curboId}}

	_, err = coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Print("Error updating car: " + id)
	}
	log.Println("Updated car: " + id)
}