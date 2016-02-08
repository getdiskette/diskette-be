package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo"
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

	// GET /db/col?st={sessionToken}&q={query}
	// examples:
	// http localhost:5025/test/user
	// http localhost:5025/test/user?q='{"name":"dfreire"}'
	// http localhost:5025/test/user?q='{"name":{"$ne":"dfreire"}}'
	e.Get("/:database/:collection", func(c *echo.Context) error {
		database := c.Param("database")
		collection := c.Param("collection")
		// sessionToken := c.Query("st")

		var query map[string]interface{}
		queryStr := c.Query("q")
		if queryStr != "" {
			if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
				return c.JSON(http.StatusInternalServerError, createErrorResponse(err.Error()))
			}
		}

		var documents []interface{}
		err := session.DB(database).C(collection).Find(query).All(&documents)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, createErrorResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, createOkResponse(documents))
	})

	// POST /db/col?st={sessionToken} BODY={doc}
	// examples:
	// http POST localhost:5025/test/user name=dfreire email=dario.freire@gmail.com
	e.Post("/:database/:collection", func(c *echo.Context) error {
		database := c.Param("database")
		collection := c.Param("collection")
		// sessionToken := c.Query("st")

		var document map[string]interface{}
		c.Bind(&document)

		err := session.DB(database).C(collection).Insert(document)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, createErrorResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, createOkResponse(document))
	})

	// PUT /db/col?st={sessionToken}&q={query} BODY={partialDoc}
	// examples:
	// http PUT localhost:5025/test/user?q='{"name":"dfreire"}' \$set:='{"name":"dariofreire"}'
	e.Put("/:database/:collection", func(c *echo.Context) error {
		database := c.Param("database")
		collection := c.Param("collection")
		// sessionToken := c.Query("st")

		queryStr := c.Query("q")
		var query map[string]interface{}
		if queryStr != "" {
			if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
				return c.JSON(http.StatusInternalServerError, createErrorResponse(err.Error()))
			}
		}

		var partialDoc map[string]interface{}
		c.Bind(&partialDoc)

		_, err := session.DB(database).C(collection).UpdateAll(query, partialDoc)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, createErrorResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, createOkResponse(partialDoc))
	})

	// DELETE /db/col?st={sessionToken}&q={query}
	// examples:
	// http DELETE localhost:5025/test/user?q='{"name":"dfreire"}'
	e.Delete("/:database/:collection", func(c *echo.Context) error {
		database := c.Param("database")
		collection := c.Param("collection")
		// sessionToken := c.Query("st")

		queryStr := c.Query("q")
		var query map[string]interface{}
		if queryStr != "" {
			if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
				return c.JSON(http.StatusInternalServerError, createErrorResponse(err.Error()))
			}
		}

		_, err := session.DB(database).C(collection).RemoveAll(query)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, createErrorResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, createOkResponse(nil))
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
