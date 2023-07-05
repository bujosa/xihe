package main

import (
	"github.com/bujosa/xihe/transformation"
)


func main() {

	// Execute the aggregation pipeline.
	// transformation.Brand()
    transformation.Model()

    // // Set up the client and connect to the database.
    // client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(
	// 	"mongodb://localhost:27017",
	// ))
    // if err != nil {
    //     panic(err)
    // }

    // // Get the database and collection.
    // db := client.Database("supercarros")
    // coll := db.Collection("cars")

    // // Define a filter that matches all documents in the collection.
    //    pipeline := []bson.M{
    //     {
    //         "$match": bson.M{
    //             "year": 2023,
    //         },
    //     },
    //     {
    //         "$limit": 10,
    //     },
	// 	{
	// 		"$project": bson.M{
	// 			"year": 1,
	// 			"title": 1,
	// 		},
	// 	},
    // }
    // // Find all documents in the collection.
    // cursor, err := coll.Aggregate(context.TODO(), pipeline)
    // if err != nil {
    //     panic(err)
    // }

	// cars := []Car{}

    // // Iterate through the cursor and print the documents.
    // for cursor.Next(context.TODO()) {
    //     var doc Car
    //     err := cursor.Decode(&doc)
    //     if err != nil {
    //         panic(err)
    //     }
	// 	cars = append(cars, doc)
    // }

	// // Iterate Cars
	// for _, car := range cars {
	// 	fmt.Printf("Year: %d\nTitle: %s\nId: %s\n\n", car.Year, car.Title, car.Id)
	// }
}