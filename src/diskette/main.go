package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"diskette/admin"
	"diskette/collections"
	"diskette/middleware"
	"diskette/rest"
	"diskette/user"

	"github.com/labstack/echo"
	"labix.org/v2/mgo"
)

type config struct {
	Database string `json:"database"`
	JwtKey   string `json:"jwtKey"`
}

func main() {
	cfg := readConfig()
	jwtKey := []byte(cfg.JwtKey)

	mongoSession := createMongoSession()
	defer mongoSession.Close()

	db := mongoSession.DB(cfg.Database)
	userCollection := db.C(collections.UserCollectionName)

	e := echo.New()

	userService := user.NewService(userCollection, jwtKey)
	user := e.Group("/user")
	user.Post("/signup", userService.Signup)
	user.Post("/confirm", userService.ConfirmSignup)
	user.Post("/signin", userService.Signin)
	user.Post("/forgot-password", userService.ForgotPassword)
	user.Post("/reset-password", userService.ResetPassword)

	sessionMiddleware := middleware.CreateSessionMiddleware(userCollection, jwtKey)
	private := e.Group("/private", sessionMiddleware)
	private.Post("/signout", userService.Signout)
	private.Post("/change-password", userService.ChangePassword)
	private.Post("/change-email", userService.ChangeEmail)
	private.Post("/update-profile", userService.UpdateProfile)

	adminService := admin.NewService(userCollection, jwtKey)
	// adminSessionMiddleware := middleware.CreateAdminSessionMiddleware(userCollection, jwtKey)
	admin := e.Group("/admin") //, adminSessionMiddleware)
	admin.Get("/get-users", adminService.GetUsers)

	restService := rest.NewService(db)
	e.Get("/:collection", restService.Get)
	e.Post("/:collection", restService.Post)
	e.Put("/:collection", restService.Put)
	e.Delete("/:collection", restService.Delete)

	fmt.Println("Listening at http://localhost:5025")
	e.Run(":5025")
}

func readConfig() config {
	var cfg config
	cfgData, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

func createMongoSession() *mgo.Session {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	return session
}
