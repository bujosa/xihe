package scripts

import (
	"log"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/storage"
)

func UploadPictures(car database.Car, createCatInput *api.CreateCarInput) {
	log.Print("Starting upload pictures script... for car: " + car.Id)

	pictures := car.Pictures
	storage := storage.Storage{}
	mainPicture, err := storage.Upload(car.MainPicture)

	if err != nil {
		panic(err)
	}

	createCatInput.MainPicture = mainPicture
	for _, picture := range pictures {
		newPicture, err := storage.Upload(picture)

		if err != nil {
			log.Println("Error uploading picture: " + picture)
			continue
		}

		createCatInput.ExteriorPictures = append(createCatInput.ExteriorPictures, newPicture)
		createCatInput.InteriorPictures = append(createCatInput.InteriorPictures, newPicture)
	}
}