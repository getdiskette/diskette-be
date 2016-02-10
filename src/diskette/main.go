package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"diskette/collections"
	"diskette/middleware"
	"diskette/restservice"
	"diskette/userservice"

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

	userService := userservice.NewUserService(userCollection, jwtKey)
	public := e.Group("/public")
	public.Post("/signup", userService.Signup)
	public.Post("/confirm", userService.ConfirmSignup)
	public.Post("/signin", userService.Signin)
	public.Post("/forgot-password", userService.ForgotPassword)
	public.Post("/reset-password", userService.ResetPassword)

	sessionMiddleware := middleware.CreateSessionMiddleware(userCollection, jwtKey)
	private := e.Group("/private", sessionMiddleware)
	private.Post("/signout", userService.Signout)
	// private.Post("/change-password", userService.ChangePassword)
	// private.Post("/update-profile", userService.UpdateProfile)

	restService := restservice.NewRestService(db)
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
