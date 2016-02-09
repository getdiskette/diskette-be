package rest

import (
	"diskette/util"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"labix.org/v2/mgo"
)

type RestService interface {
	Get(c *echo.Context) error
	Post(c *echo.Context) error
	Put(c *echo.Context) error
	Delete(c *echo.Context) error
}

type impl struct {
	db *mgo.Database
}

func NewRestService(db *mgo.Database) RestService {
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
			return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
		}
	}

	var documents []interface{}
	err := self.db.C(collection).Find(query).All(&documents)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(documents))
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
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(document))
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
			return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
		}
	}

	var partialDoc map[string]interface{}
	c.Bind(&partialDoc)

	_, err := self.db.C(collection).UpdateAll(query, partialDoc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(partialDoc))
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
			return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
		}
	}

	_, err := self.db.C(collection).RemoveAll(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
