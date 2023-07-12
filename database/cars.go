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

type Car struct {
	Id               string   `bson:"_id"`
	Year             int      `bson:"year"`
	Title            string   `bson:"title"`
	Trim             string   `bson:"trim"`
	InteriorColor    string   `bson:"interiorColor"`
	ExteriorColor    string   `bson:"exteriorColor"`
	Mileage          int      `bson:"mileage"`
	LicensePlate     string   `bson:"licensePlate"`
	Pictures         []string `bson:"pictures"`
	MainPicture      string   `bson:"mainPicture"`
	Currency         string   `bson:"currency"`
	Price            int      `bson:"price"`
	Url              string   `bson:"url"`
	DealerId         string   `bson:"dealerId"`
	PicturesUploaded bool     `bson:"picturesUploaded"`
	Spot             string   `bson:"spot"`
	ExteriorPictures []string `bson:"exteriorPictures"`
	InteriorPictures []string `bson:"interiorPictures"`
}

type UpdateCarInfo struct {
	Car Car    `bson:"car"`
	Set bson.M `bson:"$set"`
}

type Filter struct {
	Match bson.M `bson:"$match"`
}

func GetCars() []Car {
	filter := Filter{
		Match: bson.M{
			"trimMatched": true,
			"uploaded":    false,
			"success": bson.M{
				"$ne": "failed",
			},
		},
	}

	return BaseGetCars(filter)
}

func UpdateCar(updateCarInfo UpdateCarInfo) {
	log.Println("Updating car: " + updateCarInfo.Car.Id)

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

	objectId, err := utils.ToObjectId(updateCarInfo.Car.Id)
	if err != nil {
		log.Println("Error converting id to ObjectId: " + updateCarInfo.Car.Id)
		return
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": updateCarInfo.Set,
	}

	update["$set"].(bson.M)["updatedAt"] = time.Now().UTC()

	_, err = coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		maxRetries := 3
		retryInterval := time.Second * 3

		for i := 0; i < maxRetries; i++ {
			_, err = coll.UpdateOne(context.Background(), filter, update)
			if err == nil {
				log.Println("Updated car: " + updateCarInfo.Car.Id)
				return
			}

			log.Printf("Error updating car: %s. Retrying in %s...", updateCarInfo.Car.Id, retryInterval)
			log.Printf("Error: %s with %s", err, update)
			time.Sleep(retryInterval)
		}

		log.Println("Failed to update car after multiple attempts")
		return
	}
	log.Println("Updated car: " + updateCarInfo.Car.Id)
}

func GetCarsWithPublishDealers() []Car {

	filter := Filter{
		Match: bson.M{
			"trimMatched":         true,
			"uploaded":            false,
			"dealerObject.status": "PUBLISHED",
			"status":              bson.M{"$ne": "failed"},
		},
	}

	return BaseGetCars(filter)
}

func BaseGetCars(filter Filter) []Car {
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
			"$match": filter.Match,
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
				"mileage":          1,
				"licensePlate":     1,
				"pictures":         1,
				"mainPicture":      1,
				"currency":         1,
				"price":            1,
				"url":              1,
				"dealerId":         1,
				"picturesUploaded": 1,
				"spot":             1,
				"exteriorPictures": 1,
				"interiorPictures": 1,
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
