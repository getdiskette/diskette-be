package main

import (
	"github.com/labstack/echo"

	"labix.org/v2/mgo"
)

func main() {

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	e := echo.New()
	e.Run(":5025")
}
