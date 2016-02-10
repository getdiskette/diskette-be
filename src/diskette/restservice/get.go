package restservice

import (
	"diskette/util"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
)

// examples:
// http localhost:5025/user
// http localhost:5025/user?q='{"name":"dfreire"}'
// http localhost:5025/user?q='{"name":{"$ne":"dfreire"}}'
func (self *impl) Get(c *echo.Context) error {
	collection := c.Param("collection")

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
