package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/jonfriesen/subscriber-tracker-api/model"
	"github.com/jonfriesen/subscriber-tracker-api/storage/postgresql"
)

var (
	DatabaseURL = "postgresql://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
)

func main() {
	DatabaseURL = os.Getenv("DATABASE_URL")

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, os.Kill)

		sig := <-sigquit
		log.Printf("caught sig: %+v", sig)
		log.Printf("Gracefully shutting down server...")

		wg.Done()
	}()

	go func() {
		var adb *postgresql.PostgreSQL
		err := errors.New("database is not ready yet")
		for err != nil {
			log.Println("Checking if database is ready yet.")
			adb, err = postgresql.NewConnection(DatabaseURL)
			if err != nil {
				time.Sleep(5 * time.Second)
			}
		}

		for {

			_, err := adb.AddSubscriber(model.Subscriber{
				Name:  randomdata.FullName(randomdata.RandomGender),
				Email: randomdata.Email(),
			})
			if err != nil {
				log.Printf("Error occurred, ignoring: %s", err.Error())
			}
			time.Sleep(2 * time.Second)
		}
	}()

	wg.Wait()
}
