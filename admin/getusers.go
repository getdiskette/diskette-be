package admin

import (
	"diskette/util"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
)

// http localhost:5025/admin/get-users?q=<query> X-Diskette-Session-Token:<session_token>
func (service *serviceImpl) GetUsers(c *echo.Context) error {
	var query map[string]interface{}
	queryStr := c.Query("q")
	if queryStr != "" {
		if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
			return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
		}
	}

	var documents []interface{}
	err := service.userCollection.Find(query).All(&documents)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(documents))
}
