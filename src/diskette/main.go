package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"diskette/rest"
	"diskette/user"

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

	restService := rest.NewRestService(db)
	userService := user.NewUserService(db, []byte(config.JwtKey))

	e := echo.New()

	// e.Use(func(c *echo.Context) error {
	// 	return nil
	// })

	e.Get("/:collection", restService.Get)
	e.Post("/:collection", restService.Post)
	e.Put("/:collection", restService.Put)
	e.Delete("/:collection", restService.Delete)

	e.Post("/user/signup", userService.Signup)

	fmt.Println("Listening at http://localhost:5025")
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
