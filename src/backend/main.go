package main

import (
	"net/http"

	"github.com/labstack/echo"

	"log"

	"labix.org/v2/mgo"
)

func main() {

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	e := echo.New()

	// e.Use(func(c *echo.Context) error {
	// 	return nil
	// })

	e.Post("/:database/:collection/new", func(c *echo.Context) error {
		database := c.Param("database")
		collection := c.Param("collection")
		var document map[string]interface{}
		c.Bind(&document)

		err := session.DB(database).C(collection).Insert(document)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, createErrorResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, createOkResponse(document))
	})

	e.Get("/:database/:collection/list", func(c *echo.Context) error {
		database := c.Param("database")
		collection := c.Param("collection")

		var documents []interface{}
		err := session.DB(database).C(collection).Find(nil).All(&documents)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, createErrorResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, createOkResponse(documents))
	})

	e.Run(":5025")
}

func createOkResponse(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m["ok"] = true
	m["data"] = data
	return m
}

func createErrorResponse(error string) map[string]interface{} {
	m := make(map[string]interface{})
	m["ok"] = false
	m["error"] = error
	return m
}
