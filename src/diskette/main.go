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
	"diskette/session"
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

	restService := rest.NewService(db)
	restGroup := e.Group("/collection")
	restGroup.Get("/:collection", restService.Get)
	restGroup.Post("/:collection", restService.Post)
	restGroup.Put("/:collection", restService.Put)
	restGroup.Delete("/:collection", restService.Delete)

	userService := user.NewService(userCollection, jwtKey)
	userGroup := e.Group("/user")
	userGroup.Post("/signup", userService.Signup)
	userGroup.Post("/confirm", userService.ConfirmSignup)
	userGroup.Post("/signin", userService.Signin)
	userGroup.Post("/forgot-password", userService.ForgotPassword)
	userGroup.Post("/reset-password", userService.ResetPassword)

	sessionService := session.NewService(userCollection, jwtKey)
	sessionMiddleware := middleware.CreateSessionMiddleware(userCollection, jwtKey)
	sessionGroup := e.Group("/session", sessionMiddleware)
	sessionGroup.Post("/signout", sessionService.Signout)
	sessionGroup.Post("/change-password", sessionService.ChangePassword)
	sessionGroup.Post("/change-email", sessionService.ChangeEmail)
	sessionGroup.Post("/set-profile", sessionService.SetProfile)

	adminService := admin.NewService(userCollection, jwtKey)
	adminSessionMiddleware := middleware.CreateAdminSessionMiddleware(userCollection, jwtKey)
	adminGroup := e.Group("/admin", adminSessionMiddleware)
	adminGroup.Get("/get-users", adminService.GetUsers)
	adminGroup.Post("/create-user", adminService.CreateUser)
	adminGroup.Post("/change-user-password", adminService.ChangeUserPassword)
	adminGroup.Post("/change-user-email", adminService.ChangeUserEmail)
	adminGroup.Post("/set-user-roles", adminService.SetUserRoles)
	adminGroup.Post("/set-user-profile", adminService.SetUserProfile)
	adminGroup.Delete("/remove-users", adminService.RemoveUsers)
	adminGroup.Post("/signout-users", adminService.SignoutUsers)
	adminGroup.Post("/suspend-users", adminService.SuspendUsers)
	adminGroup.Post("/unsuspend-users", adminService.UnsuspendUsers)

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
