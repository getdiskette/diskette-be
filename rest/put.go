package rest

import (
	"github.com/getdiskette/diskette/util"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

// examples:
// http PUT localhost:5025/collection/user?q='{"name":"dfreire"}' \$set:='{"name":"dariofreire"}'
func (service *serviceImpl) Put(c *echo.Context) error {
	collection := c.Param("collection")

	queryStr := c.Query("q")
	if queryStr == "" {

		return c.JSON(http.StatusForbidden, util.CreateErrResponse(errors.New("Missing parameter 'q' (for query)")))
	}

	var query map[string]interface{}
	if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	var partialDoc map[string]interface{}
	c.Bind(&partialDoc)

	_, err := service.db.C(collection).UpdateAll(query, partialDoc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(partialDoc))
}
