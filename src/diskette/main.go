package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"diskette/vendor/labix.org/v2/mgo"

	"github.com/labstack/echo"
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

	e := echo.New()

	// e.Use(func(c *echo.Context) error {
	// 	return nil
	// })

	// GET /collection?st={sessionToken}&q={query}
	// examples:
	// http localhost:5025/user
	// http localhost:5025/user?q='{"name":"dfreire"}'
	// http localhost:5025/user?q='{"name":{"$ne":"dfreire"}}'
	e.Get("/:collection", func(c *echo.Context) error {
		collection := c.Param("collection")
		// sessionToken := c.Query("st")

		var query map[string]interface{}
		queryStr := c.Query("q")
		if queryStr != "" {
			if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		var documents []interface{}
		err := db.C(collection).Find(query).All(&documents)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, createOkResponse(documents))
	})

	// POST /collection?st={sessionToken} BODY={doc}
	// examples:
	// http POST localhost:5025/user name=dfreire email=dario.freire@gmail.com
	e.Post("/:collection", func(c *echo.Context) error {
		collection := c.Param("collection")
		// sessionToken := c.Query("st")

		var document map[string]interface{}
		c.Bind(&document)

		err := db.C(collection).Insert(document)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, createOkResponse(document))
	})

	// PUT /collection?st={sessionToken}&q={query} BODY={partialDoc}
	// examples:
	// http PUT localhost:5025/user?q='{"name":"dfreire"}' \$set:='{"name":"dariofreire"}'
	e.Put("/:collection", func(c *echo.Context) error {
		collection := c.Param("collection")
		// sessionToken := c.Query("st")

		queryStr := c.Query("q")
		var query map[string]interface{}
		if queryStr != "" {
			if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		var partialDoc map[string]interface{}
		c.Bind(&partialDoc)

		_, err := db.C(collection).UpdateAll(query, partialDoc)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, createOkResponse(partialDoc))
	})

	// DELETE /collection?st={sessionToken}&q={query}
	// examples:
	// http DELETE localhost:5025/user?q='{"name":"dfreire"}'
	e.Delete("/:collection", func(c *echo.Context) error {
		collection := c.Param("collection")
		// sessionToken := c.Query("st")

		queryStr := c.Query("q")
		var query map[string]interface{}
		if queryStr != "" {
			if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		_, err := db.C(collection).RemoveAll(query)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, createOkResponse(nil))
	})

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

func createOkResponse(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m["ok"] = true
	if data != nil {
		m["data"] = data
	}
	return m
}
