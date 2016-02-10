package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"diskette/collections"
	"diskette/middleware"
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
	jwtKey := []byte(config.JwtKey)

	mongoSession := createMongoSession()
	defer mongoSession.Close()

	db := mongoSession.DB(config.Database)
	userCollection := db.C(collections.USER_COLLECTION_NAME)

	e := echo.New()

	userService := user.NewAuthenticationService(userCollection, jwtKey)
	user := e.Group("/user")
	user.Post("/signup", userService.Signup)
	user.Post("/confirm", userService.ConfirmSignup)
	user.Post("/signin", userService.Signin)
	user.Post("/forgot-password", userService.ForgotPassword)
	user.Post("/reset-password", userService.ResetPassword)

	sessionMiddleware := middleware.CreateSessionMiddleware(userCollection, jwtKey)
	private := e.Group("/private", sessionMiddleware)
	private.Post("/signout", userService.Signout)
	// private.Post("/change-password", userService.ChangePassword)
	// private.Post("/update-profile", userService.UpdateProfile)

	restService := rest.NewRestService(db)
	e.Get("/:collection", restService.Get)
	e.Post("/:collection", restService.Post)
	e.Put("/:collection", restService.Put)
	e.Delete("/:collection", restService.Delete)

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
