package rest

import (
	"encoding/json"
	"errors"
	"github.com/getdiskette/diskette/util"
	"net/http"

	"github.com/labstack/echo"
)

// examples:
// http DELETE localhost:5025/collection/user?q='{"name":"dfreire"}'
func (service *serviceImpl) Delete(c echo.Context) error {
	collection := c.Param("collection")

	queryStr := c.QueryParam("q")
	if queryStr == "" {

		return c.JSON(http.StatusForbidden, util.CreateErrResponse(errors.New("Missing parameter 'q' (for query)")))
	}

	var query map[string]interface{}
	if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	_, err := service.db.C(collection).RemoveAll(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
