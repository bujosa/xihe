package scripts

import (
	"log"
	"time"

	"github.com/bujosa/xihe/api"
	"github.com/bujosa/xihe/database"
	"github.com/bujosa/xihe/storage"
	"github.com/bujosa/xihe/utils"
)

func UploadPictures(storage *storage.Storage, car database.Car, createCarInput *api.CreateCarInput) error {
	log.Print("Starting upload pictures script... for car: " + car.Id)

	pictures := car.Pictures
	mainPicture, err := storage.Upload(car.MainPicture)

	if err != nil {
		retry := 5
		for retry > 0 {
			mainPicture, err = storage.Upload(car.MainPicture)

			if err == nil {
				break
			}

			if retry < 3 {
				result := storage.RestartConnection()
				if result != nil {
					log.Println("Error restarting connection to storage")
				}
			}

			retry--
			time.Sleep(time.Duration(retry) * time.Second)
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
			retry := 5
			for retry > 0 {
				newPicture, err = storage.Upload(picture)

				if err == nil {
					break
				}

				if retry < 3 {
					utils.CleanDns()
					result := storage.RestartConnection()
					if result != nil {
						log.Println("Error restarting connection to storage")
					}
				}

				retry--
				time.Sleep(time.Duration(retry) * time.Second)
			}

			if err != nil {
				log.Println("Error picture: " + picture)
				continue
			}
		}

		createCarInput.ExteriorPictures = append(createCarInput.ExteriorPictures, newPicture)
		createCarInput.InteriorPictures = append(createCarInput.InteriorPictures, newPicture)

		time.Sleep(time.Second * 1)
	}

	return nil
}
