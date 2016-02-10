package restservice

import (
	"diskette/util"
	"net/http"

	"github.com/labstack/echo"
)

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
