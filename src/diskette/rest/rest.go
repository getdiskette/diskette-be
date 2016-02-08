package rest

import (
	"encoding/json"
	"net/http"

	"diskette/vendor/github.com/labstack/echo"
	"diskette/vendor/labix.org/v2/mgo"
)

type Rest interface {
	Get(c *echo.Context) error
	Post(c *echo.Context) error
	Put(c *echo.Context) error
	Delete(c *echo.Context) error
}

type impl struct {
	db *mgo.Database
}

func NewRest(db *mgo.Database) Rest {
	return &impl{db}
}

// GET /collection?st={sessionToken}&q={query}
// examples:
// http localhost:5025/user
// http localhost:5025/user?q='{"name":"dfreire"}'
// http localhost:5025/user?q='{"name":{"$ne":"dfreire"}}'
func (self *impl) Get(c *echo.Context) error {
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
	err := self.db.C(collection).Find(query).All(&documents)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, createOkResponse(documents))
}

// POST /collection?st={sessionToken} BODY={doc}
// examples:
// http POST localhost:5025/user name=dfreire email=dario.freire@gmail.com
func (self *impl) Post(c *echo.Context) error {
	collection := c.Param("collection")
	// sessionToken := c.Query("st")

	var document map[string]interface{}
	c.Bind(&document)

	err := self.db.C(collection).Insert(document)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, createOkResponse(document))
}

// PUT /collection?st={sessionToken}&q={query} BODY={partialDoc}
// examples:
// http PUT localhost:5025/user?q='{"name":"dfreire"}' \$set:='{"name":"dariofreire"}'
func (self *impl) Put(c *echo.Context) error {
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

	_, err := self.db.C(collection).UpdateAll(query, partialDoc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, createOkResponse(partialDoc))
}

// DELETE /collection?st={sessionToken}&q={query}
// examples:
// http DELETE localhost:5025/user?q='{"name":"dfreire"}'
func (self *impl) Delete(c *echo.Context) error {
	collection := c.Param("collection")
	// sessionToken := c.Query("st")

	queryStr := c.Query("q")
	var query map[string]interface{}
	if queryStr != "" {
		if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	_, err := self.db.C(collection).RemoveAll(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, createOkResponse(nil))
}

func createOkResponse(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m["ok"] = true
	if data != nil {
		m["data"] = data
	}
	return m
}
