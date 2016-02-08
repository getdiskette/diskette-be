package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"diskette/rest"

	"github.com/labstack/echo"
	"labix.org/v2/mgo"
)

type Config struct {
	Database string `json:"database"`
	JwtKey   string `json:"jwtKey"`
}

func main() {
	config := readConfig()

	session := createMongoSession()
	defer session.Close()

	db := session.DB(config.Database)

	rest := rest.NewRest(db)

	e := echo.New()
	// e.Use(func(c *echo.Context) error {
	// 	return nil
	// })
	e.Get("/:collection", rest.Get)
	e.Post("/:collection", rest.Post)
	e.Put("/:collection", rest.Put)
	e.Delete("/:collection", rest.Delete)
	e.Run(":5025")
}

func readConfig() Config {
	var config Config
	configData, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func createMongoSession() *mgo.Session {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	return session
}
