package rest

import (
	"encoding/json"
	"github.com/getdiskette/diskette/util"
	"net/http"

	"github.com/labstack/echo"
)

// examples:
// http localhost:5025/collection/user
// http localhost:5025/collection/user?q='{"name":"dfreire"}'
// http localhost:5025/collection/user?q='{"name":{"$ne":"dfreire"}}'
func (service *serviceImpl) Get(c echo.Context) error {
	collection := c.Param("collection")

	var query map[string]interface{}
	queryStr := c.QueryParam("q")
	if queryStr != "" {
		if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
			return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
		}
	}

	var documents []interface{}
	err := service.db.C(collection).Find(query).All(&documents)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(documents))
}
