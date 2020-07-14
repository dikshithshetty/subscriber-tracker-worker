package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/jonfriesen/subscriber-tracker-worker/model"
	"github.com/jonfriesen/subscriber-tracker-worker/storage/postgresql"
)

var (
	DatabaseURLDefault = "postgresql://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
)

func main() {
	DatabaseURL := os.Getenv("DATABASE_URL")
	if DatabaseURL == "" {
		DatabaseURL = DatabaseURLDefault
	}

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
		log.Println("Checking if database is ready yet.")
		for err != nil {
			adb, err = postgresql.NewConnection(DatabaseURL)
			if err != nil {
				log.Println("...no database found yet")
				err = err
				time.Sleep(2 * time.Second)
			} else {
				log.Println("Found DB!")
				adb = adb
				break
			}
		}

		for {
			newSub := &model.Subscriber{
				Name:  randomdata.FullName(randomdata.RandomGender),
				Email: randomdata.Email(),
			}
			_, err := adb.AddSubscriber(newSub)
			log.Printf("Added: %s <%s>", newSub.Name, newSub.Email)
			if err != nil {
				log.Printf("Error occurred, ignoring: %s", err.Error())
			}
			time.Sleep(5 * time.Second)
		}
	}()

	wg.Wait()
}
