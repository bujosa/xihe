package scripts

import (
	"log"
	"time"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/storage"
)

func UploadPictures(storage *storage.Storage, car database.Car, createCarInput *api.CreateCarInput) error {
	log.Print("Starting upload pictures script... for car: " + car.Id)

	pictures := car.Pictures
	mainPicture, err := storage.Upload(car.MainPicture)

	if err != nil {
		retry := 0
		for retry < 5 {
			mainPicture, err = storage.Upload(car.MainPicture)

			if err == nil {
				break
			}

			retry++
			time.Sleep(time.Second * 2)
		}

		if err != nil {
			log.Println("Error uploading main picture: " + car.MainPicture)
			return err
		}
	}

	createCarInput.MainPicture = mainPicture
	for _, picture := range pictures {
		newPicture, err := storage.Upload(picture)

		if err != nil {
			log.Println("Error uploading picture: " + picture)
			continue
		}

		createCarInput.ExteriorPictures = append(createCarInput.ExteriorPictures, newPicture)
		createCarInput.InteriorPictures = append(createCarInput.InteriorPictures, newPicture)

		time.Sleep(time.Millisecond * 300)
	}

	return nil
}
